package logic

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"io"
	"mime/multipart"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
)

const (
	OrderedRead = iota // 顺序读取
	RandomRead         // 随机读取
)

func UploadFile(ctx *gin.Context, planId int32, file multipart.FileHeader) (err error) {
	err = utils.UploadPlanFile(planId, file)
	if err != nil {
		log.Logger.Error("logic.plan.UploadFileToOSS, err:", err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Logger.Error("logic.plan.fileOpen, err:", err)
		return err
	}
	defer src.Close()
	// 第一次读取前3个字节，检查是否与 UTF-8 BOM 相同
	reader := bufio.NewReader(src)
	bom := []byte{0xEF, 0xBB, 0xBF}
	// 读取文件的前3个字节
	fileBom, err := reader.Peek(3)
	if err != nil && err != io.EOF {
		return err
	}
	// 如果读取到的前3个字节与 BOM 相等，跳过这3个字节
	if bytes.Equal(fileBom, bom) {
		reader.Discard(3)
	}
	// 解析CSV文件，获取列名
	csvReader := csv.NewReader(reader)
	headerRow, err := csvReader.Read()
	// column字段拼装
	var columnList []rao.Column

	// 读取第二行数据
	values, err := csvReader.Read()
	if err != nil {
		log.Logger.Error("failed to read CSV data row: %w", err)
	}
	for idx, col := range headerRow {
		columnList = append(columnList, rao.Column{
			Col:          col,
			FileName:     file.Filename,
			FirstLineVal: values[idx],
			ColIndex:     idx + 1,
			ReadType:     OrderedRead,
		})
	}
	// 计算文件行数
	rowCount := int32(1)
	for {
		_, readErr := csvReader.Read()
		if readErr != nil {
			break
		}
		rowCount++ // 增加行数
	}
	column, _ := json.Marshal(columnList)
	fileInfo := &model.ParameterFile{
		PlanID:       planId,
		FileName:     file.Filename,
		Size:         int32(file.Size),
		Rows:         rowCount,
		Column:       string(column),
		Status:       true,
		CreateUserID: utils.GetCurrentUserID(ctx),
	}
	conditions := make([]gen.Condition, 0)
	// TODO： 删除改为软删除
	tx := dal.GetQuery().ParameterFile
	conditions = append(conditions, tx.PlanID.Eq(fileInfo.PlanID))
	conditions = append(conditions, tx.FileName.Eq(fileInfo.FileName))
	_, err = tx.WithContext(ctx).Where(conditions...).Delete()
	if err != nil {
		log.Logger.Error("failed to delete file", err)
	}
	_ = tx.WithContext(ctx).Create(fileInfo)
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceFile,
		SourceID:      fileInfo.ID,
		OperationType: rao.CreateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    fileInfo,
	})
	return

}

func DownloadFile(ctx *gin.Context, path string) (content []byte, err error) {
	bucket, err := dal.GetTaishanBucket()
	if err != nil {
		log.Logger.Error("logic.file.DownFile.GetTaishanBucket, err:", err)
		return
	}
	object, err := bucket.GetObject(path)
	if err != nil {
		log.Logger.Error("logic.file.DownFile.GetObject, err:", err)
		return
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	content = buf.Bytes()
	return
}

func GetPlanDataSource(ctx *gin.Context, planId int32) (fileInfo []rao.FileInfo, err error) {

	fileInfo = make([]rao.FileInfo, 0)

	tx := dal.GetQuery().ParameterFile
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.PlanID.Eq(planId))
	conditions = append(conditions, tx.IsDelete.Zero())

	fileList, err := tx.WithContext(ctx).Where(conditions...).Order(tx.ID.Desc()).Find()
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDataSource.Find err:", err)
		return
	}

	for _, v := range fileList {
		var columns []rao.Column
		json.Unmarshal([]byte(v.Column), &columns)
		file := rao.FileInfo{
			ID:          v.ID,
			PlanID:      v.PlanID,
			Name:        v.FileName,
			Size:        v.Size,
			Rows:        v.Rows,
			Column:      columns,
			Status:      v.Status,
			CreatedTime: v.CreateTime.Format(FullTimeFormat),
			UpdatedTime: v.UpdateTime.Format(FullTimeFormat),
		}
		fileInfo = append(fileInfo, file)
	}
	return fileInfo, nil

}

