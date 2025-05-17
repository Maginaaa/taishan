package rao

const (
	RuleTypeAvgRT = iota + 1
	RuleTypeAvgRps
	RuleTypeAvgSuccessRate
)

const (
	AssertTypeLte = iota + 1
	AssertTypeGte
)

type NormalizedApiBindReq struct {
	ID       int32                  `json:"id"`
	PlanID   int32                  `json:"plan_id"`
	Url      string                 `json:"url"`
	RuleList []NormalizedAssertRule `json:"rule_list"`
}

type NormalizedAssertRule struct {
	RuleType    int32 `json:"rule_type"`
	AssertType  int32 `json:"assert_type"`
	ExpectValue int32 `json:"expect_value"`
}

type NormalizedApiListReq struct {
	Url      string `json:"url"`
	Page     int    `json:"page" binding:"required,gt=0"`      // 分页的页码（必填字段，必须大于 0）
	PageSize int    `json:"page_size" binding:"required,gt=0"` // 每页显示的项目数（必填字段，必须大于 0）

}

type NormalizedPressReq struct {
	WorkflowID   int      `json:"workflow_id"`
	WorkflowName string   `json:"workflow_name"`
	WorkflowUrl  string   `json:"workflow_url"`
	ApiList      []string `json:"api_list"`
	HistoryID    int32    `json:"history_id"`
}

type FsNormalizedResTemplate struct {
	WorkflowName string           `json:"workflow_name"`
	WorkflowUrl  string           `json:"workflow_url"`
	AssertData   []NormalizedTask `json:"assert_data"`
}

type NormalizedTask struct {
	ID            int32                  `json:"id"`
	ParentID      int32                  `json:"parent_id"`
	Url           string                 `json:"url"`
	AssertDetail  []NormalizedAssertItem `json:"assert_detail"`
	ReportID      int32                  `json:"report_id"`
	ReportUrl     string                 `json:"report_url"`
	AssertType    int32                  `json:"assert_type"`
	OperationType bool                   `json:"operation_type"`
	Remark        string                 `json:"remark"`
	CreateTime    string                 `json:"create_time"`
}

type AssertResEntity struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

type NormalizedResultListReq struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    bool   `json:"status"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type NormalizedResultList struct {
	ID         int32            `json:"id"`
	WorkflowID int32            `json:"workflow_id"`
	Status     bool             `json:"status"`
	TaskList   []NormalizedTask `json:"task_list"`
	CreateTime string           `json:"create_time"`
	UpdateTime string           `json:"update_time"`
}

type NormalizedResultHistogramRes struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type NormalizedApiAnalysis struct {
	Url     string `json:"url"`
	Total   int32  `json:"total"`
	Success int32  `json:"success"`
	Failure int32  `json:"failure"`
	Ignored int32  `json:"ignore"`
}

type NormalizedAssertItem struct {
	RuleType    int32   `json:"rule_type"`
	AssertType  int32   `json:"assert_type"`
	ExpectValue float64 `json:"expect_value"`
	ActualValue float64 `json:"actual_value"`
	Result      bool    `json:"result"`
}

type DisposeTaskReq struct {
	ID     int32  `json:"id"`
	Remark string `json:"remark"`
}

type TaskListReq struct {
	AssertType    int32  `json:"assert_type"`
	OperationType bool   `json:"operation_type"`
	Page          int    `json:"page" binding:"required,gt=0"` // 分页的页码（必填字段，必须大于 0）
	PageSize      int    `json:"page_size" binding:"required,gt=0"`
	Url           string `json:"url"`
}

type UrlHistoryReq struct {
	Url string `json:"url"`
}
