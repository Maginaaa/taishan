package rao

// DatabaseConnectionInfo 数据库连接所需的信息
type DatabaseConnectionInfo struct {
	Host     string `json:"host" binding:"required"`     // 数据库服务器的主机名或 IP 地址（必填字段）
	Port     string `json:"port" binding:"required"`     // 数据库服务器的端口号（必填字段）
	User     string `json:"user" binding:"required"`     // 数据库连接的用户名（必填字段）
	Password string `json:"password" binding:"required"` // 数据库连接的密码（必填字段）
	DBName   string `json:"dbname" binding:"required"`   // 要连接的数据库名称（必填字段）
	Remark   string `json:"remark"`                      // 数据库连接的备注
}

// CreateConnReq 创建新数据库连接的请求结构
type CreateConnReq struct {
	ID         int32                  `json:"id"`          // 连接的唯一 ID
	Type       int32                  `json:"type"`        // 连接的类型
	PingStatus int32                  `json:"ping_status"` // 连接的 ping 状态
	ConnInfo   DatabaseConnectionInfo `json:"conn_info"`   // 数据库连接信息
}

// DBManageResponse 数据库管理操作的响应结构
type DBManageResponse struct {
	IsSuccess bool `json:"is_success"` // 操作是否成功
}

// DBConnList 检索数据库连接列表的请求结构
type DBConnList struct {
	CreateUserId int32 `json:"create_user_id"`                    // 创建数据库连接的用户 ID
	Page         int   `json:"page" binding:"required,gt=0"`      // 分页的页码（必填字段，必须大于 0）
	PageSize     int   `json:"page_size" binding:"required,gt=0"` // 每页显示的项目数（必填字段，必须大于 0）
}

// DBConnListResp 检索数据库连接列表的响应结构
type DBConnListResp struct {
	ID           int32       `json:"id"`                                // 连接的唯一 ID
	Type         int32       `json:"type"`                              // 连接的类型
	ConnInfo     interface{} `json:"conn_info"`                         // 数据库连接信息
	CreateUserId int32       `json:"create_user_id"`                    // 创建数据库连接的用户 ID
	Page         int         `json:"page" binding:"required,gt=0"`      // 分页的页码（必填字段，必须大于 0）
	PageSize     int         `json:"page_size" binding:"required,gt=0"` // 每页显示的项目数（必填字段，必须大于 0）
}
