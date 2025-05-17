package rao

type CaseTracesReq struct {
	CaseID   int32  `json:"case_id"`
	ReportID int32  `json:"report_id"`
	TraceID  string `json:"trace_id"`
}

type TraceInfosEntity struct {
	Duration      int64  `json:"Duration"`
	OperationName string `json:"OperationName"`
	ServiceIp     string `json:"ServiceIp"`
	ServiceName   string `json:"ServiceName"`
	SpanID        string `json:"SpanID"`
	Timestamp     int64  `json:"Timestamp"`
	TraceID       string `json:"TraceID"`
}
