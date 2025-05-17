package execution

import (
	"engine/internal/biz/log"
	"engine/model"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DisposeScene 执行单场景
func DisposeScene(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg, goroutineId int64, executeCount *int64, currentRps int64) {

	s := sceneAction.ScenePressInfo.Scene.DeepCopy()

	// 获取所有参数化文件中的一行参数，放入场景的VariablePool
	s.LoadFileVariable()
	s.LoadSceneVariablePool()
outerLoop:
	for _, cs := range s.MarshalCases {
		switch cs.Type {
		case model.HttpCaseType:
			caseData := cs.HttpCase.DeepCopy()
			caseData.LoadVariablePool(s.VariablePool)
			caseData.LoadDefaultHeader(s.DefaultHeader)
			if executeCount != nil {
				if *executeCount >= currentRps {
					break outerLoop
				}
				atomic.AddInt64(executeCount, 1)
			}
			wg.Add(1)
			err := DisposeRequest(wg, &caseData, sceneAction, dataMsgCh, goroutineId)
			if err != nil {
				continue
			}
			s.StoreVariablePool(&caseData)
			if !caseData.ResponseData.RequestSuccess || !caseData.ResponseData.AssertSuccess {
				break outerLoop
			}
		case model.LogicControlType:
			logicControl := cs.LogicControl
			switch logicControl.ControlType {
			case model.LoopType:
				loopCount, err := strconv.Atoi(logicControl.ControlVal)
				if err != nil {
					continue
				}
				for i := 0; i < loopCount; i++ {
					for _, c := range logicControl.Children {
						caseData := c.DeepCopy()
						caseData.LoadVariablePool(s.VariablePool)
						caseData.LoadDefaultHeader(s.DefaultHeader)
						if executeCount != nil {
							if *executeCount >= currentRps {
								break outerLoop
							}
							atomic.AddInt64(executeCount, 1)
						}
						wg.Add(1)
						err = DisposeRequest(wg, &caseData, sceneAction, dataMsgCh, goroutineId)
						if err != nil {
							continue
						}
						s.StoreVariablePool(&caseData)
						if !caseData.ResponseData.RequestSuccess || !caseData.ResponseData.AssertSuccess {
							break outerLoop
						}
					}
				}
			case model.IfType:
				model.ReplaceVariables(&logicControl.ParamOne, s.VariablePool.VariableMap)
				model.ReplaceVariables(&logicControl.ParamTwo, s.VariablePool.VariableMap)
				if model.Compare(logicControl.ParamOne, logicControl.ParamTwo, logicControl.CheckingRule) {
					for _, caseData := range logicControl.Children {
						caseData.LoadVariablePool(s.VariablePool)
						caseData.LoadDefaultHeader(s.DefaultHeader)
						if executeCount != nil {
							if *executeCount >= currentRps {
								break outerLoop
							}
							atomic.AddInt64(executeCount, 1)
						}
						wg.Add(1)
						err := DisposeRequest(wg, &caseData, sceneAction, dataMsgCh, goroutineId)
						if err != nil {
							continue
						}
						s.StoreVariablePool(&caseData)
						if !caseData.ResponseData.RequestSuccess || !caseData.ResponseData.AssertSuccess {
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

}

func DisposeRequest(wg *sync.WaitGroup, httpCase *model.HttpCase, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg, goroutineId int64) (err error) {
	defer wg.Done()

	// 请求的数据处理
	requestResults := &model.HttpResultDataMsg{
		CaseID:   httpCase.CaseID,
		CaseName: httpCase.CaseName,
	}

	var (
		errType               = int64(0)
		rallyPointConcurrency = int64(0)
	)

	// 集合点
	if httpCase.RallyPoint != nil && httpCase.RallyPoint.Enable {
		rallyPointConcurrency = httpCase.RallyPoint.Concurrency
		begin := time.Now().UnixMilli()
		duration := int64(10) * 1000
		if httpCase.RallyPoint.TimeoutPeriod != 0 {
			duration = httpCase.RallyPoint.TimeoutPeriod
		}
		key, err := model.ReportRdb.EvalSha(httpCase.RallyPoint.LuaScriptSHA, []string{}).Result()
		if err != nil {
			log.Logger.Error("model.request_redis.DisposeRequest.EvalSha()", "err", err.Error())
		}
		newKey := fmt.Sprintf("report:index:%d:%d:%d", sceneAction.PlanID, httpCase.CaseID, key)
		for begin+duration >= time.Now().UnixMilli() {
			countStr, err := model.ReportRdb.Get(newKey).Result()
			if err != nil {
				log.Logger.Error("model.request_redis.DisposeRequest.Get(newKey)", "err", err.Error(), "newKey", newKey)
				break
			}
			countInt, err := strconv.Atoi(countStr)
			if err != nil {
				log.Logger.Error("model.request_redis.DisposeRequest.strconv.Atoi(countStr)", "err", err.Error(), "countStr", countStr)
				break
			}
			if int64(countInt) == httpCase.RallyPoint.Concurrency {
				break
			}
		}
	}

	// 前置等待
	if httpCase.WaitingConfig != nil && httpCase.WaitingConfig.PreWaitingSwitch {
		waitingTime := httpCase.WaitingConfig.PreWaitingTime
		if waitingTime > model.MaxWaitingTime {
			waitingTime = model.MaxWaitingTime
		}
		time.Sleep(time.Millisecond * time.Duration(waitingTime))
	}

	httpCase.DoRequest()

	// 后置等待
	if httpCase.WaitingConfig != nil && httpCase.WaitingConfig.PostWaitingSwitch {
		waitingTime := httpCase.WaitingConfig.PostWaitingTime
		if waitingTime > model.MaxWaitingTime {
			waitingTime = model.MaxWaitingTime
		}
		time.Sleep(time.Millisecond * time.Duration(waitingTime))
	}

	// 采样
	samplingType := sceneAction.SamplingInfo.SamplingType
	samplingRate := sceneAction.SamplingInfo.SamplingRate
	isSuccess := httpCase.ResponseData.RequestSuccess && httpCase.ResponseData.AssertSuccess
	needSampling := false
	if samplingType > 0 {
		coincidenceSamplingRate := (httpCase.ResponseData.StartTime%100)*100+httpCase.ResponseData.EndTime%100 < samplingRate
		if samplingType == model.SamplingWithPercentage && coincidenceSamplingRate {
			//model.GetMongoCollection(fmt.Sprintf("%d:%d", sceneAction.ReportID, sceneAction.SceneID)).InsertOne(context.Background(), httpCase.ResponseData)
			needSampling = true
		} else if samplingType == model.SamplingWithError && !isSuccess && coincidenceSamplingRate {
			needSampling = true
		}
	}
	if needSampling {
		requestResults.DebugInfo = httpCase.ResponseData
	}

	if !httpCase.ResponseData.RequestSuccess {
		errType = model.RequestErr
	} else if !httpCase.ResponseData.AssertSuccess {
		errType = model.AssertErr
	}
	urlArr := strings.Split(httpCase.URL, "?")
	uri := ""
	if len(urlArr) > 0 {
		uri = urlArr[0]
	}
	if dataMsgCh != nil {
		requestResults.StatusCode = httpCase.ResponseData.StatusCode
		requestResults.RequestTime = httpCase.ResponseData.ResponseTime
		requestResults.ErrorType = errType
		requestResults.IsSuccess = isSuccess
		requestResults.SendBytes = httpCase.ResponseData.SendBytes
		requestResults.ReceivedBytes = httpCase.ResponseData.ReceiverBytes
		requestResults.StartTime = httpCase.ResponseData.StartTime
		requestResults.EndTime = httpCase.ResponseData.EndTime
		requestResults.RallyPoint = rallyPointConcurrency
		requestResults.GoroutineID = goroutineId
		requestResults.Url = uri
		dataMsgCh <- &model.ResultDataMsg{
			MsgType: model.TypeHttpResultData,
			GlobalInfo: &model.GlobalInfo{
				ReportID:     sceneAction.ReportID,
				Timestamp:    time.Now().UnixMilli(),
				NeedSampling: needSampling,
			},
			SceneInfo: &model.SceneInfo{
				SceneID: sceneAction.SceneID,
			},
			HttpResultDataMsg: requestResults,
		}
	}
	return
}
