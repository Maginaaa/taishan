package model

import "sync"

type VariablePool struct {
	VariableList []*Variable `json:"variable_list"` // 变量列表
	VariableMap  *sync.Map   `json:"variable_pool"` // 变量池
}

func (v *VariablePool) init() {
	if v == nil {
		v = &VariablePool{
			VariableMap:  new(sync.Map),
			VariableList: make([]*Variable, 0),
		}
	} else {
		if v.VariableMap == nil {
			v.VariableMap = new(sync.Map)
		}
	}
	for _, variable := range v.VariableList {
		v.VariableMap.Store(variable.VariableName, variable.VariableVal)
	}
}

func (v *VariablePool) Save(variable *Variable) {
	v.init()
	v.VariableMap.Store(variable.VariableName, variable.VariableVal)
}

func (v *VariablePool) SaveList(VariableList []*Variable) {
	v.init()
	for _, variable := range VariableList {
		v.VariableMap.Store(variable.VariableName, variable.VariableVal)
	}
}

func (v *VariablePool) Get(variableName string) (variableVal string, ok bool) {
	res, ok := v.VariableMap.Load(variableName)
	if !ok {
		return "", false
	}
	return res.(string), true
}
