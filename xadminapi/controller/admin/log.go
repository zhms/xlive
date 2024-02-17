package controller_admin

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_admin "xadminapi/service/admin"
	"xcom/enum"

	"github.com/gin-gonic/gin"
)

type ControllerAdminLog struct {
	service *service_admin.ServiceAdmin
}

func (this *ControllerAdminLog) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.GET("/get_login_log", middleware.Authorization("系统管理", "登录日志", "查", ""), this.get_login_log)
	router.GET("/get_opt_log", middleware.Authorization("系统管理", "操作日志", "查", ""), this.get_opt_log)
}

// @Router /admin_log/get_login_log [get]
// @Tags 后台日志
// @Summary 获取登录日志
// @Param x-token header string true "token"
// @Param query query service_admin.GetAdminLoginLogReq false "筛选参数"
// @Success 200 {object} []model.XAdminLoginLog "成功"
func (this *ControllerAdminLog) get_login_log(ctx *gin.Context) {
	var reqdata service_admin.GetAdminLoginLogReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	total, data, merr, err := this.service.GetLoginLogList(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakePageSucess(total, data))
}

// @Router /admin_log/get_opt_log [get]
// @Tags 后台日志
// @Summary 获取操作日志
// @Param query query service_admin.GetAdminOptLogReq false "筛选参数"
// @Success 200 {object} []model.XAdminOptLog "成功"
func (this *ControllerAdminLog) get_opt_log(ctx *gin.Context) {
	var reqdata service_admin.GetAdminOptLogReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	total, data, merr, err := this.service.GetOptLogList(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakePageSucess(total, data))
}
