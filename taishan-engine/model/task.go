package model

const (
	StopPlan     = 1
	DebugStatus  = 2
	ReportChange = 3
	SceneRelease = 4
)

// ReportStatusChange 订阅压测计划状态变更
type ReportStatusChange struct {
	Type             int              `json:"type"` // 1: stopPlan; 2: debug; 3.报告变更
	StopPlan         string           `json:"stop_plan"`
	Debug            string           `json:"debug"`
	ActionChangeInfo ActionChangeInfo `json:"action_change_info"`
}

type ActionChangeInfo struct {
	Concurrency int64 `json:"concurrency"`
	Rps         int64 `json:"rps"`
	//Duration    int64 `json:"duration"`
}
