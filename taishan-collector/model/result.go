package model

type PlanTestResultDataMsg struct {
	ReportID  int32  `json:"report_id"`
	PlanID    int32  `json:"plan_id"`
	MachineIP string `json:"machine_ip"`
	End       bool   `json:"end"`

	SceneResults map[int32]*SceneResultDataMsg `json:"scene_results"`
	TimeStamp    int64                         `json:"time_stamp"` // 单位为毫秒
}

type SceneResultDataMsg struct {
	TargetConcurrency int64                           `json:"target_concurrency" bson:"target_concurrency"` // 记录当前的目标并发数
	SceneType         int32                           `json:"scene_type" bson:"scene_type"`
	End               bool                            `json:"end" bson:"end"`
	CaseResults       map[int32]*ApiTestResultDataMsg `json:"case_results"`
}

type ApiTestResultDataMsg struct {
	// 以下为固定配置数据
	CaseID     int32  `json:"case_id"`
	Url        string `json:"url"`
	RallyPoint int64  `json:"rally_point"` // 集合点个数

	RequestNum  int64 `json:"request_num"`  // 总请求数
	RequestTime int64 `json:"request_time"` // 总响应时间
	SuccessNum  int64 `json:"success_num"`
	ErrorNum    int64 `json:"error_num"` // 错误数

	RequestErrorNum int64 `json:"request_error_num"`
	AssertErrorNum  int64 `json:"assert_error_num"`

	MaxRequestTime                 float64       `json:"max_request_time"`
	MinRequestTime                 float64       `json:"min_request_time"` // 毫秒
	FiftyRequestTimeLineValue      float64       `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     float64       `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64       `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64       `json:"ninety_nine_request_time_line_value"`
	SendBytes                      float64       `json:"send_bytes"`     // 发送字节数
	ReceivedBytes                  float64       `json:"received_bytes"` // 接收字节数
	StartTime                      int64         `json:"start_time"`
	EndTime                        int64         `json:"end_time"`
	StatusCodeCounter              map[int]int64 `json:"status_code_counter"`

	ActualConcurrency int64 `json:"actual_concurrency" bson:"actual_concurrency"`
}

type StageCaseResult struct {
	CaseID    int32     `json:"case_id"`
	SceneID   int32     `json:"scene_id"`
	SceneType int32     `json:"scene_type"`
	CalcData  *CalcData // key为机器ip
	BaseData  *BaseData
}

// CalcData : key为机器ip,不把CalcData改为map[string]CalcData的原因是：每次消费不一定能采集到所有机器的所有case信息，将ip放在CalcData内部，能保障只有采集到的数据才会在parseBatchMachineData时计为分母
type CalcData struct {
	FiftyRequestTimeMap      map[string]float64
	NinetyRequestTimeMap     map[string]float64
	NinetyFiveRequestTimeMap map[string]float64
	NinetyNineRequestTimeMap map[string]float64
	//DurationMap map[string][][2]int64 //   key为机器ip,[[duration1. reqNum1], [duration2. reqNum2]]

	ActualConcurrencyMap map[string]int64
}

type BaseData struct {
	MaxRequestTime                 float64 `json:"max_request_time"`
	MinRequestTime                 float64 `json:"min_request_time"` // 毫秒
	FiftyRequestTimeLineValue      float64 `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"`
	Normal1xxCodeNum               int64   `json:"1xx_code_num"`
	Normal2xxCodeNum               int64   `json:"2xx_code_num"`
	Normal3xxCodeNum               int64   `json:"3xx_code_num"`
	Normal4xxCodeNum               int64   `json:"4xx_code_num"`
	Normal5xxCodeNum               int64   `json:"5xx_code_num"`
	ActualConcurrency              int64   `json:"actual_concurrency"`
	StartTime                      int64   `json:"start_time"`
	EndTime                        int64   `json:"end_time"`
	//Duration                       int64   `json:"duration"`
	RequestNum      int64   `json:"request_num"`
	RequestTime     int64   `json:"request_time"`
	SuccessNum      int64   `json:"success_num"`
	RequestErrorNum int64   `json:"request_error_num"`
	AssertErrorNum  int64   `json:"assert_error_num"`
	ErrorNum        int64   `json:"error_num"`
	SendBytes       float64 `json:"send_bytes"`
	ReceivedBytes   float64 `json:"received_bytes"`
	Url             string  `json:"url"`
}
