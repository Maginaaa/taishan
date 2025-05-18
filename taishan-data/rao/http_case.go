package rao

const (
	JSON = "json"
)

type HttpCase struct {
	SceneID        int32           `json:"scene_id"`
	CaseID         int32           `json:"case_id"`
	CaseName       string          `json:"case_name"`
	CaseType       int32           `json:"case_type"`
	MethodType     string          `json:"method_type"`
	URL            string          `json:"url"`
	ParamsForm     []*ParamsForm   `json:"params_form"` //  url param
	Body           *Body           `json:"body"`
	HeadersForm    []*ParamsForm   `json:"headers_form"` // 响应头
	AssertForm     []*AssertForm   `json:"assert_form"`  // 断言
	VariableForm   []VariableForm  `json:"variable_form"`
	RallyPoint     *RallyPoint     `json:"rally_point"`
	WaitingConfig  *WaitingConfig  `json:"waiting_config"`
	OvertimeConfig *OvertimeConfig `json:"overtime_config"`

	ResponseData *HttpResponse `json:"response_data"`
}

type Assert struct {
	ExtractMethod string `json:"extract_method"`
	ExtractRule   string `json:"extract_rule"`
	ExpectedValue string `json:"expected_value"` // 预期值
	ExtractValue  string `json:"extract_value"`  // 实际提取值
	AssertPass    bool   `json:"assert_pass"`    // 断言结果
}

type Body struct {
	BodyType  string `json:"body_type"`
	BodyValue string `json:"body_value"`
}

type RallyPoint struct {
	Enable        bool   `json:"enable"`
	Concurrency   int64  `json:"concurrency"`
	TimeoutPeriod int64  `json:"timeout_period"`
	LuaScriptSHA  string `json:"lua_script_sha"`
}

type WaitingConfig struct {
	PreWaitingSwitch  bool  `json:"pre_waiting_switch"`
	PreWaitingTime    int64 `json:"pre_waiting_time"`
	PostWaitingSwitch bool  `json:"post_waiting_switch"`
	PostWaitingTime   int64 `json:"post_waiting_time"`
}

type OvertimeConfig struct {
	Enable        bool  `json:"enable"`
	TimeoutPeriod int64 `json:"timeout_period"`
}
