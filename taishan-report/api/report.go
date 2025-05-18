package api

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"report/internal/errno"
	"report/internal/response"
	"report/logic"
	"report/rao"
)

func CreateReport(ctx *gin.Context) {
	var req rao.CreateReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	reportId, err := logic.CreateReport(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, reportId)
	return
}

func GetReportList(ctx *gin.Context) {
	var req rao.ReportListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}

	res, total, err := logic.GetReportList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.PageResponse{
		List:  res,
		Total: total,
	})
}

func GetReportDetail(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	data, err := logic.GetReportDetail(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
	return
}

func GetReportData(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	data, err := logic.GetReportData(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
	return
}

func GetReportCaseData(ctx *gin.Context) {
	var req rao.CaseReportDataReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}

	data, err := logic.GetCaseData(ctx, req.ReportID, req.CaseID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
	return
}

func StopPressTest(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err = logic.StopPress(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportChange, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func StopPressTestHard(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err = logic.ReportPressDown(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportChange, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func UpdatePressTest(ctx *gin.Context) {
	var req rao.ConcurrencyChange
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err := logic.UpdatePressCurrency(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportChange, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func ReleasePreScene(ctx *gin.Context) {
	var req rao.ConcurrencyChange
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err := logic.ReleasePreScene(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportChange, err.Error())
		return
	}
	response.Success(ctx)
	return

}

func UpdateReportName(ctx *gin.Context) {
	var req rao.UpdateReportNameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err := logic.UpdateReportName(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportChange, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func GetReportRps(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	status, err := logic.GetReportRps(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, status)
	return
}

func GetReportTargetRps(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	rps, err := logic.GetReportTargetRps(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, rps)
	return
}
