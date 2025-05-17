package model

import (
	"engine/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	ossClient *oss.Client
)

func MustInitOSS() {
	// 读取配置文件
	c := config.Conf
	ossConfig := c.OSSConfig

	// 创建 OSS 客户端
	endpoint := ossConfig.Endpoint
	accessKeyID := ossConfig.AccessKeyID
	accessKeySecret := ossConfig.AccessKeySecret
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic(fmt.Errorf("fatal error OSS init: %w", err))
	}

	// 将 OSS 客户端保存到全局变量中，以供后续使用
	ossClient = client

	fmt.Println("OSS initialized")
}

func GetOSSClient() *oss.Client {
	return ossClient
}

func GetTaishanBucket() (*oss.Bucket, error) {
	// 获取存储
	bucketName := config.Conf.OSSConfig.BucketName
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
