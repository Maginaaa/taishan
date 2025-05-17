package model

import "sync"

const (
	NormalScene = iota
	PreFixScene
)

type SceneCaseTree struct {
	*SceneCase
	Children []*SceneCaseTree `json:"children,omitempty"`
}

type SceneCase struct {
	CaseID   int32       `json:"case_id"`   // case id
	Title    string      `json:"title"`     // case名
	ParentId int32       `json:"parent_id"` // 父节点id
	Type     CaseType    `json:"type"`
	SceneId  int32       `json:"scene_id"`
	Disabled bool        `json:"disabled"` // 是否执行
	Extend   interface{} `json:"extend"`
}

type CaseType int32

func (t CaseType) ToInt32() int32 {
	return int32(t)
}

type SceneRate struct {
	SceneId     int32   `json:"scene_id"`
	Rate        float32 `json:"rate"`
	Concurrency int32   `json:"concurrency"`
}

type Scene struct {
	SceneID        int32            `json:"scene_id"`         // 场景id
	SceneName      string           `json:"scene_name"`       // 场景名
	SceneType      int32            `json:"scene_type"`       // 场景类型
	ExportDataInfo ExportDataInfo   `json:"export_data_info"` // 导出数据信息
	CreateUserID   int32            `json:"create_user_id"`   // 创建人id
	CreateUserName string           `json:"create_user_name"` // 创建人
	UpdateUserID   int32            `json:"update_user_id"`   // 最后修改人id
	UpdateUserName string           `json:"update_user_name"` // 最后修改人
	CreateTime     string           `json:"create_time"`      // 创建时间
	UpdateTime     string           `json:"update_time"`      // 最后修改时间
	Remark         string           `json:"remark"`           // 场景描述
	VariablePool   VariablePool     `json:"variable_pool"`
	Cases          []*SceneCaseTree `json:"cases"` // 场景case
	MarshalCases   []MarshalCase    `json:"marshal_cases"`
	FileList       *FileList        `json:"file_list"` // 自定义文件
	//Result SceneExecutionResult `json:"execution_result"` // 场景运行结果
	DefaultHeader []*ParamsForm `json:"default_header"`
}

type ExportDataInfo struct {
	VariableList []string `json:"variable_list"`
	ExportTimes  int64    `json:"export_times"`
	Concurrency  int64    `json:"concurrency"`
	HasCache     bool     `json:"has_cache"`
	DisableCache bool     `json:"disable_cache"`
}

type ScenePressInfo struct {
	Scene       *Scene `json:"scene"`
	Rate        int64  `json:"rate"`
	Concurrency int64  `json:"concurrency"`
	Iteration   int64  `json:"iteration"`
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

const (
	HttpCaseType CaseType = iota + 1
	_
	_
	_
	_
	_
	_
	_
	_
	_
	LogicControlType
)

func (s *Scene) variablePoolInit() {
	if s.VariablePool.VariableList == nil {
		s.VariablePool.VariableList = make([]*Variable, 0)
	}
	if s.VariablePool.VariableMap == nil {
		s.VariablePool.VariableMap = new(sync.Map)
	}
}

func (s *Scene) LoadFileVariable() {
	s.variablePoolInit()
	variableMap := s.FileList.GetVariable()
	for k, v := range variableMap {
		s.VariablePool.VariableMap.Store(k, v)
	}
}

func (s *Scene) LoadSceneVariablePool() {
	s.variablePoolInit()
	for _, v := range s.VariablePool.VariableList {
		s.VariablePool.VariableMap.Store(v.VariableName, v.VariableVal)
	}
}

func (s *Scene) StoreVariablePool(caseData *HttpCase) {
	caseData.VariablePool.VariableMap.Range(func(key, value any) bool {
		s.VariablePool.VariableMap.Store(key, value)
		return true
	})
}

// CasesMarshal
// 压测过程中进行Unmarshal会有较大资源开销，压测前统一反序列化可节省资源
// 实测DeepCopy逻辑可将Unmarshal逻辑cpu消耗由24.2%降低至1.9%
func (s *Scene) CasesMarshal() {
	for _, cs := range s.Cases {
		if cs.Disabled {
			continue
		}
		switch cs.SceneCase.Type {
		case HttpCaseType:
			var caseData HttpCase
			caseData.Unmarshal(cs.SceneCase)
			s.MarshalCases = append(s.MarshalCases, MarshalCase{
				Type:     HttpCaseType,
				HttpCase: caseData,
			})
		case LogicControlType:
			var logicControl LogicControl
			logicControl.Unmarshal(cs.SceneCase.Extend)
			for _, c := range cs.Children {
				var caseData HttpCase
				caseData.Unmarshal(c.SceneCase)
				logicControl.Children = append(logicControl.Children, caseData)
			}
			s.MarshalCases = append(s.MarshalCases, MarshalCase{
				Type:         LogicControlType,
				LogicControl: logicControl,
			})
		default:
			continue
		}
	}
}

type MarshalCase struct {
	Type         CaseType `json:"type"`
	HttpCase     HttpCase
	LogicControl LogicControl
}

func (s *Scene) DeepCopy() (dst Scene) {
	dst = Scene{
		SceneID:        s.SceneID,
		SceneName:      s.SceneName,
		SceneType:      s.SceneType,
		ExportDataInfo: s.ExportDataInfo,
		VariablePool: VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]*Variable, 0),
		},
		MarshalCases:  s.MarshalCases,
		FileList:      s.FileList,
		DefaultHeader: s.DefaultHeader,
	}
	if s.VariablePool.VariableList != nil {
		s.VariablePool.VariableMap.Range(func(key, value any) bool {
			dst.VariablePool.VariableMap.Store(key, value)
			return true
		})
		for _, param := range s.VariablePool.VariableList {
			vrb := &Variable{}
			*vrb = *param
			dst.VariablePool.VariableList = append(dst.VariablePool.VariableList, vrb)
		}
	}
	return
}
