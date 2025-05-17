package rao

import (
	"encoding/json"
	"scene/internal/biz/log"
	"time"
)

type PlanTestResultData struct {
	PlanResult   []PlanTestResultCalcData           `json:"plan_result"`
	ScenesResult map[int32][]SceneTestResultDataMsg `json:"scenes_result"`
}

// PlanTestResultCalcData 计划的单次(stage)测试结果
type PlanTestResultCalcData struct {
	End                    bool    `json:"end"`
	Time                   string  `json:"time"`
	StageTotalRequestNum   int64   `json:"stage_total_request_num"`
	StageTotalResponseTime int64   `json:"stage_total_response_time"`
	StageAvgResponseTime   float64 `json:"stage_avg_response_time"`
	ActualConcurrency      int64   `json:"actual_concurrency"`
	StageSendBytes         float64 `json:"stage_send_bytes"`     // 发送字节数
	StageReceivedBytes     float64 `json:"stage_received_bytes"` // 接收字节数

	TotalResponseTime int64   `json:"total_response_time"`
	TotalRequestNum   int64   `json:"total_request_num"`
	TotalSuccessNum   int64   `json:"total_success_num"`
	ApiSuccessRate    float64 `json:"api_success_rate"` // total数据
	PlanTps           float64 `json:"plan_tps"`         // total数据
}

