package rao

type Role struct {
	RoleID   int    `json:"role_id"`   // 角色ID
	RoleName string `json:"role_name"` // 角色名称
	Code     string `json:"code"`      // 角色代码
	RoleType int    `json:"role_type"` // 角色类型
}

type User struct {
	ID          int    `json:"id"`              // 用户ID
	Name        string `json:"name"`            // 姓名
	Nickname    string `json:"nickname"`        // 昵称
	Channel     int    `json:"channel"`         // 渠道
	Phone       string `json:"phone"`           // 手机号码
	Email       string `json:"email,omitempty"` // 邮箱
	AvatarURL   string `json:"avatarUrl"`       // 头像URL
	IsDeleted   int    `json:"isDeleted"`       // 是否已删除
	Status      int    `json:"status"`          // 状态
	CreatedTime string `json:"created_time"`    // 创建时间
	FSUserID    string `json:"fs_user_id"`      // FS用户ID
	OpenID      string `json:"open_id"`         // Open ID
	Roles       []Role `json:"roles"`           // 角色列表
}

type Response struct {
	Users []User `json:"users"` // 用户列表
}
