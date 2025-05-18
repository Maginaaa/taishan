package api

const (
	Ok = 0

	ErrParam = 10001 //参数错误

)

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok: "成功",

	ErrParam: "参数错误",
}

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok: "success",

	ErrParam: "param error",
}
