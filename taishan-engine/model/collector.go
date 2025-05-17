package model

import "encoding/json"

type CaseTestResultDataMsg struct {
	CaseResults map[int32]*ApiTestResultDataMsg `json:"case_results"`
}

// 仅记录当前stage结果，不记录总数据
type PlanTestResultDataMsg struct {
	ReportID  int32  `json:"report_id"`
	PlanID    int32  `json:"plan_id"`
	MachineIP string `json:"machine_ip"`
	End       bool   `json:"end"`

	SceneResults map[int32]*SceneResultDataMsg `json:"scene_results"`
	TimeStamp    int64                         `json:"time_stamp"` // 单位为毫秒
}

func (p *PlanTestResultDataMsg) ToByte() (msg []byte) {
	msg, _ = json.Marshal(p)
	return
}

type SceneResultDataMsg struct {
	TargetConcurrency int64                           `json:"target_concurrency" bson:"target_concurrency"` // 记录当前的目标并发数
	SceneType         int32                           `json:"scene_type" bson:"scene_type"`
	End               bool                            `json:"end" bson:"end"`
	CaseResults       map[int32]*ApiTestResultDataMsg `json:"case_results"`
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
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

type RequestTimeList []int64

func (rt RequestTimeList) Len() int {
	return len(rt)
}

func (rt RequestTimeList) Less(i int, j int) bool {
	return rt[i] < rt[j]
}
func (rt RequestTimeList) Swap(i int, j int) {
	rt[i], rt[j] = rt[j], rt[i]
}

// TimeLineCalculate 根据响应时间线，计算该线的值
func TimeLineCalculate(line int64, requestTimeList RequestTimeList) (requestTime float64) {
	if line > 0 && line < 100 {
		proportion := float64(line) / 100
		value := proportion * float64(len(requestTimeList))
		requestTime = float64(requestTimeList[int(value)])
	}
	return

}
