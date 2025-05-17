package model

import (
	"strconv"
	"strings"
)

type HttpDebugResponse struct {
	AssertAllPass     bool            `json:"assert_all_pass"`
	ExtractAllSuccess bool            `json:"extract_all_success"`
	CaseName          string          `json:"case_name"`
	ResponseContent   string          `json:"response_content"`
	ReqHeader         []HeadersForm   `json:"request_header"`
	RespHeader        []HeadersForm   `json:"response_header"`
	StatusCode        int             `json:"status_code"`
	ResponseTime      int64           `json:"response_time"`
	RequestContent    string          `json:"request_content"`
	URL               string          `json:"url"`
	MethodType        string          `json:"method_type"`
	ResponseSize      ResponseSize    `json:"response_size"`
	AssertRes         []*AssertItem   `json:"assert_res"`
	VariableRes       []*VariableItem `json:"variable_res"`
}

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

func Compare(extractValue string, expectValue string, checkingRule int) bool {
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

const (
	JsonPath = 1
	Regex    = 2
	Xpath    = 3
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
