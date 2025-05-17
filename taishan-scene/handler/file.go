package handler

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
	"strconv"
)

func FileUpload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanUploadFileFailed, "file参数解析失败")
		return
	}
	planIdStr := ctx.PostForm("plan_id")
	planId, err := strconv.ParseInt(planIdStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}

	running, _ := logic.PlanIsRunning(ctx, int32(planId))
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	err = logic.UploadFile(ctx, int32(planId), *file)

	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanUploadFileFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func DownloadPlanFile(ctx *gin.Context) {
	var req rao.FileDownload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, errno.CodeAlertMap[errno.ErrParam])
		return
	}
	content, err := logic.DownloadFile(ctx, filepath.Join(req.Path, req.FileName))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrDownloadFileFailed, err.Error())
		return
	}
	fileContentDisposition := "attachment;filename=\"" + req.FileName + "\""
	ctx.Header("Content-Disposition", fileContentDisposition)
	ctx.Data(http.StatusOK, "", content)
	return
}

// GetPlanDataSource 参数化文件列表
func GetPlanDataSource(ctx *gin.Context) {
	planId, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}

	res, err := logic.GetPlanDataSource(ctx, int32(planId))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, res)
	return
}

func ColumnUpdate(ctx *gin.Context) {
	var req rao.ColumnUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, errno.CodeAlertMap[errno.ErrPlanParam])
		return
	}
	err := logic.ColumnUpdate(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrColumnRename, errno.CodeAlertMap[errno.ErrColumnRename])
	}
	response.Success(ctx)
	return
}

func DeleteFile(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}

	b, err := logic.DeleteParameterFile(ctx, int32(id))
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}
