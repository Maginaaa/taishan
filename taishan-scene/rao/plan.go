package rao

import (
	"encoding/json"
	"scene/internal/model"
	"sort"
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
	ConcurrentModel = iota // 并发数模式
	StepModel              // 阶梯模式
	CaseRpsMode
	PlanRpsRateMode
	FixedFrequency
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
	PlanInfo     string  `json:"plan_info"`      // 计划名或id，非必填
	CaseInfo     string  `json:"case_info"`      // 接口名/路径，非必填
	CreateUserId int32   `json:"create_user_id"` // 创建者id，非必填
	TagList      []int32 `json:"tag_list"`
	Page         int     `json:"page" binding:"required,gt=0"`      // 页码
	PageSize     int     `json:"page_size" binding:"required,gt=0"` // 每页条数
}

type GetPlanListResponse struct {
	PlanList []*Plan `json:"plan_list"`
	Total    int64   `json:"total"`
}

type SceneInfo struct {
	SceneId        int32            `json:"scene_id"`
	Sort           int32            `json:"sort"`
	SceneType      int32            `json:"scene_type"`
	ExportDataInfo *ExportDataInfo  `json:"export_data_info,omitempty"`
	Disabled       bool             `json:"disabled"`
	SceneName      string           `json:"scene_name"`
	Rate           int64            `json:"rate,omitempty"`
	RpsRate        int64            `json:"rps_rate,omitempty"`
	Concurrency    int64            `json:"concurrency,omitempty"`
	Iteration      int64            `json:"iteration,omitempty"`
	CaseTree       []*SceneCaseTree `json:"case_tree"`
	VariableList   []*Variable      `json:"variable_list"`
}

type Plan struct {
	PlanID         int32         `json:"plan_id"`
	PlanName       string        `json:"plan_name"`
	IsRunning      bool          `json:"is_running"`
	ReportID       int32         `json:"report_id"`
	SceneList      []SceneInfo   `json:"scene_list"`
	Partition      int32         `json:"partition"`
	EngineCount    int32         `json:"engine_count"`
	PressInfo      PressInfo     `json:"press_info"`
	SamplingInfo   SamplingInfo  `json:"sampling_info"`
	GlobalVariable []*ParamsForm `json:"global_variable"`
	DefaultHeader  []*ParamsForm `json:"default_header"`
	BreakType      int32         `json:"break_type"`
	BreakValue     float32       `json:"break_value"`
	TaskInfo       TaskInfo      `json:"task_info"`
	ServerInfo     [][]string    `json:"server_info"` // 管理服务信息 [[namespace, serverName]]
	CapacitySwitch bool          `json:"capacity_switch"`
	CreateUserName string        `json:"create_user_name"` // 创建人
	CreateTime     string        `json:"create_time"`
	UpdateUserName string        `json:"update_user_name"` // 最后修改人
	UpdateTime     string        `json:"update_time"`      // 最后修改时间
	TagList        []int32       `json:"tag_list"`
	Remark         string        `json:"remark"`
	PressCount     int           `json:"press_count"`     // 累积压测次数
	LastPressTime  string        `json:"last_press_time"` // 最后压测时间
	VariablePool   sync.Map      `json:"variable_pool"`   // 变量池
	DebugSuccess   bool          `json:"debug_success"`
}

type PlanVariable struct {
	VariableName string `json:"variable_name"` // 变量名
	VariableVal  string `json:"variable_val"`  // 变量值
}

func (p *Plan) InitGlobalVariable() {
	for _, vrb := range p.GlobalVariable {
		if !vrb.Enable {
			continue
		}
		p.VariablePool.Store(vrb.Key, vrb.Value)
	}
}

func (p *Plan) AppendVariablePool(variableList []Variable) {
	for _, v := range variableList {
		p.VariablePool.Store(v.VariableName, v.VariableVal)
	}
}

func (p *Plan) SetSceneList(sceneList []*model.Scene, caseMp map[int32][]*SceneCaseTree, vrbMap map[int32][]*Variable) {
	newSceneList := make([]SceneInfo, 0)
	sort.SliceStable(sceneList, func(i, j int) bool {
		return sceneList[i].SceneType > sceneList[j].SceneType
	})
	for _, v := range sceneList {
		scene := SceneInfo{
			SceneId:      v.ID,
			SceneType:    v.SceneType,
			SceneName:    v.SceneName,
			Disabled:     v.Disabled,
			CaseTree:     caseMp[v.ID],
			VariableList: vrbMap[v.ID],
		}
		if v.SceneType == PreFixScene {
			var exportInfo ExportDataInfo
			_ = json.Unmarshal([]byte(v.ExportInfo), &exportInfo)
			scene.ExportDataInfo = &exportInfo
		}
		newSceneList = append(newSceneList, scene)
	}
	p.SceneList = newSceneList
}

//func (p *Plan) SetMachineList(machineListStr string) (err error) {
//	var machineList []string
//	err = json.Unmarshal([]byte(machineListStr), &machineList)
//	if err != nil {
//		log.Logger.Error("logic.plan.SetMachineList.jsonUnmarshal, err:", err)
//		return
//	}
//	p.MachineList = machineList
//	return nil
//}

type PressInfo struct {
	PressType        int              `json:"press_type"`
	Concurrency      int64            `json:"concurrency"`
	RPS              int64            `json:"rps"`
	Duration         int64            `json:"duration"` // 单位为分钟
	StartConcurrency int64            `json:"start_concurrency"`
	StepSize         int64            `json:"step_size"`
	StepDuration     int64            `json:"step_duration"` // 单位为秒
	SceneList        []PressSceneInfo `json:"scene_list"`
}

type SamplingInfo struct {
	SamplingType int `json:"sampling_type"`
	// 采样率
	SamplingRate int64 `json:"sampling_rate"`
}

type ConcurrencyChange struct {
	ReportID    int32 `json:"report_id"`
	SceneID     int32 `json:"scene_id"`
	Concurrency int64 `json:"concurrency"`
}

type PlanDebugRecordReq struct {
	PlanID   int32 `json:"plan_id"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

type PlanDebugRecordRes struct {
	Total      int64             `json:"total"`
	RecordList []PlanDebugRecord `json:"record_list"`
}

type PlanDebugRecord struct {
	Time   string                 `json:"time"`   // 记录时间
	Passed bool                   `json:"passed"` // 当前调试是否通过
	Result []SceneExecutionResult `json:"result"`
}

type TaskInfo struct {
	TaskId int32  `json:"task_id"`
	Enable bool   `json:"enable"`
	Cron   string `json:"cron"`
}

type PressSceneInfo struct {
	SceneId     int32 `json:"scene_id"`
	SceneType   int32 `json:"scene_type"`
	Rate        int64 `json:"rate,omitempty"`
	RpsRate     int64 `json:"rps_rate,omitempty"`
	Concurrency int64 `json:"concurrency,omitempty"`
	Iteration   int64 `json:"iteration,omitempty"`
}
