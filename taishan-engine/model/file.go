package model

import (
	"encoding/csv"
	"engine/internal/biz/log"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"sync"
)

// 参数化文件读取类型
const (
	OrderedRead = iota // 顺序读取
	RandomRead         // 随机读取
	PollingOnce
	PlanFilePath = "plan"
)

type FileList struct {
	File []*FileInfo
}

const (
	FileDataType = iota
	ExportDataType
)

// TODO:做文件管理后，将planID拆出去后，需要将DataMap也拆出去
type FileInfo struct {
	Name      string   `json:"name"`
	Rows      int32    `json:"rows"`
	Column    []Column `json:"column"`
	Status    bool     `json:"status"`
	DataType  int      // 0: 文件数据， 1: 导出数据
	DataMap   map[string]*Param
	dataMutex sync.Mutex
}

type Column struct {
	Col      string `json:"col"`
	Alias    string `json:"alias"`
	ReadType int    `json:"read_type"`
}

type Param struct {
	Val      []string
	ReadType int
	index    int
	rows     int
	//indexMutex sync.Mutex
}

func (f *FileList) ParseFile(planId int32, engineSerial, engineCount int32) {
	if len(f.File) == 0 {
		return
	}
	fileBucket, getBucketErr := GetTaishanBucket()
	if getBucketErr != nil {
		log.Logger.Errorf("Failed to get Taishan bucket: %v", getBucketErr)
		return
	}
	// 多个文件遍历处理
	for _, file := range f.File {
		if file.DataType != FileDataType {
			continue
		}
		body, err := fileBucket.GetObject(filepath.Join(PlanFilePath, strconv.Itoa(int(planId)), file.Name))
		if err != nil {
			log.Logger.Error("Error downloading CSV file:", err)
			return
		}
		rows, err := csv.NewReader(body).ReadAll()
		if err != nil {
			return
		}
		file.setFileData(rows, engineSerial, engineCount)
	}
}

func (f *FileList) ParsePreSceneFile(sceneId int32, exportDataInfo ExportDataInfo, engineSerial, engineCount int32) {
	fileBucket, getBucketErr := GetTaishanBucket()
	if getBucketErr != nil {
		log.Logger.Errorf("Failed to get Taishan bucket: %v", getBucketErr)
		return
	}
	body, err := fileBucket.GetObject(fmt.Sprintf("export/scene/%d.csv", sceneId))
	rows, err := csv.NewReader(body).ReadAll()
	if err != nil {
		return
	}

	file := &FileInfo{
		DataType: ExportDataType,
		Column:   make([]Column, len(exportDataInfo.VariableList)),
	}
	for i, vrb := range rows[0] {
		file.Column[i] = Column{
			Col:      vrb,
			ReadType: OrderedRead,
		}
	}
	file.setFileData(rows, engineSerial, engineCount)
	f.File = append(f.File, file)
}

func (f *FileList) GetVariable() map[string]string {
	mp := make(map[string]string)
	for _, file := range f.File {
		file.dataMutex.Lock()
		for key, value := range file.DataMap {
			mp[key] = value.getValue()
		}
		file.dataMutex.Unlock()
	}
	return mp
}

func (f *FileInfo) setFileData(records [][]string, engineSerial, engineCount int32) {
	f.DataMap = make(map[string]*Param)
	keyArray := make([]string, 0)
	if len(records) <= 1 {
		return
	}
	maxLen := int32(len(records) - 1)
	if engineSerial > maxLen {
		engineSerial = engineSerial % maxLen
	}
	for line, row := range records {
		// 参数化文件切片计算
		if line != 0 && int32(line-1)%engineCount != engineSerial {
			continue
		}
		if line == 0 {
			for columnIndex, _ := range row {
				key := f.Column[columnIndex].Col
				if f.Column[columnIndex].Alias != "" {
					key = f.Column[columnIndex].Alias
				}
				f.DataMap[key] = &Param{
					ReadType: f.Column[columnIndex].ReadType,
				}
				keyArray = append(keyArray, key)
			}
			continue
		}
		for columnIndex, val := range row {
			f.DataMap[keyArray[columnIndex]].add(val)
		}
	}
}

func (p *Param) add(val string) {
	if p.Val == nil {
		p.Val = make([]string, 0)
	}
	p.Val = append(p.Val, val)
}

func (p *Param) getValue() string {
	if p.Val == nil || len(p.Val) == 0 {
		return ""
	}
	i := p.nextIndex()
	if i < 0 {
		return ""
	}
	return p.Val[i]
}

func (p *Param) nextIndex() (ind int) {
	if p.rows == 0 {
		p.rows = len(p.Val)
	}
	switch p.ReadType {
	case OrderedRead:
		ind = p.index
		if p.index == p.rows-1 {
			p.index = 0
		} else {
			p.index++
		}
		break
	case RandomRead:
		ind = rand.Intn(p.rows)
		break
	case PollingOnce:
		ind = p.index
		if p.index == p.rows-1 {
			p.index = -999
		} else {
			p.index++
		}
	default:
		ind = 0
	}

	return
}
