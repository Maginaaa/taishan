// Package errno 定义所有错误码
package errno

const (
	Ok = 0

	ErrActionParam = 10001 // Action参数错误
)

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok: "成功",

	ErrActionParam: "场景参数错误",
}

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok: "success",

	ErrActionParam: "action param error",
}
