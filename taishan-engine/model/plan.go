package model

import (
	"sync"
)

// LogType 日志类型
const (
	LogClose    = iota // 不开启
	LogOnlyFail        // 开启-错误日志
	LogAll             // 开启-所有日志
)

// PressType 压测类型
const (
	ConcurrentModel     = iota // 并发数模式
	StepModel                  // 阶梯模式
	RpsModel                   // 错误率模式
	RpsRateModel               //rps比例模式
	FixedIterationModel        // 固定执行次数模式
)

// SamplingType 采样类型
const (
	SamplingClose          = iota // 不开启
	SamplingWithError             // 采样-仅错误
	SamplingWithPercentage        // 采样-采样比例， 万分之XX
)

const (
	NotBreak = iota
	ErrorRate
	RtTimeout
)

type PlanListReq struct {
	PlanInfo     string `json:"plan_info"`                         // 计划名或id，非必填
	CreateUserId int32  `json:"create_user_id"`                    // 创建者id，非必填
	Page         int    `json:"page" binding:"required,gt=0"`      // 页码
	PageSize     int    `json:"page_size" binding:"required,gt=0"` // 每页条数
}

type GetPlanListResponse struct {
	PlanList []*Plan `json:"plan_list"`
	Total    int64   `json:"total"`
}

type Plan struct {
	PlanID         int32         `json:"plan_id"`
	PlanName       string        `json:"plan_name"`
	SceneList      []SceneRate   `json:"scene_list"`
	LogType        int32         `json:"log_type"` // 0: 不开启， 1
	PressInfo      PressInfo     `json:"press_info"`
	SamplingInfo   SamplingInfo  `json:"sampling_info"`
	GlobalVariable []*ParamsForm `json:"global_variable"`
	DefaultHeader  []*ParamsForm `json:"default_header"`
	BreakType      int32         `json:"break_type"`
	BreakValue     float32       `json:"break_value"`
	Remark         string        `json:"remark"`
}

type ConcurrencyChange struct {
	ReportID    int32 `json:"report_id"`
	Concurrency int32 `json:"concurrency"`
	DurationMin int64 `json:"duration_min"`
}

type Concurrency struct {
	Concurrency int64 `json:"concurrency"`
	DurationMin int64 `json:"duration_min"`
}

type StepPressInfo struct {
	StartConcurrent int `json:"start_concurrent"`
	StepSize        int `json:"step_size"`
	StepDuration    int `json:"step_duration"`
	MaxConcurrent   int `json:"max_concurrent"`
	DurationMin     int `json:"duration_min"`
}

type Configuration struct {
	Mu sync.Mutex `json:"mu"`
}

type PressInfo struct {
	PressType        int   `json:"press_type"`
	Concurrency      int64 `json:"concurrency"` // 最大并发数/最大RPS
	Duration         int64 `json:"duration"`    // 单位为分钟
	StartConcurrency int64 `json:"start_concurrency"`
	StepSize         int64 `json:"step_size"`
	StepDuration     int64 `json:"step_duration"` // 单位为秒
}

type SamplingInfo struct {
	SamplingType int `json:"sampling_type"`
	// 采样率
	SamplingRate int64 `json:"sampling_rate"`
}
