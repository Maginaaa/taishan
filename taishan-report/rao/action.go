package rao

type ReportStatusChange struct {
	Type             int              `json:"type"` // 1: stopPlan; 2: debug; 3.报告变更
	StopPlan         string           `json:"stop_plan"`
	Debug            string           `json:"debug"`
	ActionChangeInfo ActionChangeInfo `json:"action_change_info"`
}

type ActionChangeInfo struct {
	Concurrency int64 `json:"concurrency"`
	Duration    int64 `json:"duration"`
}
