package rao

import (
	"strconv"
	"sync"
	"time"
)

const (
	NormalScene = iota
	PreFixScene

	MaxWaitingTime = 30 * 1000 // case前后置等待时间，单位ms

)

type Scene struct {
	SceneID        int32                `json:"scene_id"`         // 场景id
	SceneName      string               `json:"scene_name"`       // 场景名
	SceneType      int32                `json:"scene_type"`       // 场景类型
	PlanID         int32                `json:"plan_id"`          // 计划id
	Sort           int32                `json:"sort"`             // 排序
	ExportDataInfo ExportDataInfo       `json:"export_data_info"` // 导出数据信息
	Disabled       bool                 `json:"disabled"`         // 是否禁用
	CreateUserID   int32                `json:"create_user_id"`   // 创建人id
	CreateUserName string               `json:"create_user_name"` // 创建人
	UpdateUserID   int32                `json:"update_user_id"`   // 最后修改人id
	UpdateUserName string               `json:"update_user_name"` // 最后修改人
	CreateTime     string               `json:"create_time"`      // 创建时间
	UpdateTime     string               `json:"update_time"`      // 最后修改时间
	VariablePool   VariablePool         `json:"variable_pool"`
	Cases          []*SceneCaseTree     `json:"cases"`            // 场景case
	Result         SceneExecutionResult `json:"execution_result"` // 场景运行结果，仅供scene调试执行
	DefaultHeader  []*ParamsForm        `json:"default_header"`
}

type ExportDataInfo struct {
	VariableList []string   `json:"variable_list"`
	ExportTimes  int32      `json:"export_times"`
	Concurrency  int32      `json:"concurrency"`
	ExportList   []Variable `json:"-"`
	HasCache     bool       `json:"has_cache"`
	DisableCache bool       `json:"disable_cache"`
}

type SceneCreator struct {
	CreateUserID   int32  `json:"create_user_id"`   // 创建人id
	CreateUserName string `json:"create_user_name"` // 创建人
}

type Variable struct {
	VariableID   int32  `json:"variable_id"`   // 变量id
	VariableName string `json:"variable_name"` // 变量名
	VariableVal  string `json:"variable_val"`  // 变量值
	Remark       string `json:"remark"`        // 变量描述
}

type SceneVariable struct {
	SceneID int32 `json:"scene_id"`
	Variable
}

type SceneExecutionResult struct {
	SceneID        int32           `json:"scene_id"`       // 场景id
	SceneName      string          `json:"scene_name"`     // 场景名称
	Passed         bool            `json:"passed"`         // 当前场景是否执行通过
	TotalRequests  int             `json:"total_requests"` // 场景请求总数
	SuccessCount   int             `json:"success_count"`  // 成功接口数
	StartTime      string          `json:"start_time"`     // 开始时间
	EndTime        string          `json:"end_time"`       // 结束时间
	TotalTime      float64         `json:"total_time"`     // 总耗时
	CaseResultList []*HttpResponse `json:"case_result_list"`
}

type SceneListResponse struct {
	SceneList []*Scene `json:"scene_list"`
	Total     int64    `json:"total"`
}

type SceneDebugRecord struct {
	Time       string                 `json:"time"`   // 记录时间
	Passed     bool                   `json:"passed"` // 当前调试是否通过
	ResultList []SceneExecutionResult `json:"case_result_list"`
}

const (
	// FullTimeFormat 时间格式
	FullTimeFormat = "2006-01-02 15:04:05"
	Nanosecond     = "2006-01-02 15:04:05.000000000"
)

func (s *Scene) variablePoolInit() {
	if s.VariablePool.VariableList == nil {
		s.VariablePool.VariableList = make([]Variable, 0)
	}
	if s.VariablePool.VariableMap == nil {
		s.VariablePool.VariableMap = new(sync.Map)
	}
}

func (s *Scene) LoadPlanVariablePool(p *Plan) {
	s.variablePoolInit()
	p.VariablePool.Range(func(key, value interface{}) bool {
		s.VariablePool.VariableMap.Store(key, value)
		return true
	})
}

func (s *Scene) LoadSceneVariablePool() {
	s.variablePoolInit()
	for _, v := range s.VariablePool.VariableList {
		s.VariablePool.VariableMap.Store(v.VariableName, v.VariableVal)
	}
}

func (s *Scene) RunScene() {
	s.LoadSceneVariablePool()
	s.Result.Passed = true
	starTime := time.Now()
outerLoop:
	for _, cs := range s.Cases {
		if cs.Disabled {
			continue
		}
		switch cs.SceneCase.Type {
		case HttpCaseType:
			if !s.RunHttpCase(cs) {
				break outerLoop
			}
		case LogicControlType:
			var logicControl LogicControl
			logicControl.Unmarshal(cs.SceneCase.Extend)
			switch logicControl.ControlType {
			case LoopType:
				// str转int
				loopCount, err := strconv.Atoi(logicControl.ControlVal)
				if err != nil {
					continue
				}
				for i := 0; i < loopCount; i++ {
					for _, c := range cs.Children {
						if !s.RunHttpCase(c) {
							break outerLoop
						}
					}
				}
			case IfType:
				replaceVariables(&logicControl.ParamOne, s.VariablePool.VariableMap)
				replaceVariables(&logicControl.ParamTwo, s.VariablePool.VariableMap)
				if compare(logicControl.ParamOne, logicControl.ParamTwo, logicControl.CheckingRule) {
					for _, c := range cs.Children {
						if !s.RunHttpCase(c) {
							break outerLoop
						}
					}
				}
				continue
			}
		default:
			continue
		}
	}
	endTime := time.Now()

	s.Result.SceneID = s.SceneID
	s.Result.SceneName = s.SceneName
	s.Result.StartTime = starTime.Format(FullTimeFormat)
	s.Result.TotalTime = endTime.Sub(starTime).Seconds()
	s.Result.EndTime = endTime.Format(FullTimeFormat)

	if s.SceneType == PreFixScene {
		var variableList []Variable
		for _, key := range s.ExportDataInfo.VariableList {
			var variableVal string
			if val, ok := s.VariablePool.VariableMap.Load(key); ok {
				variableVal = val.(string)
			}
			variableList = append(variableList, Variable{
				VariableName: key,
				VariableVal:  variableVal,
			})
		}
		s.ExportDataInfo.ExportList = variableList
	}
}

