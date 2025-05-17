package model

import (
	"crypto/tls"
	"encoding/json"
	"engine/internal/biz/log"
	"engine/internal/utils"
	"fmt"
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	JSON = "json"

	MaxWaitingTime = 30 * 1000 // case前后置等待时间，单位ms

	DefaultLimitRt = int64(5 * 1000)
	MinLimitRt     = int64(10)
	MaxLimitRt     = int64(60 * 1000)

	RequestErr = 1
	AssertErr  = 2
)

type HttpResponse struct {
	// 基础信息
	CaseId   int32  `json:"case_id" bson:"case_id"`
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

type HttpCase struct {
	CaseID         int32           `json:"case_id"`
	CaseName       string          `json:"case_name"`
	MethodType     string          `json:"method_type"`
	URL            string          `json:"url"`
	VariablePool   *VariablePool   `json:"variable_pool"`
	ParamsForm     []*ParamsForm   `json:"params_form"` //  url param
	Body           *Body           `json:"body"`
	HeadersForm    []*ParamsForm   `json:"headers_form"` // 响应头
	AssertForm     []*AssertForm   `json:"assert_form"`  // 断言
	VariableForm   []VariableForm  `json:"variable_form"`
	RallyPoint     *RallyPoint     `json:"rally_point"`
	WaitingConfig  *WaitingConfig  `json:"waiting_config"`
	OvertimeConfig *OvertimeConfig `json:"overtime_config"`
	ResponseData   *HttpResponse   `json:"response_data"`
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

func (h *HttpCase) DeepCopy() (dst HttpCase) {
	dst = HttpCase{
		CaseID:     h.CaseID,
		CaseName:   h.CaseName,
		MethodType: h.MethodType,
		URL:        h.URL,
		VariablePool: &VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]*Variable, 0),
		},
		ParamsForm:    make([]*ParamsForm, 0),
		Body:          &Body{},
		HeadersForm:   make([]*ParamsForm, 0),
		AssertForm:    make([]*AssertForm, 0),
		VariableForm:  make([]VariableForm, 0),
		RallyPoint:    &RallyPoint{},
		WaitingConfig: &WaitingConfig{},
		ResponseData:  &HttpResponse{},
	}
	for _, param := range h.ParamsForm {
		dst.ParamsForm = append(dst.ParamsForm, param)
	}

	for _, param := range h.HeadersForm {
		header := &ParamsForm{}
		*header = *param
		dst.HeadersForm = append(dst.HeadersForm, header)
	}
	for _, param := range h.AssertForm {
		assert := &AssertForm{}
		*assert = *param
		dst.AssertForm = append(dst.AssertForm, assert)
	}
	for _, param := range h.VariableForm {
		dst.VariableForm = append(dst.VariableForm, param)
	}
	if h.VariablePool != nil {
		*dst.VariablePool = *h.VariablePool
	}
	if h.Body != nil {
		*dst.Body = *h.Body
	}
	if h.RallyPoint != nil {
		*dst.RallyPoint = *h.RallyPoint
	}
	if h.WaitingConfig != nil {
		*dst.WaitingConfig = *h.WaitingConfig
	}
	return
}

func (h *HttpCase) Unmarshal(cs *SceneCase) {
	arr, err := json.Marshal(cs.Extend)
	if err != nil {
		log.Logger.Error("json.Marshal err", err)
		return
	}
	err = json.Unmarshal(arr, &h)
	h.CaseID = cs.CaseID
	if err != nil {
		log.Logger.Error("json.Unmarshal err", err)
		return
	}
}

func (h *HttpCase) initVariablePool() {
	if h.VariablePool == nil {
		h.VariablePool = &VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]*Variable, 0),
		}
	}
}

var (
	KeepAliveClient *fasthttp.Client
	once            sync.Once
)

