package rao

type TraceListReq struct {
	ReportId int32 `json:"report_id"`
	CaseId   int32 `json:"case_id"`
}

type TraceListResp struct {
	Duration      int32  `json:"Duration"`
	OperationName string `json:"OperationName"`
	ServiceIp     string `json:"ServiceIp"`
	ServiceName   string `json:"ServiceName"`
	SpanID        string `json:"SpanID"`
	Timestamp     int64  `json:"Timestamp"`
	TraceID       string `json:"TraceID"`
}

type TraceDetailReq struct {
	ReportId int32  `json:"report_id"`
	TraceId  string `json:"trace_id"`
}

type TraceTreeNode struct {
	Id            string          `json:"id"`
	ParentId      string          `json:"parent_id"`
	Duration      int64           `json:"duration"`
	Children      []TraceTreeNode `json:"children"`
	App           string          `json:"app"`
	OperationName string          `json:"operation_name"`
	IP            string          `json:"ip"`
	Timestamp     int64           `json:"timestamp"`
	Status        string          `json:"status"`
}