// SceneTestResultDataMsg 单次场景的测试结果
type SceneTestResultDataMsg struct {
	End               bool    `json:"end"`
	ReportID          int32   `json:"report_id"`
	PlanID            int32   `json:"plan_id"`            // 任务ID
	SceneID           int32   `json:"scene_id"`           // 场景
	SceneType         int32   `json:"scene_type"`         // 场景类型： 0:普通场景  1:前置场景
	TargetConcurrency int64   `json:"target_concurrency"` // 当前接口实际并发(通过协程数获取)
	ActualConcurrency int64   `json:"actual_concurrency"` // 记录当前的目标并发数
	StageSceneTps     float64 `json:"stage_scene_tps"`
	StageSceneRps     float64 `json:"stage_scene_rps"`
	Stage             int64   `json:"stage"` // 采集次数(总请求时间 / 采集频率 + 1)

	SceneSuccessRate         float64 `json:"scene_success_rate"` // 场景成功率
	TotalSceneRequestNum     int64   `json:"total_scene_request_num"`
	TotalSceneRequestCaseNum int64   `json:"total_scene_request_case_num"`
	TotalSceneRequestTime    int64   `json:"total_scene_request_time"`
	MinSceneRequestTime      float64 `json:"min_scene_request_time"`
	MaxSceneRequestTime      float64 `json:"max_scene_request_time"`
	AvgSceneRequestTime      float64 `json:"avg_scene_request_time"`
	TotalSceneSuccessNum     int64   `json:"total_scene_success_num"`
	TotalTps                 float64 `json:"total_tps"`

	TotalConcurrency int64 `json:"total_concurrency"` // 实际并发数做累加，用于计算vum

	TotalSendBytes     float64 `json:"total_send_bytes"`     // 发送字节数
	TotalReceivedBytes float64 `json:"total_received_bytes"` // 接收字节数

	StageSceneRequestNum  int64 `json:"stage_scene_request_num"` // 场景执行次数
	StageSceneRequestTime int64 `json:"stage_scene_request_time"`
	//StageTotalSceneRequestCaseNum int64 `json:"stage_total_scene_request_case_num"` // 场景内接口的执行次数(未完全执行完的场景不会计入该统计)，区别于StageTotalCaseRequestNum，数据会有滞后性

	StageTotalSceneSuccessNum int64 `json:"stage_total_scene_success_num"`          // 成功的场景数
	ParameterCount            int64 `json:"parameter_count" bson:"parameter_count"` // 参数个数

	TotalApiSuccessNum int64   `json:"total_api_success_num"` // 接口成功数
	ApiSuccessRate     float64 `json:"api_success_rate"`      // 请求成功率(接口级别)

	StageTotalCaseRequestNum  int64   `json:"stage_total_case_request_num"`  // 步骤内case总请求次数
	StageTotalCaseRequestTime int64   `json:"stage_total_case_request_time"` // 步骤内case总请求时间
	StageTotalCaseSuccessNum  int64   `json:"stage_total_case_success_num"`  // 步骤内case成功请求次数
	StageAvgCaseResponseTime  float64 `json:"stage_avg_case_response_time"`  // 步骤内case的平均响应时间

	Results   map[int32]*ApiTestResultDataMsg `json:"results"`
	TimeStamp int64                           `json:"time_stamp"` // 单位为秒
	Time      string                          `json:"time"`       // mm-dd hh:mm:ss格式时间

	StartTime int64 `json:"start_time"` // stage的开始时间
	EndTime   int64 `json:"end_time"`   // stage的结束时间
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
type ApiTestResultDataMsg struct {
	// 以下为固定配置数据
	CaseID     int32  `json:"case_id"`
	CaseName   string `json:"case_name"`
	RallyPoint int64  `json:"rally_point"` // 集合点个数

	// 以下都是汇总数据total
	TotalRequestNum  int64 `json:"total_request_num"`  // 总请求数
	TotalRequestTime int64 `json:"total_request_time"` // 总响应时间
	SuccessNum       int64 `json:"success_num"`
	ErrorNum         int64 `json:"error_num"` // 错误数

	AvgRequestTime                 float64 `json:"avg_request_time"` // 平均响应事件
	MaxRequestTime                 float64 `json:"max_request_time"`
	MinRequestTime                 float64 `json:"min_request_time"` // 毫秒
	TargetRps                      int64   `json:"target_rps"`       // 目标rps
	ConRps                         float64 `json:"con_rps"`          //当前rps
	CustomRequestTimeLine          int64   `json:"custom_request_time_line"`
	FiftyRequestTimeline           int64   `json:"fifty_request_time_line"`
	NinetyRequestTimeLine          int64   `json:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine      int64   `json:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine      int64   `json:"ninety_nine_request_time_line"`
	CustomRequestTimeLineValue     float64 `json:"custom_request_time_line_value"`
	FiftyRequestTimelineValue      float64 `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value"`
	TotalSendBytes                 float64 `json:"total_send_bytes"`     // 发送字节数
	TotalReceivedBytes             float64 `json:"total_received_bytes"` // 接收字节数
	StartTime                      int64   `json:"start_time"`
	EndTime                        int64   `json:"end_time"`
	StatusCodeCounter              []int64 `json:"status_code_counter"`

	Rps  float64 `json:"rps"`
	SRps float64 `json:"s_rps"` // success_rps
	ERps float64 `json:"e_rps"` // error_reps

	// 以下都是当前步骤数据
	ActualConcurrency int64 `json:"actual_concurrency" bson:"actual_concurrency"`

	StageStartTime        int64   `json:"stage_start_time"` // 某段时间内，开始时间
	StageEndTime          int64   `json:"stage_end_time"`
	StageTotalRequestNum  int64   `json:"stage_total_request_num"` // 某段时间内的，总请求数
	StageTotalRequestTime int64   `json:"stage_total_request_time"`
	StageAvgResponseTime  float64 `json:"stage_avg_response_time"`
	StageSuccessNum       int64   `json:"stage_success_num"`
	StageErrorNum         int64   `json:"stage_error_num"`
	StageSendBytes        float64 `json:"stage_send_bytes"`
	StageReceivedBytes    float64 `json:"stage_received_bytes"`

	StageRps  float64 `json:"stage_rps"`
	StageSRps float64 `json:"stage_s_rps"` // success_rps
	StageERps float64 `json:"stage_e_rps"` // error_reps

}

func (s *SceneTestResultDataMsg) ToJson() string {
	res, err := json.Marshal(s)
	if err != nil {
		log.Logger.Error("c_model.SceneTestResultDataMsg.jsonMarshal()： ", err)
		return ""
	}
	return string(res)
}

func (s *SceneTestResultDataMsg) FormatTime() {
	s.Time = time.Unix(s.TimeStamp, 0).Format("01-02 15:04:05")
}
