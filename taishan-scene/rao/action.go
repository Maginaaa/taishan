package rao

type Action struct {
	Plan               *Plan             `json:"plan"`
	ReportID           int32             `json:"report_id"`
	ReportName         string            `json:"report_name"`
	ScenePressInfoList []*ScenePressInfo `json:"scene_press_info_list"`
	FileInfo           []FileInfo        `json:"file_info"`
	EngineSerialNumber int32             `json:"engine_serial_number"` // 压测引擎序号
	EngineCount        int32             `json:"engine_count"`         // 压测引擎的总数量
	PartitionID        int32             `json:"partition_id"`         // 压测引擎的分区id
}

type SceneAction struct {
	PlanID     int32          `json:"plan_id"`
	PlanName   string         `json:"plan_name"`
	ReportID   int32          `json:"report_id"`
	ReportName string         `json:"report_name"`
	Scene      ScenePressInfo `json:"scene_list"`
	LogType    int32          `json:"log_type"` // 0: 不开启， 1,开启-仅错误日志， 2，开启-所有日志
	PressInfo  PressInfo      `json:"press_info"`
	BreakType  int32          `json:"break_type"`
	BreakValue float32        `json:"break_value"`
}

type ScenePressInfo struct {
	Scene       *Scene `json:"scene"`
	Rate        int64  `json:"rate"`
	Concurrency int64  `json:"concurrency"`
	Iteration   int64  `json:"iteration"`
	PartitionID int32  `json:"partition_id"`
}

type SceneCases []*SceneCaseTree

type Concurrency struct {
	Concurrency int64 `json:"concurrency"`
	DurationMin int64 `json:"duration_min"`
}

type StepPressInfo struct {
	StartConcurrent int `json:"start_concurrent"`
	StepSize        int `json:"step_size"`
	StepDuration    int `json:"step_duration"`
	//MaxConcurrent   int `json:"max_concurrent"`
	DurationMin int `json:"duration_min"`
}

// ReportStatusChange 订阅压测计划状态变更
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
