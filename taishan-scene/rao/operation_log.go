package rao

const (
	CreateOperation = iota
	UpdateOperation
	CopyOperation
	DeleteOperation
	_
	_
	_
	_
)

const (
	DisabledOperation = iota + 11
	EnableOperation
)

const (
	DebugOperation = iota + 21
	RunOperation
)

const (
	SourceSceneCase        = "SceneCase"
	SourceScene            = "Scene"
	SourcePlan             = "Plan"
	SourceFile             = "File"
	SourceNormalizedSwitch = "NormalizedSwitch"
)

type OperationLogReq struct {
	SourceName string `json:"source_name"`
	SourceID   int32  `json:"source_id"`
	Page       int    `json:"page"`      // 页码
	PageSize   int    `json:"page_size"` // 每页条数
}

type OperationLogResp struct {
	ID            int32                           `json:"id"`
	SourceId      int32                           `json:"source_id"`
	SourceName    string                          `json:"source_name"`
	OperationType int32                           `json:"operation_type"`
	OperatorName  string                          `json:"operator_name"`
	ValueDiff     map[string]OperationDetail[any] `json:"value_diff"`
	CreateTime    string                          `json:"create_time"`
}

type OperationLog struct {
	SourceName    string      `json:"source_name"`
	SourceID      int32       `json:"source_id"`
	OperationType int32       `json:"operation_type"`
	OperatorID    int32       `json:"operator_id"`
	ValueBefore   interface{} `json:"value_before"`
	ValueAfter    interface{} `json:"value_after"`
	ValueDiff     string      `json:"value_diff"`
}

type OperationDetail[T any] struct {
	Before T `json:"before"`
	After  T `json:"after"`
}