func ColumnUpdate(ctx *gin.Context, req rao.ColumnUpdateReq) (err error) {
	tx := dal.GetQuery().ParameterFile
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.PlanID.Eq(req.PlanID))
	conditions = append(conditions, tx.FileName.Eq(req.FileName))
	conditions = append(conditions, tx.IsDelete.Zero())
	// 根据计划id、文件名查询数据库记录
	file, err := tx.WithContext(ctx).Where(conditions...).First()
	if err != nil {
		return err
	}
	// 解析Column字段的JSON数据
	var columns []rao.Column
	err = json.Unmarshal([]byte(file.Column), &columns)
	if err != nil {
		return err
	}
	// 查找指定index并更新alias字段
	for i, col := range columns {
		if col.ColIndex == req.ColIndex {
			columns[i].Alias = req.Alias
			columns[i].ReadType = req.ReadType
			break
		}
	}

	updatedColumnData, err := json.Marshal(columns)
	if err != nil {
		return err
	}
	// 更新数据库记录
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(file.ID)).UpdateSimple(tx.Column.Value(string(updatedColumnData)))
	afterFile, err := tx.WithContext(ctx).Where(tx.ID.Eq(file.ID)).First()
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceFile,
		SourceID:      file.ID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   file,
		ValueAfter:    afterFile,
	})
	return
}

func DeleteParameterFile(ctx *gin.Context, id int32) (isDelete bool, err error) {
	tx := dal.GetQuery().ParameterFile
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(id)).UpdateSimple(tx.IsDelete.Value(true))
	if err != nil {
		log.Logger.Error("logic.plan.DeleteParameterFile，err:", err)
		return false, err
	}
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceFile,
		SourceID:      id,
		OperationType: rao.DeleteOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})
	return true, nil
}

func GetPlanDebugFileVariable(ctx *gin.Context, planId int32) (res []rao.Variable, err error) {
	fileList := getParameterFileInfo(ctx, planId)
	res = make([]rao.Variable, 0)
	for _, file := range fileList {
		var columns []rao.Column
		err = json.Unmarshal([]byte(file.Column), &columns)
		if err != nil {
			log.Logger.Error("logic.plan.GetPlanFileVariable.jsonUnmarshal(Parameter), err：", err)
			return
		}
		for _, column := range columns {
			key := column.Col
			if column.Alias != "" {
				key = column.Alias
			}
			res = append(res, rao.Variable{
				VariableName: key,
				VariableVal:  column.FirstLineVal,
			})
		}
	}
	return
}
func copyPlanFile(ctx *gin.Context, targetPlanId int32, oriPlanId int32) (err error) {
	fileList := getParameterFileInfo(ctx, oriPlanId)
	if len(fileList) == 0 {
		return
	}
	tx := dal.GetQuery().ParameterFile
	for _, file := range fileList {
		fileInfo := &model.ParameterFile{
			PlanID:       targetPlanId,
			FileName:     file.FileName,
			Size:         file.Size,
			Rows:         file.Rows,
			Column:       file.Column,
			Status:       file.Status,
			CreateUserID: utils.GetCurrentUserID(ctx),
		}
		err = tx.WithContext(ctx).Create(fileInfo)
		if err != nil {
			return
		}
		err = utils.CopyPLanFile(oriPlanId, targetPlanId, file.FileName)
		if err != nil {
			return
		}
	}
	return
}

func getParameterFileInfo(ctx *gin.Context, planId int32) (fileList []*model.ParameterFile) {
	tx := dal.GetQuery().ParameterFile
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.PlanID.Eq(planId))
	fileList, err := tx.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		log.Logger.Error("getParameterFileInfo.Find err:", err)
		return
	}
	return
}
