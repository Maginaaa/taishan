package logic

import (
	"report/rao"
)

func GetUserList() ([]rao.User, error) {
	return []rao.User{
		{
			ID:        1,
			Name:      "简单随风",
			AvatarURL: "https://avatars.githubusercontent.com/u/19279437",
		},
	}, nil
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