func (s *Scene) RunHttpCase(c *SceneCaseTree) bool {
	if c.Disabled {
		return true
	}
	var caseData HttpCase
	caseData.Unmarshal(c.Extend)
	caseData.LoadVariablePool(s.VariablePool)
	caseData.LoadDefaultHeader(s.DefaultHeader)
	// 前置等待
	if caseData.WaitingConfig != nil && caseData.WaitingConfig.PreWaitingSwitch {
		waitingTime := caseData.WaitingConfig.PreWaitingTime
		if waitingTime > MaxWaitingTime {
			waitingTime = MaxWaitingTime
		}
		time.Sleep(time.Millisecond * time.Duration(waitingTime))
	}
	err := caseData.DoRequest()
	if err != nil {
		// TODO:记录错误信息
		return false
	}
	// 后置等待
	if caseData.WaitingConfig != nil && caseData.WaitingConfig.PostWaitingSwitch {
		waitingTime := caseData.WaitingConfig.PostWaitingTime
		if waitingTime > MaxWaitingTime {
			waitingTime = MaxWaitingTime
		}
		time.Sleep(time.Millisecond * time.Duration(waitingTime))
	}

	s.Result.CaseResultList = append(s.Result.CaseResultList, caseData.ResponseData)
	s.StoreVariablePool(&caseData)
	s.Result.TotalRequests++
	if caseData.ResponseData.RequestSuccess && caseData.ResponseData.AssertSuccess {
		s.Result.SuccessCount++
		return true
	} else {
		s.Result.Passed = false
		return false
	}
}

func (s *Scene) StoreVariablePool(caseData *HttpCase) {
	caseData.VariablePool.VariableMap.Range(func(key, value any) bool {
		s.VariablePool.VariableMap.Store(key, value)
		return true
	})
}

type ScenePressResult struct {
	SceneID           int32                  `json:"scene_id"`            // 场景id
	SceneName         string                 `json:"scene_name"`          // 场景名称
	TotalCount        int64                  `json:"total_count"`         // 场景执行次数
	SuccessCount      int64                  `json:"success_count"`       // 场景执行成功次数
	TotalResponseTime int64                  `json:"total_response_time"` // 总响应时间
	AvgResponseTime   float64                `json:"avg_response_time"`   // 平均响应时间
	MinResponseTime   float64                `json:"min_response_time"`   // 最小响应时间
	MaxResponseTime   float64                `json:"max_response_time"`   // 最大响应时间
	TotalSendBytes    float64                `json:"total_send_bytes"`    // 发送流量
	TotalReceiveBytes float64                `json:"total_receive_bytes"` // 接受流量
	Tps               float64                `json:"tps"`
	Data              []ApiTestResultDataMsg `json:"data"`
}

type PlanInfo struct {
	PlanId   int32               `json:"plan_id"`
	ReportId int32               `json:"report_id"`
	Scenes   []*SceneInformation `json:"scenes"`
}
type SceneInformation struct {
	SceneId   int32          `json:"scene_id"`
	SceneName string         `json:"scene_name"`
	SceneRps  int64          `json:"scene_rps"`
	Cases     []*Measurement `json:"cases"`
}

type Measurement struct {
	CaseId                         int32   `json:"case_id"`
	CaseName                       string  `json:"case_name"`
	Time                           int64   `json:"time"`
	StartTime                      int64   `json:"start_time"`
	EndTime                        int64   `json:"end_time"`
	RallyPoint                     int64   `json:"rally_point"`
	FiftyRequestTimeLineValue      int64   `json:"fifty_request_time_line_value"`
	NinetyRequestTimeLineValue     int64   `json:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue int64   `json:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue int64   `json:"ninety_nine_request_time_line_value"`
	XxCodeNum                      int64   `json:"1xx_code_num"`
	XxCodeNum1                     int64   `json:"2xx_code_num"`
	XxCodeNum2                     int64   `json:"4xx_code_num"`
	XxCodeNum3                     int64   `json:"5xx_code_num"`
	ActualConcurrency              int64   `json:"actual_concurrency"`
	RequestNum                     int64   `json:"request_num"`
	RequestTime                    int64   `json:"request_time"`
	SuccessNum                     int64   `json:"success_num"`
	ErrorNum                       int64   `json:"error_num"`
	SendBytes                      float64 `json:"send_bytes"`
	ReceivedBytes                  float64 `json:"received_bytes"`
	TargetRps                      int64   `json:"target_rps"`
	CurrentRps                     float64 `json:"current_rps"`
}
