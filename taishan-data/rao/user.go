package rao

type User struct {
	AvatarURL   string      `json:"avatarUrl"`
	Channel     int         `json:"channel"`
	CreatedTime string      `json:"created_time"`
	Email       string      `json:"email"`
	FsUserID    string      `json:"fs_user_id"`
	ID          int32       `json:"id"`
	IsDeleted   int         `json:"isDeleted"`
	Name        string      `json:"name"`
	Nickname    string      `json:"nickname"`
	OpenID      string      `json:"open_id"`
	Phone       string      `json:"phone"`
	Roles       interface{} `json:"roles"`
	Status      int         `json:"status"`
}