func (h *HttpCase) DoRequest() {
	var client *fasthttp.Client
	timeout := DefaultLimitRt
	if h.OvertimeConfig != nil {
		if h.OvertimeConfig.Enable && h.OvertimeConfig.TimeoutPeriod != 0 {
			timeout = h.OvertimeConfig.TimeoutPeriod
			if timeout < MinLimitRt {
				timeout = MinLimitRt
			}
			if timeout > MaxLimitRt {
				timeout = MaxLimitRt
			}
		}
	}
	req := fasthttp.AcquireRequest()
	req.SetTimeout(time.Duration(timeout) * time.Millisecond)
	newKeepAlive()
	client = KeepAliveClient
	h.ReplaceVariables()
	h.setHeader(req)
	req.SetRequestURI(strings.TrimSpace(h.URL))
	req.SetBody([]byte(h.Body.BodyValue))
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	startTime := time.Now()
	err := client.Do(req, resp)
	endTime := time.Now()

	hr := &HttpResponse{
		CaseId:        h.CaseID,
		CaseName:      h.CaseName,
		URL:           h.URL,
		MethodType:    h.MethodType,
		RequestHeader: h.HeadersForm,
	}
	respHeaders := make([]*ParamsForm, 0)
	if err == nil {
		resp.Header.VisitAll(func(key, value []byte) {
			respHeaders = append(respHeaders, &ParamsForm{
				Key:   string(key),
				Value: string(value),
			})
		})
		hr.ResponseContent = string(resp.Body())
		_ = json.Unmarshal(resp.Body(), &hr.ResponseContentInterface) // 前置反序列化，cpu消耗由19.7%降低至6.7%
		hr.StatusCode = resp.StatusCode()
		// 响应码>=400时，视为请求失败
		if resp.StatusCode() < 400 {
			hr.RequestSuccess = true
		}
	} else {
		hr.Err = err.Error()
		hr.ResponseContent = err.Error()
	}
	hr.ResponseHeader = respHeaders
	hr.ResponseTime = endTime.Sub(startTime).Milliseconds()
	hr.SendBytes = float64(len(hr.URL)+req.Header.ContentLength()+len(req.Header.String())) / 1024
	hr.StartTime = startTime.UnixMilli()
	hr.EndTime = endTime.UnixMilli()
	hr.ResponseSize = ResponseSize{
		BodySize:   float64(len(resp.Body())) / 1024,
		HeaderSize: float64(len(resp.Header.String())) / 1024,
		TotalSize:  float64(len(resp.Body())+len(resp.Header.String())) / 1024,
	}
	hr.ReceiverBytes = hr.ResponseSize.TotalSize
	hr.RequestContent = h.Body.BodyValue
	h.ResponseData = hr

	h.VariableExtract()
	h.ResponseBodyAssert()
}

func newKeepAlive() {
	once.Do(func() {
		KeepAliveClient = &fasthttp.Client{
			TLSConfig:                 &tls.Config{InsecureSkipVerify: true},
			MaxConnsPerHost:           10000,
			MaxIdemponentCallAttempts: 1,
			MaxIdleConnDuration:       time.Duration(5) * time.Second,
			ReadTimeout:               time.Duration(5) * time.Second,
			WriteTimeout:              time.Duration(5) * time.Second,
		}
	})
}

func (h *HttpCase) setHeader(req *fasthttp.Request) {
	req.Header.SetMethod(h.MethodType)
	switch h.Body.BodyType {
	case JSON:
		req.Header.Set("Content-Type", "application/json")
	}
	for _, header := range h.HeadersForm {
		if header.Enable && len(header.Key) > 0 && len(header.Value) > 0 {
			req.Header.Set(header.Key, header.Value)
		}
	}
}

func (h *HttpCase) LoadVariablePool(pool VariablePool) {
	h.initVariablePool()
	if pool.VariableList != nil && len(pool.VariableList) > 0 {
		for _, variable := range pool.VariableList {
			h.VariablePool.VariableList = append(h.VariablePool.VariableList, variable)
			h.VariablePool.VariableMap.Store(variable.VariableName, variable.VariableVal)
		}
	}
	if pool.VariableMap != nil {
		pool.VariableMap.Range(func(key, value interface{}) bool {
			h.VariablePool.VariableMap.Store(key, value)
			return true
		})
	}
}

