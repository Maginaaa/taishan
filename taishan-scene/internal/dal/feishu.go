package dal

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"scene/internal/conf"
)

var (
	FsClient *lark.Client
)

func InitFeishu() {
	FsClient = lark.NewClient(conf.Conf.FsConfig.AppID, conf.Conf.FsConfig.AppSecret)
}
