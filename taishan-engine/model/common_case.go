package model

const (
	// FullTimeFormat 时间格式
	FullTimeFormat = "2006-01-02 15:04:05"
	//Nanosecond     = "2006-01-02 15:04:05.000000000"
)

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