func (h *HttpCase) LoadDefaultHeader(header []*ParamsForm) {
	h.HeadersForm = append(h.HeadersForm, header...)

}

func (h *HttpCase) ReplaceVariables() {
	ReplaceVariables(&(h.URL), h.VariablePool.VariableMap)
	ReplaceVariables(&(h.Body.BodyValue), h.VariablePool.VariableMap)
	for _, header := range h.HeadersForm {
		ReplaceVariables(&header.Value, h.VariablePool.VariableMap)
	}
}

// 支持预期值使用变量
func (h *HttpCase) replaceAssertVariables() {
	for _, as := range h.AssertForm {
		ReplaceVariables(&(as.ExpectValue), h.VariablePool.VariableMap)
	}
}

// ResponseBodyAssert 断言
func (h *HttpCase) ResponseBodyAssert() {
	expectedRes := make([]*AssertItem, 0)
	isSuccess := true
	if h.ResponseData.RequestSuccess {
		for _, as := range h.AssertForm {
			if !as.Enable {
				continue
			}
			if as.ExpectValue == "" && as.ExtractExpress == "" {
				continue
			}
			expectedItem := &AssertItem{
				ExtractRule:   as.ExtractExpress,
				ExpectedValue: as.ExpectValue,
				CheckingRule:  as.CheckingRule,
			}
			switch as.ExtractType {
			case JsonPath:
				res, err := utils.GetByJsonPath(h.ResponseData.ResponseContentInterface, as.ExtractExpress)
				if err != nil {
					expectedItem.ExtractValue = "提取失败"
					expectedItem.AssertPass = false
					isSuccess = false
				} else {
					resStr := fmt.Sprintf("%v", res)
					if res == nil {
						resStr = "null"
					}
					expectedItem.ExtractValue = resStr
					expectedItem.AssertPass = Compare(resStr, as.ExpectValue, as.CheckingRule)
					isSuccess = isSuccess && expectedItem.AssertPass
				}
				expectedRes = append(expectedRes, expectedItem)
				break
			case Regex:
				break
			case Xpath:
				break
			}
		}
	} else {
		isSuccess = false
	}
	h.ResponseData.AssertRes = expectedRes
	h.ResponseData.AssertSuccess = isSuccess
}

// VariableExtract 变量提取
func (h *HttpCase) VariableExtract() {
	variableArr := make([]*VariableItem, 0)
	extractAllSuccess := true
	if h.ResponseData.RequestSuccess {
		for _, v := range h.VariableForm {
			if !v.Enable {
				continue
			}
			if v.Key == "" || v.Value == "" {
				continue
			}
			variableItem := &VariableItem{
				ExtractType:  v.ExtractType,
				ExtractRule:  v.Key,
				VariableName: v.Value,
			}

			switch v.ExtractType {
			case JsonPath:
				res, err := utils.GetByJsonPath(h.ResponseData.ResponseContentInterface, v.Key)
				if err != nil {
					variableItem.ActualRes = "提取失败"
					variableItem.ExtractSuccess = false
					extractAllSuccess = false
					break
				} else {
					variableItem.ActualRes = anyToStr(res)
					variableItem.ExtractSuccess = true
				}
				break
			case Regex:
				break
			case Xpath:
				break
			}
			if variableItem.ExtractSuccess {
				h.VariablePool.VariableMap.Store(variableItem.VariableName, variableItem.ActualRes)
			}
			variableArr = append(variableArr, variableItem)
		}
	}
	h.ResponseData.VariableRes = variableArr
	h.ResponseData.ExtractAllSuccess = extractAllSuccess
}

func anyToStr(i interface{}) string {
	if i == nil {
		return ""
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Complex128:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Array:
		str, _ := json.Marshal(i)
		return string(str)
	default:
		return ""
	}
}
