package rao

const (
	PreRunning = iota + 1
	Running
	PreDown
	Down
)

type ReportListReq struct {
	PlanID       int32 `json:"plan_id"`                           // 计划id，非必填
	CreateUserId int32 `json:"create_user_id"`                    // 创建者id，非必填
	Page         int   `json:"page" binding:"required,gt=0"`      // 页码
	PageSize     int   `json:"page_size" binding:"required,gt=0"` // 每页条数
}

type Report struct {
	ReportID       int32   `json:"report_id"`       // 报告id
	ReportName     string  `json:"report_name"`     // 报告名
	Status         bool    `json:"status"`          // 压测状态
	PlanID         int32   `json:"plan_id"`         // 计划id
	PlanName       string  `json:"plan_name"`       // 计划名
	PressType      int32   `json:"press_type"`      // 压测模式
	Duration       int32   `json:"duration"`        // 预计持续时间
	ActualDuration float64 `json:"actual_duration"` // 实际持续时间
	Concurrency    int64   `json:"concurrency"`     // 并发数
	MachineNum     int32   `json:"machine_num"`     // 机器数量
	Remark         string  `json:"remark"`

	TotalCount  int64   `json:"total_count"`  // 总请求数
	SuccessRate float64 `json:"success_rate"` // 成功率
	StartTime   string  `json:"start_time"`   // 开始时间
	EndTime     string  `json:"end_time"`     // 结束时间

	TotalSendBytes    float64 `json:"total_send_bytes"`    // 发送流量
	TotalReceiveBytes float64 `json:"total_receive_bytes"` // 接受流量

	//MaxRps               float64            `json:"max_rps"`                 // rps
	//MaxConcurrent        int64              `json:"max_concurrent"`          // 最大并发
	ScenePressResultList []ScenePressResult `json:"scene_press_result_list"` // 场景压测结果
	CreateUserName       string             `json:"create_user_name"`        // 创建人
	UpdateUserName       string             `json:"update_user_name"`        // 最后修改人
	UpdateTime           string             `json:"update_time"`             // 最后修改时间
}

type SamplingDataReq struct {
	ReportID   int32   `json:"report_id" bson:"report_id"`
	SceneID    int32   `json:"scene_id" bson:"scene_id"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	CaseIdList []int32 `json:"case_id_list"`
	CodeList   []int32 `json:"code_list"`
	StatusList []int32 `json:"status_list"`
	MinRt      int32   `json:"min_rt"`
	MaxRt      int32   `json:"max_rt"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	IsSuccess  bool    `json:"is_success" bson:"is_success"`
}

type SamplingDataResp struct {
	Total int64          `json:"total"`
	List  []SamplingData `json:"list"`
}

type DoReportRpsModifyReq struct {
	ReportId int32 `json:"report_id"`
	CaseId   int32 `json:"case_Id"`
	NewValue int64 `json:"new_rps"`
}

type SamplingData struct {
	Success          bool         `json:"success"`
	StartTime        string       `json:"start_time"`
	CaseName         string       `json:"case_name"`
	CaseID           int32        `json:"case_id"`
	ResponseTime     int64        `json:"response_time"`
	StatusCode       int          `json:"status_code"`
	HttpResponseData HttpResponse `json:"http_response_data"`
}

// ReportData report服务返回的data
type ReportData struct {
	End bool `json:"end"`
	// stage数据
	StageRequestTime int64 `json:"stage_request_time"`
	StageRequestNum  int64 `json:"stage_request_num"` // 步骤请求总数
	StageSuccessNum  int64 `json:"stage_success_num"` // 步骤请求成功数
	StageErrorNum    int64 `json:"stage_error_num"`   // 步骤请求错误数

	Concurrency      int64   `json:"concurrency"`
	StageRps         float64 `json:"stage_rps"`
	StageRT          float64 `json:"stage_rt"`
	StageSuccessRate float64 `json:"stage_success_rate"`
	StageStartTime   int64   `json:"stage_start_time"` // 请求开始时间
	StageEndTime     int64   `json:"stage_end_time"`   // 请求结束时间

	// total数据
	TotalRequestTime int64   `json:"total_request_time"` // 总请求时长
	TotalRequestNum  int64   `json:"total_request_num"`  // 总请求数
	TotalSuccessNum  int64   `json:"total_success_num"`  // 总请求成功数
	TotalErrorNum    int64   `json:"total_error_num"`    // 总请求失败数
	TotalSuccessRate float64 `json:"total_success_rate"` // 接口成功率

	TotalSendBytes     float64 `json:"total_send_bytes"`
	TotalReceivedBytes float64 `json:"total_received_bytes"`
	TotalStartTime     int64   `json:"total_start_time"` // 请求开始时间
	TotalEndTime       int64   `json:"total_end_time"`   // 请求结束时间

	//Graph *ReportDataGraphEntity `json:"graph"`
	//
	Scenes []*SceneResultData `json:"scenes"`
}

