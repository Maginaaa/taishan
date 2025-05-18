package rao

// LogType 日志类型
const (
	LogClose    = iota // 不开启
	LogOnlyFail        // 开启-错误日志
	LogAll             // 开启-所有日志
)

// PressType 压测类型
const (
	ConcurrentModel = iota // 并发数模式
	StepModel              // 阶梯模式
	CaseRpsMode     = 2
	PlanRpsRateMode = 3 // 错误率模式
)

type PlanListReq struct {
	PlanInfo     string `json:"plan_info"`                         // 计划名或id，非必填
	CreateUserId int32  `json:"create_user_id"`                    // 创建者id，非必填
	Page         int    `json:"page" binding:"required,gt=0"`      // 页码
	PageSize     int    `json:"page_size" binding:"required,gt=0"` // 每页条数
}

type ConcurrencyChange struct {
	ReportID    int32 `json:"report_id"`
	SceneID     int32 `json:"scene_id"`
	Concurrency int64 `json:"concurrency"`
}

type UpdateReportNameReq struct {
	ReportID   int32  `json:"report_id"`
	ReportName string `json:"report_name"`
}

type Plan struct {
	PlanID         int32       `json:"plan_id"`
	PlanName       string      `json:"plan_name"`
	IsRunning      bool        `json:"is_running"`
	ReportID       int32       `json:"report_id"`
	SceneList      []SceneInfo `json:"scene_list"`
	Partition      int32       `json:"partition"`
	EngineCount    int32       `json:"engine_count"`
	PressInfo      PressInfo   `json:"press_info"`
	BreakType      int32       `json:"break_type"`
	BreakValue     float32     `json:"break_value"`
	CreateUserName string      `json:"create_user_name"` // 创建人
	CreateTime     string      `json:"create_time"`
	UpdateUserName string      `json:"update_user_name"` // 最后修改人
	UpdateTime     string      `json:"update_time"`      // 最后修改时间
	Remark         string      `json:"remark"`
	PressCount     int         `json:"press_count"`     // 累积压测次数
	LastPressTime  string      `json:"last_press_time"` // 最后压测时间
}

type PressInfo struct {
	PressType        int         `json:"press_type"`
	Concurrency      int64       `json:"concurrency"`
	RPS              int64       `json:"rps"`
	Duration         int64       `json:"duration"` // 单位为分钟
	StartConcurrency int64       `json:"start_concurrency"`
	StepSize         int64       `json:"step_size"`
	StepDuration     int64       `json:"step_duration"` // 单位为秒
	SceneList        []SceneInfo `json:"scene_list"`
}
