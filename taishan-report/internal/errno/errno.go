package errno

const (
	Ok = 0

	ErrSceneParam            = 10001 //场景参数错误
	ErrSceneID               = 10002 //场景id错误
	ErrSceneIDNotExist       = 10003 //场景id不存在
	ErrSceneNameAlreadyExist = 10004 //场景名称已存在
	ErrEmptySceneTestCase    = 10005 //场景用例不能为空
	ErrSceneCaseNameIsExist  = 10006 //同一场景下用例名称不能重复
	ErrSceneMysqlFailed      = 10007 //场景mysql操作失败
	ErrSceneCaseParam        = 10008 //场景case参数错误
	ErrSceneCaseID           = 10009 //场景case id错误
	ErrSceneCaseMysqlFailed  = 10010 //场景case mysql操作失败

	ErrPlanParam            = 20001 //计划参数错误
	ErrPlanID               = 20002 //计划id错误
	ErrPlanIDNotExist       = 20003 //计划id不存在
	ErrPlanNameAlreadyExist = 20004 //计划名称已存在
	ErrPlanMysqlFailed      = 20005 //计划mysql操作失败

	ErrReportParam  = 30001 //报告参数错误
	ErrReportFailed = 30002 // 获取报告失败
	ErrReportChange = 30003 // 压测过程修改失败
)

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok: "成功",

	ErrSceneParam:            "场景参数错误",
	ErrSceneID:               "场景id错误",
	ErrSceneNameAlreadyExist: "场景名称已存在",
	ErrEmptySceneTestCase:    "场景用例不能为空",
	ErrSceneCaseNameIsExist:  "同一场景下用例名称不能重复",
	ErrSceneMysqlFailed:      "场景mysql操作失败",
	ErrSceneCaseParam:        "场景case参数错误",
	ErrSceneCaseID:           "场景case id错误",
	ErrSceneCaseMysqlFailed:  "场景case mysql操作失败",

	ErrPlanParam:            "计划参数错误",
	ErrPlanID:               "计划id错误",
	ErrPlanIDNotExist:       "计划id不存在",
	ErrPlanNameAlreadyExist: "计划名称已存在",
	ErrPlanMysqlFailed:      "计划mysql操作失败",

	ErrReportParam:  "报告参数错误",
	ErrReportFailed: "获取报告失败",
}

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok: "success",

	ErrSceneParam:            "scene param error",
	ErrSceneID:               "scene id error",
	ErrSceneNameAlreadyExist: "scene name already exist",
	ErrEmptySceneTestCase:    "scene cases cannot be empty",
	ErrSceneCaseNameIsExist:  "scene casename already exist",
	ErrSceneMysqlFailed:      "scene mysql operate failed",

	ErrPlanParam:            "plan param error",
	ErrPlanID:               "plan id error",
	ErrPlanIDNotExist:       "plan id not exist",
	ErrPlanNameAlreadyExist: "plan name already exist",
	ErrPlanMysqlFailed:      "plan mysql operate failed",
}
