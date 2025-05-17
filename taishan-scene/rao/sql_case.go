package rao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"scene/utils"
	"strings"
	"sync"
)

type SqlCase struct {
	SceneID          int32                   `json:"scene_id"`
	CaseID           int32                   `json:"case_id"`
	CaseName         string                  `json:"case_name"`
	CaseType         int32                   `json:"case_type"`
	DbID             int32                   `json:"db_id"`
	Sql              string                  `json:"sql"`
	DbInfo           *DatabaseConnectionInfo `json:"db_info"`
	UseSceneVariable bool                    `json:"use_scene_variable"` // 使用场景变量
	VariablePool     *VariablePool           `json:"variable_pool"`
	AssertForm       []AssertForm            `json:"assert_form"` // 断言
	VariableForm     []VariableForm          `json:"variable_form"`

	ResponseData *HttpResponse `json:"response_data"`
}

func (s *SqlCase) InitVariablePool() {
	if s.VariablePool == nil {
		s.VariablePool = &VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]Variable, 0),
		}
	}
}

func (s *SqlCase) LoadVariablePool(pool VariablePool) {
	s.InitVariablePool()
	if pool.VariableList != nil && len(pool.VariableList) > 0 {
		for _, variable := range pool.VariableList {
			s.VariablePool.VariableMap.Store(variable.VariableName, variable.VariableVal)
		}
	}
	if pool.VariableMap != nil {
		pool.VariableMap.Range(func(key, value interface{}) bool {
			s.VariablePool.VariableMap.Store(key, value)
			return true
		})
	}
}

func (s *SqlCase) DoSql() error {
	//s.InitVariablePool()
	//startTime := time.Now()
	//isSuccess, sqlResult, err := DBPing(s.DbInfo, s.Sql)
	//if !isSuccess {
	//	return err
	//}
	//endTime := time.Now()
	//
	//dataBytes, err := json.Marshal(sqlResult)
	//if err != nil {
	//	return err
	//}
	//
	//respHeaders := make([]ParamsForm, 0)
	//
	//hr := &HttpResponse{
	//	CaseName:        s.CaseName,
	//	RequestContent:  s.Sql,
	//	ResponseContent: string(dataBytes),
	//	ResponseHeader:  respHeaders,
	//}
	//s.ResponseData = hr
	//
	//s.VariableExtract()
	//s.ResponseBodyAssert()
	//
	//s.ResponseData.ResponseTime = endTime.Sub(startTime).Milliseconds()
	//s.ResponseData.SendBytes = 0
	////s.ResponseData.StartTime = startTime.Format(FullTimeFormat)
	////s.ResponseData.EndTime = endTime.Format(FullTimeFormat)
	//s.ResponseData.ResponseHeader = respHeaders
	//s.ResponseData.StatusCode = 200
	//s.ResponseData.ResponseSize = ResponseSize{
	//	BodySize:   float64(len(hr.ResponseContent)) / 1024,
	//	HeaderSize: 0,
	//	TotalSize:  float64(len(hr.ResponseContent)) / 1024,
	//}
	return nil
}

// ResponseBodyAssert 断言
func (s *SqlCase) ResponseBodyAssert() {
	expectedRes := make([]*AssertItem, 0)
	isSuccess := true
	for _, as := range s.AssertForm {
		if !as.Enable {
			continue
		}
		if as.ExpectValue == "" && as.ExtractExpress == "" {
			continue
		}
		expectedItem := &AssertItem{
			ExpectedValue: as.ExpectValue,
			CheckingRule:  as.CheckingRule,
		}
		switch as.ExtractType {
		case JsonPath:
			res, err := utils.GetByJsonPath(s.ResponseData.ResponseContent, as.ExtractExpress)
			if err != nil {
				expectedItem.ExtractValue = "提取失败"
				expectedItem.AssertPass = false
				isSuccess = false
			} else {
				resStr := fmt.Sprintf("%v", res)
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
	s.ResponseData.AssertRes = expectedRes
	s.ResponseData.RequestSuccess = isSuccess
}

// VariableExtract 变量提取
func (s *SqlCase) VariableExtract() {
	variableArr := make([]*VariableItem, 0)
	extractAllSuccess := true
	for _, v := range s.VariableForm {
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
			res, err := utils.GetByJsonPath(s.ResponseData.ResponseContent, v.Key)
			if err != nil {
				variableItem.ActualRes = "提取失败"
				variableItem.ExtractSuccess = false
				extractAllSuccess = false
				break
			} else {
				variableItem.ActualRes = fmt.Sprintf("%v", res)
				variableItem.ExtractSuccess = true
			}
			break
		case Regex:
			break
		case Xpath:
			break
		}
		variableArr = append(variableArr, variableItem)
	}
	s.ResponseData.VariableRes = variableArr
	s.ResponseData.ExtractAllSuccess = extractAllSuccess
}

func DBPing(info *DatabaseConnectionInfo, sql string) (isSuccess bool, sqlResult []map[string]interface{}, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&timeout=5s", info.User, info.Password, info.Host, info.Port, info.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return false, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false, nil, err
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		return false, nil, err
	}
	if sql != "" {
		// 确保SQL语句以SELECT开头
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
			// 执行sql语句
			db.Raw(sql).Scan(&sqlResult)
		}
	}
	return true, sqlResult, nil
}
