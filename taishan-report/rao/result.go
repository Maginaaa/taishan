package rao

type PlanResultData struct {
	End                bool   `json:"end"`
	PrefixSceneEndTime string `json:"prefix_scene_end_time"` // 前置场景是否在执行

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

	TargetRps int64 `json:"target_rps"`

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

	// 错误信息分布
	ErrorUrlArray     ApiDistributionList `json:"error_url_array"`
	ErrorCodeArray    ApiDistributionList `json:"error_code_array"`
	RequestErrorArray ApiDistributionList `json:"request_error_array"`
	AssertErrorArray  ApiDistributionList `json:"assert_error_array"`

	Graph *ReportDataGraphEntity `json:"graph"`

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

	//MaxRt                          float64 `json:"max_rt"`                              // 最大响应时间
	//MinRt                          float64 `json:"min_rt"`                              // 最小响应时间
	FiftyRequestTimeLineValue      float64 `json:"fifty_request_time_line_value"`       // 50线
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`      // 90线
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"` // 95线
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"` // 99线

	// total数据
	TotalRequestTime     int64   `json:"total_request_time"`      // 总请求时长
	TotalRequestNum      int64   `json:"total_request_num"`       // 总请求数
	TotalSuccessNum      int64   `json:"total_success_num"`       // 总请求成功数
	TotalRequestErrorNum int64   `json:"total_request_error_num"` // 总请求失败数
	TotalAssertErrorNum  int64   `json:"total_assert_error_num"`  // 总断言失败数
	TotalErrorNum        int64   `json:"total_error_num"`         // 总请求失败数
	TotalSuccessRate     float64 `json:"total_success_rate"`      // 接口成功率
	TwoXxCodeNum         int64   `json:"two_xx_code_num"`         // 2xx数
	ThreeXxCodeNum       int64   `json:"three_xx_code_num"`       // 3xx数
	FourXxCodeNum        int64   `json:"four_xx_code_num"`        // 4xx数
	FiveXxCodeNum        int64   `json:"five_xx_code_num"`        // 5xx数
	OtherCodeNum         int64   `json:"other_code_num"`          // 其他响应码数
	TotalRps             float64 `json:"total_rps"`               // 平均RPS
	TotalAvgRt           float64 `json:"total_avg_rt"`            // 平均RT
	SendBytes            float64 `json:"send_bytes"`              // 发送流量
	ReceivedBytes        float64 `json:"received_bytes"`          // 接收流量
	TotalStartTime       int64   `json:"total_start_time"`        // 请求开始时间
	TotalEndTime         int64   `json:"total_end_time"`          // 请求结束时间

	TargetRps float64 `json:"target_rps"` // 目标Rps
}

type ReportRecordEntity struct {
	CaseId    int32 `json:"case_id"`
	SceneId   int32 `json:"scene_id"`
	SceneType int32 `json:"scene_type"`
	CollectorCaseData
}

type ReportRecordBaseEntity struct {
	CaseId    string `json:"case_id"`
	SceneId   string `json:"scene_id"`
	SceneType string `json:"scene_type"`
	CollectorCaseData
}

type CollectorCaseData struct {
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
	RequestNum                     int64   `json:"request_num"`
	RequestTime                    int64   `json:"request_time"`
	SuccessNum                     int64   `json:"success_num"`
	RequestErrorNum                int64   `json:"request_error_num"`
	AssertErrorNum                 int64   `json:"assert_error_num"`
	ErrorNum                       int64   `json:"error_num"`
	SendBytes                      float64 `json:"send_bytes"`
	ReceivedBytes                  float64 `json:"received_bytes"`
	Url                            string  `json:"url"`
}

// ReportStatusCaseEntity 用例状态实体
type ReportStatusCaseEntity struct {
	CaseId                         int32   `json:"case_id"`
	CaseName                       string  `json:"case_name"`
	Time                           int64   `json:"time"`
	StartTime                      int64   `json:"start_time"`
	EndTime                        int64   `json:"end_time"`
	ActualConcurrency              int64   `json:"actual_concurrency"`
	FiftyRequestTimeLineValue      float64 `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"`
	OneXxCodeNum                   int64   `json:"one_xx_code_num"`
	TwoXxCodeNum                   int64   `json:"tow_xx_code_num"`
	ThreeXxCodeNum                 int64   `json:"three_xx_code_num"`
	FourXxCodeNum                  int64   `json:"four_xx_code_num"`
	FiveXxCodeNum                  int64   `json:"five_xx_code_num"`
	RequestNum                     int64   `json:"request_num"`
	RequestTime                    int64   `json:"request_time"`
	SuccessNum                     int64   `json:"success_num"`
	ErrorNum                       int64   `json:"error_num"`
	SendBytes                      float64 `json:"send_bytes"`
	ReceivedBytes                  float64 `json:"received_bytes"`
	CurrentRps                     float64 `json:"current_rps"`
}

// ReportStatusSceneEntity 场景状态实体
type ReportStatusSceneEntity struct {
	SceneId   int32                    `json:"scene_id"`
	SceneName string                   `json:"scene_name"`
	Cases     []ReportStatusCaseEntity `json:"cases"`
}

// ReportRpsResponse 报告状态数据
type ReportRpsResponse struct {
	ReportID int32                     `json:"report_id"`
	Scenes   []ReportStatusSceneEntity `json:"scenes"`
}

// ReportDataGraphEntity 图表数据
type ReportDataGraphEntity struct {
	Rps         [][3]any `json:"rps"`
	RequestTime [][2]any `json:"request_time"`
	SuccessRate [][2]any `json:"success_rate"`
}

type ApiDistributionList []ApiDistribution
type ApiDistribution struct {
	CaseID     int32  `json:"case_id"`
	Url        string `json:"url"`
	ErrorCode  string `json:"error_code,omitempty"`
	Count      int64  `json:"count"`
	ErrorCount int64  `json:"error_count"`
}

func (rt ApiDistributionList) Len() int {
	return len(rt)
}

// Less 页面需要错误数从大到小排序
func (rt ApiDistributionList) Less(i int, j int) bool {
	return rt[i].Count > rt[j].Count
}
func (rt ApiDistributionList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}
