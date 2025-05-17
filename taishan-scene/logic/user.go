package logic

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"scene/internal/biz/errno"
	"scene/internal/biz/log"
	"scene/internal/conf"
	"scene/rao"
	"strings"
)

func GetUserList() ([]rao.User, error) {
	// 准备 URL 和请求参数
	url := conf.Conf.Url.Account + "/user/list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// 发起 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: 从redis中获取数据
		return nil, err
	}
	defer resp.Body.Close()

	//var userListRes rao.CommonResponse[[]rao.User]
	//err = json.NewDecoder(resp.Body).Decode(&userListRes)
	//if err != nil {
	//	return nil, err
	//}
	return []rao.User{
		{
			ID:        1,
			Name:      "简单随风",
			AvatarURL: "https://avatars.githubusercontent.com/u/19279437",
		},
	}, nil
}

func UserLogout(ctx *gin.Context) error {
	// 从请求头中获取 Auth 字段的值
	auth := ctx.Request.Header.Get("Authorization")

	// 构造请求数据
	jsonData := map[string]interface{}{
		"token": auth,
	}
	requestBody, err := json.Marshal(jsonData)

	// 创建一个客户端对象
	client := &http.Client{}

	// 创建 POST 请求
	url := conf.Conf.Url.Account + "/auth/logout"
	request, err := http.NewRequest("POST", url, strings.NewReader(string(requestBody)))
	if err != nil {
		log.Logger.Error(http.StatusInternalServerError, err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Logger.Error(http.StatusInternalServerError, err)
		return err
	}
	defer response.Body.Close()

	// 解析响应体
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		// 处理响应体解析错误
		log.Logger.Error(http.StatusInternalServerError, err)
		return err
	}
	code, _ := result["code"]
	if code != float64(0) {
		return errors.New(errno.CodeAlertMap[errno.ErrLogout])
	}
	return nil
}

func GetNameByID(users []rao.User, id int32) string {
	for _, user := range users {
		if int32(user.ID) == id {
			return user.Name
		}
	}
	// 如果找不到ID，则返回一个空字符串
	return ""
}
