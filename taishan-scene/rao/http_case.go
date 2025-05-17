package rao

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"reflect"
	"scene/internal/biz/log"
	"scene/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
	VariablePool   *VariablePool   `json:"variable_pool"`
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
	ExtractExpress string `json:"extract_express" bson:"extract_express"`
	ExtractType    int    `json:"extract_type" bson:"extract_type"`
}

type VariableForm struct {
	Enable      bool   `json:"enable" bson:"enable"`
	ExtractType int    `json:"extract_type" bson:"extract_type"`
	Key         string `json:"key" bson:"key"`
	Value       string `json:"value" bson:"value"`
	Desc        string `json:"desc,omitempty" bson:"desc"`
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

type ResponseSize struct {
	BodySize   float64 `json:"body_size" bson:"body_size"`     // 单位KB
	HeaderSize float64 `json:"header_size" bson:"header_size"` // 单位KB
	TotalSize  float64 `json:"total_size" bson:"total_size"`   // 单位KB
}

// Unmarshal 对ext内容进行反序列化
func (h *HttpCase) Unmarshal(s interface{}) {
	arr, err := json.Marshal(s)
	if err != nil {
		log.Logger.Error("json.Marshal err", err)
		return
	}
	err = json.Unmarshal(arr, &h)
	if err != nil {
		log.Logger.Error("json.Unmarshal err", err)
		return
	}
}

func (h *HttpCase) initVariablePool() {
	if h.VariablePool == nil {
		h.VariablePool = &VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]Variable, 0),
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

func (h *HttpCase) DoRequest() error {
	h.initVariablePool()
	client := &fasthttp.Client{
		TLSConfig:                 &tls.Config{InsecureSkipVerify: true},
		MaxIdemponentCallAttempts: 1,
	}
	timeout := int64(5 * 1000)
	if h.OvertimeConfig != nil {
		if h.OvertimeConfig.Enable {
			timeout = h.OvertimeConfig.TimeoutPeriod
		}
	}
	req := fasthttp.AcquireRequest()
	req.SetTimeout(time.Duration(timeout) * time.Millisecond)
	h.ReplaceVariables()
	h.setHeader(req)
	req.SetRequestURI(strings.TrimSpace(h.URL))
	req.SetBody([]byte(h.Body.BodyValue))

	resp := fasthttp.AcquireResponse()
	defer func() {
		// 释放请求和响应对象
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	startTime := time.Now()
	err := client.Do(req, resp)
	endTime := time.Now()

	hr := &HttpResponse{
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
		json.Unmarshal(resp.Body(), &hr.ResponseContentInterface)
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

	return nil
}

func (h *HttpCase) setHeader(req *fasthttp.Request) {
	req.Header.SetMethod(h.MethodType)
	if h.Body == nil {
		h.Body = &Body{
			BodyType:  JSON,
			BodyValue: "",
		}
	}
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

//func (h *HttpCase) LoadVariablePool() {
//	s.VariablePool.Range(func(key, value interface{}) bool {
//		h.VariablePool.LoadOrStore(key, value)
//		return true
//	})
//}

func (h *HttpCase) ReplaceVariables() {
	replaceVariables(&(h.URL), h.VariablePool.VariableMap)
	replaceVariables(&(h.Body.BodyValue), h.VariablePool.VariableMap)
	for _, header := range h.HeadersForm {
		replaceVariables(&header.Value, h.VariablePool.VariableMap)
	}
}

// 支持预期值使用变量
func (h *HttpCase) replaceAssertVariables() {
	for _, as := range h.AssertForm {
		replaceVariables(&(as.ExpectValue), h.VariablePool.VariableMap)
	}
}

// ResponseBodyAssert 断言
func (h *HttpCase) ResponseBodyAssert() {
	h.replaceAssertVariables()
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
					expectedItem.AssertPass = compare(resStr, as.ExpectValue, as.CheckingRule)
					isSuccess = isSuccess && expectedItem.AssertPass
				}
				expectedRes = append(expectedRes, expectedItem)
				break
			case Regex:
				fmt.Println("regex")
				break
			case Xpath:
				fmt.Println("xpath")
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
