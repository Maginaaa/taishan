package rao

import (
	"reflect"
	"scene/internal/biz/log"
	"strconv"
	"strings"
)

type VariableItem struct {
	ExtractType    int    `json:"extract_method" bson:"extract_type"`
	ExtractRule    string `json:"extract_rule" bson:"extract_rule"`
	VariableName   string `json:"variable_name" bson:"variable_name"`
	ActualRes      string `json:"actual_res" bson:"actual_res"`
	ExtractSuccess bool   `json:"extract_success" bson:"extract_success"`
}

type AssertItem struct {
	ExtractMethod string `json:"extract_method" bson:"extract_method"`
	ExtractRule   string `json:"extract_rule" bson:"extract_rule"`
	ExpectedValue string `json:"expected_value" bson:"expected_value"` // 预期值
	ExtractValue  string `json:"extract_value" bson:"extract_value"`   // 实际提取值
	CheckingRule  int    `json:"checking_rule" bson:"checking_rule"`
	AssertPass    bool   `json:"assert_pass" bson:"assert_pass"` // 断言结果
}

type HttpResponse struct {
	// 基础信息
	CaseID   int32  `json:"case_id" bson:"case_id"`
	CaseName string `json:"case_name" bson:"case_name"`

	// 请求信息
	URL            string        `json:"url" bson:"url"`
	MethodType     string        `json:"method_type" bson:"method_type"`
	RequestHeader  []*ParamsForm `json:"request_header" bson:"request_header"`
	RequestContent string        `json:"request_content" bson:"request_content"`

	// 响应信息
	StatusCode               int             `json:"status_code" bson:"status_code"`
	ResponseContent          string          `json:"response_content" bson:"response_content"`
	ResponseHeader           []*ParamsForm   `json:"response_header" bson:"response_header"`
	AssertRes                []*AssertItem   `json:"assert_res" bson:"assert_res"`
	VariableRes              []*VariableItem `json:"variable_res" bson:"variable_res"`
	ResponseContentInterface interface{}     `json:"response_content_interface" bson:"response_content_interface"`

	// 统计信息
	AssertSuccess     bool         `json:"assert_success" bson:"assert_success"`
	ExtractAllSuccess bool         `json:"extract_all_success" bson:"extract_all_success"`
	ResponseSize      ResponseSize `json:"response_size" bson:"response_size"`
	ResponseTime      int64        `json:"response_time" bson:"response_time"`
	SendBytes         float64      `json:"send_bytes" bson:"send_bytes"`
	ReceiverBytes     float64      `json:"receiver_bytes" bson:"receiver_bytes"`
	RequestSuccess    bool         `json:"request_success" bson:"request_success"`
	Err               string       `json:"err" bson:"err"`
	StartTime         int64        `json:"start_time" bson:"start_time"`
	EndTime           int64        `json:"end_time" bson:"end_time"`
}

func compare(extractValue string, expectValue string, checkingRule int) bool {
	switch checkingRule {
	case Equals:
		return extractValue == expectValue
	case NotEquals:
		return extractValue != expectValue
	case Contains:
		return strings.Contains(extractValue, expectValue)
	case NotContains:
		return !strings.Contains(extractValue, expectValue)
	case Greater:
		v1, err := strconv.ParseFloat(extractValue, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseFloat(expectValue, 64)
		if err != nil {
			return false
		}
		return v1 > v2
	case Less:
		v1, err := strconv.ParseFloat(extractValue, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseFloat(expectValue, 64)
		if err != nil {
			return false
		}
		return v1 < v2
	case GreaterOrEquals:
		v1, err := strconv.ParseFloat(extractValue, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseFloat(expectValue, 64)
		if err != nil {
			return false
		}
		return v1 >= v2
	case LessOrEquals:
		v1, err := strconv.ParseFloat(extractValue, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseFloat(expectValue, 64)
		if err != nil {
			return false
		}
		return v1 <= v2
	}
	return false
}

func getValue(extractValue interface{}, expectValue string) (v1, v2 interface{}) {
	switch reflect.TypeOf(extractValue).Kind() {
	case reflect.String:
		return extractValue.(string), expectValue
	case reflect.Float64:
		f, _ := strconv.ParseFloat(expectValue, 64)
		return extractValue.(float64), f
	case reflect.Bool:
		b, _ := strconv.ParseBool(expectValue)
		return extractValue.(bool), b
	}
	log.Logger.Error("类型转换失败")
	return nil, nil
}

const (
	Default = iota // 外部携带的变量
	JsonPath
	Regex
	Xpath
)

const (
	Equals = iota + 1
	NotEquals
	Contains
	NotContains
	Greater
	Less
	GreaterOrEquals
	LessOrEquals
)
