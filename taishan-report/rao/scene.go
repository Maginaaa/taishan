package rao

const (
	NormalScene = iota
	PreFixScene
)

type SceneInfo struct {
	SceneId        int32          `json:"scene_id"`
	SceneType      int32          `json:"scene_type"`
	ExportDataInfo ExportDataInfo `json:"export_data_info"`
	Disabled       bool           `json:"disabled"`
	SceneName      string         `json:"scene_name"`
	Rate           float32        `json:"rate,omitempty"`
	RpsRate        int64          `json:"rps_rate,omitempty"`
	Concurrency    int64          `json:"concurrency,omitempty"`
	//CaseTree       []*SceneCaseTree `json:"case_tree"`
	//VariableList   []*Variable      `json:"variable_list"`
}

type ExportDataInfo struct {
	VariableList []string `json:"variable_list"`
	ExportTimes  int32    `json:"export_times"`
	Concurrency  int32    `json:"concurrency"`
	//ExportList   []Variable // 前端不用，用于scene间的数据传递
}

type Scene struct {
	SceneID        int32          `json:"scene_id"`         // 场景id
	SceneName      string         `json:"scene_name"`       // 场景名
	SceneType      int32          `json:"scene_type"`       // 场景类型
	PlanID         int32          `json:"plan_id"`          // 计划id
	ExportDataInfo ExportDataInfo `json:"export_data_info"` // 导出数据信息
	Disabled       bool           `json:"disabled"`         // 是否禁用
	CreateUserID   int32          `json:"create_user_id"`   // 创建人id
	CreateUserName string         `json:"create_user_name"` // 创建人
	UpdateUserID   int32          `json:"update_user_id"`   // 最后修改人id
	UpdateUserName string         `json:"update_user_name"` // 最后修改人
	CreateTime     string         `json:"create_time"`      // 创建时间
	UpdateTime     string         `json:"update_time"`      // 最后修改时间
	//VariablePool   VariablePool         `json:"variable_pool"`
	//Cases          []*SceneCaseTree     `json:"cases"`            // 场景case
	//Result         SceneExecutionResult `json:"execution_result"` // 场景运行结果，仅供scene调试执行
	//DefaultHeader  []*ParamsForm        `json:"default_header"`
}

type SceneInformation struct {
	SceneId   int32                     `json:"scene_id"`
	SceneName string                    `json:"scene_name"`
	SceneRps  int64                     `json:"scene_rps"`
	Cases     []*ReportStatusCaseEntity `json:"cases"`
}
