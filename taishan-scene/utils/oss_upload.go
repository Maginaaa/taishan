package utils

import (
	"mime/multipart"
	"path/filepath"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"strconv"
)

const (
	PlanFilePath = "plan"
)

func UploadPlanFile(planId int32, file multipart.FileHeader) (err error) {

	bucket, err := dal.GetTaishanBucket()
	if err != nil {
		log.Logger.Error("utils.oss_upload.UploadPlanFile.GetTaishanBucket, err:", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		log.Logger.Error("utils.oss_upload.UploadPlanFile.fileOpen, err:", err)
		return err
	}
	defer src.Close()

	// 上传文件
	err = bucket.PutObject(filepath.Join(PlanFilePath, strconv.Itoa(int(planId)), file.Filename), src)
	// 获取上传后的文件 URL
	if err == nil {
		log.Logger.Infof("%s文件上传成功", file.Filename)
	}
	return
}

func CopyPLanFile(srcPlanId, destPlanId int32, fileName string) (err error) {
	bucket, err := dal.GetTaishanBucket()
	if err != nil {
		log.Logger.Error("utils.oss_upload.CopyPLanFile.GetTaishanBucket, err:", err)
		return
	}
	srcPath := filepath.Join(PlanFilePath, strconv.Itoa(int(srcPlanId)), fileName)
	destPath := filepath.Join(PlanFilePath, strconv.Itoa(int(destPlanId)), fileName)
	_, err = bucket.CopyObject(srcPath, destPath)
	return
}
