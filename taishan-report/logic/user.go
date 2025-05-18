package logic

import (
	"encoding/json"
	"net/http"
	"report/conf"
	"report/rao"
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

	// 解析响应体
	var userListResponse struct {
		Code int64      `json:"code"`
		Em   string     `json:"em"`
		Et   string     `json:"et"`
		Data []rao.User `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userListResponse)
	if err != nil {
		return nil, err
	}

	return userListResponse.Data, nil
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
