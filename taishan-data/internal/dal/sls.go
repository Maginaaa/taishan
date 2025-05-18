package dal

import (
	"data/config"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	slsclient "github.com/alibabacloud-go/sls-20201230/v6/client"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	slsClient *slsclient.Client
)

func MustInitSLS() {
	// 读取配置文件
	slsConfig := config.Conf.SLSConfig
	slsConf := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(slsConfig.AccessKeyID),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(slsConfig.AccessKeySecret),
		Endpoint:        tea.String(slsConfig.EndPoint),
	}
	sls, err := slsclient.NewClient(slsConf)
	if err != nil {
		panic(fmt.Errorf("fatal error SLS init: %w", err))
	}
	slsClient = sls
	fmt.Println("sls initialized")
}

func SLSClient() *slsclient.Client {
	return slsClient
}
