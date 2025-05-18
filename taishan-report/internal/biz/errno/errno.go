// Package errno 定义所有错误码
package errno

const (
	Ok = 0

	ErrParam = 10001 //通用参数错误

	ErrToken   = 11001 //无效的token
	ErrGetUser = 11002 //获取用户信息错误",
	ErrLogout  = 11003 //登出失败

	ErrSceneParam            = 20001 //场景参数错误
	ErrSceneID               = 20002 //场景id错误
	ErrSceneIDNotExist       = 20003 //场景id不存在
	ErrSceneNameAlreadyExist = 20004 //场景名称已存在
	ErrEmptySceneTestCase    = 20005 //场景用例不能为空
	ErrSceneCaseNameIsExist  = 20006 //同一场景下用例名称不能重复
	ErrSceneCreateFailed     = 20007 //场景创建失败
	ErrSceneCaseParam        = 20008 //场景case参数错误
	ErrSceneCaseID           = 20009 //场景case id错误
	ErrSceneCaseFailed       = 20010 //场景case mysql操作失败
	ErrSceneDebugFailed      = 20011 //场景debug失败
	ErrSceneVariableID       = 20012 //场景变量id错误
	ErrSceneVariableCreate   = 20013 //场景变量创建失败
	ErrSceneVariableUpdate   = 20014 //场景变量更新失败
	ErrDBParam               = 20015 //数据连接参数错误
	ErrDBID                  = 20016 //id错误
	ErrDBMysqlFailed         = 20017 //数据库操作失败
	ErrSceneDeleteFailed     = 20018 // 场景删除失败
	ErrSceneUpdateFailed     = 20019

	ErrPlanParam            = 30001 //计划参数错误
	ErrPlanID               = 30002 //计划id错误
	ErrPlanIDNotExist       = 30003 //计划id不存在
	ErrPlanNameAlreadyExist = 30004 //计划名称已存在
	ErrPlanMysqlFailed      = 30005 //计划mysql操作失败
	ErrPlanDebugFailed      = 30006 //计划debug失败
	ErrPlanUploadFileFailed = 30007 //参数化文件上传失败
	ErrPlanExecute          = 30008 // 计划执行失败
	ErrPlanUpdate           = 30009 // 计划更新失败
	RunningPlanEnableChange = 30010
	ErrPlanDeleteFailed     = 30011
	ErrColumnRename         = 30012

	ErrDownloadFileFailed = 310000

	ErrReportParam       = 40001 //报告参数错误
	ErrReportRedisFailed = 40002 //报告redis操作失败
	ErrReportChange      = 40003 // 压测过程修改失败
)

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok:       "成功",
	ErrParam: "参数错误",

	ErrToken:   "无效的token",
	ErrGetUser: "获取用户信息错误",
	ErrLogout:  "登出失败",

	ErrSceneParam:            "场景参数错误",
	ErrSceneID:               "场景id错误",
	ErrSceneNameAlreadyExist: "场景名称已存在",
	ErrEmptySceneTestCase:    "场景用例不能为空",
	ErrSceneCaseNameIsExist:  "同一场景下用例名称不能重复",
	ErrSceneCreateFailed:     "场景创建失败",
	ErrSceneCaseParam:        "场景case参数错误",
	ErrSceneCaseID:           "场景case id错误",
	ErrSceneCaseFailed:       "场景case 操作失败",
	ErrSceneDebugFailed:      "场景debug失败",
	ErrSceneVariableID:       "场景变量id错误",
	ErrDBParam:               "数据连接参数错误",
	ErrDBID:                  "id错误",
	ErrDBMysqlFailed:         "数据库操作失败",
	ErrSceneDeleteFailed:     "场景删除失败",
	ErrSceneUpdateFailed:     "场景更新失败",

	ErrPlanParam:            "计划参数错误",
	ErrPlanID:               "计划id错误",
	ErrPlanIDNotExist:       "计划id不存在",
	ErrPlanNameAlreadyExist: "计划名称已存在",
	ErrPlanMysqlFailed:      "计划mysql操作失败",
	ErrPlanUploadFileFailed: "参数化文件上传失败",
	ErrPlanDebugFailed:      "计划debug失败",
	ErrPlanExecute:          "计划执行失败",
	ErrPlanUpdate:           "计划更新失败",
	RunningPlanEnableChange: "计划正在执行中，无法修改",
	ErrPlanDeleteFailed:     "计划删除失败",
	ErrColumnRename:         "文件变量重命名失败",

	ErrDownloadFileFailed: "文件下载失败",

	ErrReportParam:       "报告参数错误",
	ErrReportRedisFailed: "报告redis操作失败",
}

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok:       "success",
	ErrParam: "param error",

	ErrToken:   "invalid token",
	ErrGetUser: "get user info failed",
	ErrLogout:  "logout failed",

	ErrSceneParam:            "scene param error",
	ErrSceneID:               "scene id error",
	ErrSceneNameAlreadyExist: "scene name already exist",
	ErrEmptySceneTestCase:    "scene cases cannot be empty",
	ErrSceneCaseNameIsExist:  "scene casename already exist",
	ErrSceneDebugFailed:      "scene debug failed",
	ErrSceneCreateFailed:     "scene create failed",
	ErrSceneVariableID:       "scene variable id error",
	ErrDBParam:               "db param error",
	ErrDBID:                  "db id error",
	ErrDBMysqlFailed:         "db mysql operate failed",
	ErrSceneDeleteFailed:     "scene delete failed",

	ErrPlanParam:            "plan param error",
	ErrPlanID:               "plan id error",
	ErrPlanIDNotExist:       "plan id not exist",
	ErrPlanNameAlreadyExist: "plan name already exist",
	ErrPlanMysqlFailed:      "plan mysql operate failed",
	ErrPlanUploadFileFailed: "file upload failed",
	ErrPlanDebugFailed:      "plan debug failed",
	ErrPlanExecute:          "plan execute error",
	ErrPlanUpdate:           "plan update error",
	ErrPlanDeleteFailed:     "plan delete error",
	ErrColumnRename:         "file variable rename failed",
}
