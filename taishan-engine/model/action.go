package model

type Action struct {
	Plan               Plan              `json:"plan"`
	ReportID           int32             `json:"report_id"`
	ReportName         string            `json:"report_name"`
	ScenePressInfoList []*ScenePressInfo `json:"scene_press_info_list"`
	FileInfo           []*FileInfo       `json:"file_info"`
	EngineSerialNumber int32             `json:"engine_serial_number"` // 压测引擎序号
	EngineCount        int32             `json:"engine_count"`         // 压测引擎的总数量
	PartitionID        int32             `json:"partition_id"`         // 压测引擎的分区id
}

type SceneAction struct {
	ReportID           int32           `json:"report_id"`
	PlanID             int32           `json:"plan_id"`
	SceneID            int32           `json:"scene_id"`
	ScenePressInfo     *ScenePressInfo `json:"scene_list"`
	PressInfo          PressInfo       `json:"press_info"`
	SamplingInfo       SamplingInfo    `json:"sampling_info"`
	BreakType          int32           `json:"break_type"`
	BreakValue         float32         `json:"break_value"`
	PartitionID        int32           `json:"partition_id"`
	EngineSerialNumber int32           `json:"engine_serial_number"` // 压测引擎序号
	EngineCount        int32           `json:"engine_count"`         // 压测引擎的总数量
}

type SceneCases []*SceneCaseTree
