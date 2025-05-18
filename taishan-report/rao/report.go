package rao

const (
	PreRunning = iota + 1
	Running
	PreDown
	Down
)

type CreateReportReq struct {
	PlanID    int32  `json:"plan_id"`
	PlanName  string `json:"plan_name"`
	Duration  int32  `json:"duration"`
	PressType int32  `json:"press_type"`
	//MachineNum   int32  `json:"machine_num"`
	EngineList   []string `json:"engine_list"`
	CreateUserID int32    `json:"create_user_id"`
}

type ReportListReq struct {
	PlanID       int32  `json:"plan_id"`                           // 计划id，非必填
	CreateUserId int32  `json:"create_user_id"`                    // 创建者id，非必填
	Page         int    `json:"page" binding:"required,gt=0"`      // 页码
	PageSize     int    `json:"page_size" binding:"required,gt=0"` // 每页条数
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

type CaseReportDataReq struct {
	ReportID int32 `json:"report_id"`
	CaseID   int32 `json:"case_id"`
}

type Report struct {
	ReportID       int32   `json:"report_id"`       // 报告id
	ReportName     string  `json:"report_name"`     // 报告名
	Status         bool    `json:"status"`          // 压测状态
	PlanID         int32   `json:"plan_id"`         // 计划id
	PlanName       string  `json:"plan_name"`       // 计划名
	PressType      int32   `json:"press_type"`      // 压测模式
	Duration       int32   `json:"duration"`        // 预计持续时间
	ActualDuration float64 `json:"actual_duration"` // 实际持续时间
	Concurrency    int64   `json:"concurrency"`     // 并发数
	//MachineNum     int32   `json:"machine_num"`     // 机器数量
	EngineList []string `json:"engine_list"`

	TotalCount  int64   `json:"total_count"`  // 总请求数
	SuccessRate float64 `json:"success_rate"` // 成功率
	StartTime   string  `json:"start_time"`   // 开始时间
	EndTime     string  `json:"end_time"`     // 结束时间

	TotalSendBytes    float64 `json:"total_send_bytes"`    // 发送流量
	TotalReceiveBytes float64 `json:"total_receive_bytes"` // 接受流量

	CreateUserID   int32  `json:"create_user_id"`
	CreateUserName string `json:"create_user_name"` // 创建人
	UpdateUserName string `json:"update_user_name"` // 最后修改人
	UpdateTime     string `json:"update_time"`      // 最后修改时间
}

type ReportPage struct {
	ReportList []*Report `json:"scene_list"`
	Total      int64     `json:"total"`
}
