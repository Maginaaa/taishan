package logic

import (
	"context"
	"data/config"
	"data/internal/biz/log"
	"data/internal/dal"
	"data/rao"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
)

const (
	MaxSaveCaseSampling   = 10000
	SamplingDataPartition = 10
	PrefixDataPartition   = 11 // 采样日志默认写死由partition11进行传递
)

func ConsumerSamplingData() {
	consumer, err := sarama.NewConsumer(strings.Split(config.Conf.Kafka.Address, ","), sarama.NewConfig())
	pc, err := consumer.ConsumePartition(config.Conf.Kafka.Topic, int32(SamplingDataPartition), sarama.OffsetNewest)
	defer pc.AsyncClose()
	if err != nil {
		log.Logger.Error("consumer.ConsumePartition error:", err)
		return
	}

	// 记录plan对应的report
	reportMp := make(map[int32]int32)
	// 一级key为reportId
	caseMp := make(map[int32]*map[int32]int32)

	for msg := range pc.Messages() {
		var resultDataMsg rao.ResultDataMsg
		err = json.Unmarshal(msg.Value, &resultDataMsg)
		if err != nil {
			log.Logger.Errorf("pc.Messages().jsonUnmarshal error: %s, \nvalue: %s", err.Error(), string(msg.Value))
		}
		dt := resultDataMsg.HttpResultDataMsg.DebugInfo
		dt.SceneId = resultDataMsg.SceneInfo.SceneID

		planId := resultDataMsg.GlobalInfo.ReportID
		reportId := resultDataMsg.GlobalInfo.ReportID
		caseId := resultDataMsg.HttpResultDataMsg.CaseID

		// 清除历史缓存
		if oldReportId, ok := reportMp[planId]; !ok {
			reportMp[planId] = reportId
			log.Logger.Infof("开始处理报告%d的采样数据", reportId)
		} else {
			if oldReportId != reportId {
				delete(caseMp, oldReportId)
			}
		}
		// 采样计数
		if _, ok := caseMp[reportId]; !ok {
			caseMp[resultDataMsg.GlobalInfo.ReportID] = &map[int32]int32{}
		}
		if _, ok := (*caseMp[resultDataMsg.GlobalInfo.ReportID])[caseId]; !ok {
			(*caseMp[resultDataMsg.GlobalInfo.ReportID])[caseId] = 0
		}
		(*caseMp[resultDataMsg.GlobalInfo.ReportID])[caseId]++
		if (*caseMp[resultDataMsg.GlobalInfo.ReportID])[caseId] <= MaxSaveCaseSampling {
			dal.GetMongoCollection(fmt.Sprintf("%d", reportId)).InsertOne(context.Background(), dt)
		}

	}

}

func ConsumerTransferData() {
	consumer, err := sarama.NewConsumer(strings.Split(config.Conf.Kafka.Address, ","), sarama.NewConfig())
	pc, err := consumer.ConsumePartition(config.Conf.Kafka.Topic, int32(PrefixDataPartition), sarama.OffsetNewest)
	defer pc.AsyncClose()
	if err != nil {
		log.Logger.Error("consumer.ConsumePartition error:", err)
		return
	}

	dataMap := make(map[int32]*rao.ExportDataInfo) // key为reportId

	for msg := range pc.Messages() {
		var transferInfo rao.TransferInfo
		err = json.Unmarshal(msg.Value, &transferInfo)
		if err != nil {
			log.Logger.Errorf("pc.Messages().jsonUnmarshal error: %s, \nvalue: %s", err.Error(), string(msg.Value))
		}
		switch transferInfo.Type {
		case rao.TypePreSceneExport:
			if _, ok := dataMap[transferInfo.ReportID]; !ok {
				dataMap[transferInfo.ReportID] = &rao.ExportDataInfo{
					MachineIPSet: make(map[string]struct{}),
					TitleInit:    false,
					TitleArray:   make([]string, 0),
					Content:      make([][]string, 0),
				}
				log.Logger.Infof("报告: %d，开始处理前置导出数据", transferInfo.ReportID)
			}
			// machine的初始化和释放
			if transferInfo.End {
				delete(dataMap[transferInfo.ReportID].MachineIPSet, transferInfo.MachineIP)
				if len(dataMap[transferInfo.ReportID].MachineIPSet) == 0 {
					go dal.UploadCsvToOss("export/scene", fmt.Sprintf("%d.csv", transferInfo.SceneID), dataMap[transferInfo.ReportID].Content)
					go dal.UploadCsvToOss("export/report", fmt.Sprintf("%d.csv", transferInfo.ReportID), dataMap[transferInfo.ReportID].Content)
					delete(dataMap, transferInfo.ReportID)
					log.Logger.Infof("报告: %d，场景: %d，处理前置导出数据结束", transferInfo.ReportID, transferInfo.SceneID)
				}
				continue
			} else {
				if _, ok := dataMap[transferInfo.ReportID].MachineIPSet[transferInfo.MachineIP]; !ok {
					dataMap[transferInfo.ReportID].MachineIPSet[transferInfo.MachineIP] = struct{}{}
				}
			}
			// title的初始化
			if !dataMap[transferInfo.ReportID].TitleInit {
				for title := range transferInfo.Data {
					dataMap[transferInfo.ReportID].TitleArray = append(dataMap[transferInfo.ReportID].TitleArray, title)
				}
				dataMap[transferInfo.ReportID].TitleInit = true
				dataMap[transferInfo.ReportID].Content = append(dataMap[transferInfo.ReportID].Content, dataMap[transferInfo.ReportID].TitleArray)
			}
			rowData := make([]string, 0)
			for _, title := range dataMap[transferInfo.ReportID].TitleArray {
				rowData = append(rowData, transferInfo.Data[title])
			}
			dataMap[transferInfo.ReportID].Content = append(dataMap[transferInfo.ReportID].Content, rowData)
		default:
			log.Logger.Info("未知数据类型")

		}
	}
}
