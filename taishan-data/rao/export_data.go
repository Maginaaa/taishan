package rao

type ExportDataInfo struct {
	MachineIPSet map[string]struct{}
	TitleInit    bool
	TitleArray   []string
	Content      [][]string
}

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

type HttpResponse struct {
	// 基础信息
	CaseId   int32  `json:"case_id" bson:"case_id"`
	CaseName string `json:"case_name" bson:"case_name"`
	SceneId  int32  `json:"scene_id" bson:"scene_id"`

	// 请求信息
	URL            string        `json:"url" bson:"url"`
	MethodType     string        `json:"method_type" bson:"method_type"`
	RequestHeader  []*ParamsForm `json:"request_header" bson:"request_header"`
	RequestContent string        `json:"request_content" bson:"request_content"`

	// 响应信息
	StatusCode               int             `json:"status_code" bson:"status_code"`
	ResponseContent          string          `json:"response_content" bson:"response_content"`
	ResponseHeader           []*ParamsForm   `json:"response_header" bson:"response_header"`
	AssertRes                []*AssertItem   `json:"assert_res" bson:"assert_res"`
	VariableRes              []*VariableItem `json:"variable_res" bson:"variable_res"`
	ResponseContentInterface interface{}     `json:"response_content_interface" bson:"response_content_interface"`

	// 统计信息
	AssertSuccess     bool         `json:"assert_success" bson:"assert_success"`
	ExtractAllSuccess bool         `json:"extract_all_success" bson:"extract_all_success"`
	ResponseSize      ResponseSize `json:"response_size" bson:"response_size"`
	ResponseTime      int64        `json:"response_time" bson:"response_time"`
	SendBytes         float64      `json:"send_bytes" bson:"send_bytes"`
	ReceiverBytes     float64      `json:"receiver_bytes" bson:"receiver_bytes"`
	RequestSuccess    bool         `json:"request_success" bson:"request_success"`
	Err               string       `json:"err" bson:"err"`
	StartTime         int64        `json:"start_time" bson:"start_time"`
	EndTime           int64        `json:"end_time" bson:"end_time"`
}

type HeadersForm struct {
	Desc   string `json:"desc" bson:"desc"`
	Enable bool   `json:"enable" bson:"enable"`
	Key    string `json:"key" bson:"key"`
	Value  string `json:"value" bson:"value"`
}

type ParamsForm struct {
	Desc   string `json:"desc" bson:"desc"`
	Enable bool   `json:"enable" bson:"enable"`
	Key    string `json:"key" bson:"key"`
	Value  string `json:"value" bson:"value"`
}

type AssertForm struct {
	CheckingRule   int    `json:"checking_rule" bson:"checking_rule"`
	Desc           string `json:"desc" bson:"desc"`
	Enable         bool   `json:"enable" bson:"enable"`
	ExpectValue    string `json:"expect_value" bson:"expect_value"`
	ExtractExpress string `json:"extract_express" bson:"extract_express"` // 提取表达式
	ExtractType    int    `json:"extract_type" bson:"extract_type"`
}

type VariableForm struct {
	Enable      bool   `json:"enable" bson:"enable"`
	ExtractType int    `json:"extract_type" bson:"extract_type"`
	Key         string `json:"key" bson:"key"`
	Value       string `json:"value" bson:"value"`
	Desc        string `json:"desc,omitempty" bson:"desc"`
}

type ResponseSize struct {
	BodySize   float64 `json:"body_size" bson:"body_size"`     // 单位KB
	HeaderSize float64 `json:"header_size" bson:"header_size"` // 单位KB
	TotalSize  float64 `json:"total_size" bson:"total_size"`   // 单位KB
}

type VariableItem struct {
	ExtractType    int    `json:"extract_method" bson:"extract_type"`
	ExtractRule    string `json:"extract_rule" bson:"extract_rule"`
	VariableName   string `json:"variable_name" bson:"variable_name"`
	ActualRes      string `json:"actual_res" bson:"actual_res"`
	ExtractSuccess bool   `json:"extract_success" bson:"extract_success"`
}

type AssertItem struct {
	ExtractMethod string `json:"extract_method" bson:"extract_method"`
	ExtractRule   string `json:"extract_rule" bson:"extract_rule"`
	ExpectedValue string `json:"expected_value" bson:"expected_value"` // 预期值
	ExtractValue  string `json:"extract_value" bson:"extract_value"`   // 实际提取值
	CheckingRule  int    `json:"checking_rule" bson:"checking_rule"`
	AssertPass    bool   `json:"assert_pass" bson:"assert_pass"` // 断言结果
}
