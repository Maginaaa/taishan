package rao

import "encoding/json"

const (
	LoopType = iota + 1
	IfType
)

type LogicControl struct {
	ControlType  int        `json:"control_type"`
	ControlVal   string     `json:"control_val"`
	Children     []HttpCase `json:"children"`
	ParamOne     string     `json:"param_one"`
	ParamTwo     string     `json:"param_two"`
	CheckingRule int        `json:"checking_rule"`
}

func (l *LogicControl) Unmarshal(s interface{}) {
	arr, _ := json.Marshal(s)
	_ = json.Unmarshal(arr, &l)
}
