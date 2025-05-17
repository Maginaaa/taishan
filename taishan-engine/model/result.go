package model

const (
	TypePlanAction = iota
	_              // plan result预留类型
	TypeSceneAction
	TypeSceneResultData
	TypeHttpResultData
)

type ResultDataMsg struct {
	MsgType           int32              `json:"msg_type" bson:"msg_type"`
	GlobalInfo        *GlobalInfo        `json:"global_info"`
	Start             bool               `json:"start" bson:"start"`
	End               bool               `json:"end" bson:"end"` //最后一条记录标记
	SceneInfo         *SceneInfo         `json:"scene_info"`
	HttpResultDataMsg *HttpResultDataMsg `json:"http_result_data_msg" bson:"http_result_data_msg"`
}

type GlobalInfo struct {
	PartitionID  int32  `json:"partition_id" bson:"partition_id"`
	ReportID     int32  `json:"report_id" bson:"report_id"`
	PlanID       int32  `json:"plan_id" bson:"plan_id"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp"`
	MachineIP    string `json:"machine_ip" bson:"machine_ip"`
	NeedSampling bool
}

type SceneInfo struct {
	SceneID     int32 `json:"scene_id" bson:"scene_id"`
	SceneType   int32 `json:"scene_type" bson:"scene_type"`
	Concurrency int64 `json:"concurrency" bson:"concurrency"`
}

type HttpResultDataMsg struct {
	CaseID            int32  `json:"case_id" bson:"case_id"`
	CaseName          string `json:"case_name" bson:"case_name"`
	Url               string `json:"url" bson:"url"`
	ActualConcurrency int64  `json:"actual_concurrency" bson:"actual_concurrency"` // 记录当前的目标并发数

	StatusCode    int     `json:"status_code" bson:"status_code"`
	RequestTime   int64   `json:"request_time" bson:"request_time"`     // 请求响应时间
	ErrorType     int64   `json:"error_type" bson:"error_type"`         // 错误类型：1. 请求错误；2. 断言错误
	IsSuccess     bool    `json:"is_success" bson:"is_success"`         // 请求是否有错：true / false   为了计数
	SendBytes     float64 `json:"send_bytes" bson:"send_bytes"`         // 发送字节数 单位KB
	ReceivedBytes float64 `json:"received_bytes" bson:"received_bytes"` // 接收字节数 单位KB
	StartTime     int64   `json:"start_time" bson:"start_time"`
	EndTime       int64   `json:"end_time" bson:"end_time"`
	RallyPoint    int64   `json:"rally_point"` // 集合点个数

	DebugInfo   *HttpResponse `json:"debug_info,omitempty" bson:"debug_info"`
	GoroutineID int64         `json:"goroutine_id" bson:"goroutine_id"`
}