type SceneResultData struct {
	SceneID     int32             `json:"scene_id"`
	SceneType   int32             `json:"scene_type"`
	Concurrency int64             `json:"concurrency"`
	Cases       []*CaseResultData `json:"cases"`
}

type CaseResultData struct {
	CaseID int32 `json:"case_id"` // 用例ID
	// stage数据
	StageRequestTime int64   `json:"stage_request_time"`
	StageRequestNum  int64   `json:"stage_request_num"` // 步骤请求总数
	StageSuccessNum  int64   `json:"stage_success_num"` // 步骤请求成功数
	StageErrorNum    int64   `json:"stage_error_num"`   // 步骤请求错误数
	StageRps         float64 `json:"stage_rps"`
	StageSuccessRate float64 `json:"stage_success_rate"` // 每秒接口成功率
	StageAvgRt       float64 `json:"stage_avg_rt"`       // 平均响应时间
	StageStartTime   int64   `json:"stage_start_time"`   // 请求开始时间
	StageEndTime     int64   `json:"stage_end_time"`     // 请求结束时间

	MaxRt                          float64 `json:"max_rt"`                              // 最大响应时间
	MinRt                          float64 `json:"min_rt"`                              // 最小响应时间
	FiftyRequestTimeLineValue      float64 `json:"fifty_request_time_line_value"`       // 50线
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`      // 90线
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"` // 95线
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"` // 99线

	// total数据
	TotalRequestTime int64   `json:"total_request_time"` // 总请求时长
	TotalRequestNum  int64   `json:"total_request_num"`  // 总请求数
	TotalSuccessNum  int64   `json:"total_success_num"`  // 总请求成功数
	TotalErrorNum    int64   `json:"total_error_num"`    // 总请求失败数
	TotalSuccessRate float64 `json:"total_success_rate"` // 接口成功率
	TwoXxCodeNum     int64   `json:"two_xx_code_num"`    // 2xx数
	OtherCodeNum     int64   `json:"other_code_num"`     // 其他响应码数
	TotalRps         float64 `json:"total_rps"`          // 平均RPS
	TotalAvgRt       float64 `json:"total_avg_rt"`       // 平均RT
	SendBytes        float64 `json:"send_bytes"`         // 发送流量
	ReceivedBytes    float64 `json:"received_bytes"`     // 接收流量
	TotalStartTime   int64   `json:"total_start_time"`   // 请求开始时间
	TotalEndTime     int64   `json:"total_end_time"`     // 请求结束时间

	TargetRps float64 `json:"target_rps"` // 目标Rps
}

type CreateReportReq struct {
	PlanID    int32  `json:"plan_id"`
	PlanName  string `json:"plan_name"`
	Duration  int32  `json:"duration"`
	PressType int32  `json:"press_type"`
	//MachineNum   int32  `json:"machine_num"`
	EngineList   []string `json:"engine_list"`
	Rps          int32    `json:"rps"`
	CreateUserID int32    `json:"create_user_id"`
}

type CapacityListReq struct {
	Page      int    `json:"page" binding:"required,gt=0"`      // 页码
	PageSize  int    `json:"page_size" binding:"required,gt=0"` // 每页条数
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	TargetRps int64  `json:"target_rps"`
}

type CapacityListResp struct {
	ID       int32            `json:"id"`
	ReportID int32            `json:"report_id"` // 测试报告id
	Rps      float64          `json:"rps"`       // RPS
	Result   map[string]int32 `json:"result"`
}

type CapacityHistory struct {
}
