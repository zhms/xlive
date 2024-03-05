package api_admin

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_admin "xadminapi/service/admin"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiAdminLog struct {
	service *service_admin.ServiceAdmin
}

func (this *ApiAdminLog) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.POST("/get_login_log", middleware.Authorization("系统管理", "登录日志", "查", ""), this.get_login_log)
	router.POST("/get_opt_log", middleware.Authorization("系统管理", "操作日志", "查", ""), this.get_opt_log)
}

// @Router /admin_log/get_login_log [post]
// @Tags 后台日志
// @Summary 获取登录日志
// @Param x-token header string true "token"
// @Param body body service_admin.GetAdminLoginLogReq false "筛选参数"
// @Success 200 {object} service_admin.GetAdminLoginLogRes "成功"
func (this *ApiAdminLog) get_login_log(ctx *gin.Context) {
	var reqdata service_admin.GetAdminLoginLogReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetLoginLogList)
}

// @Router /admin_log/get_opt_log [post]
// @Tags 后台日志
// @Summary 获取操作日志
// @Param x-token header string true "token"
// @Param body body service_admin.GetAdminOptLogReq false "筛选参数"
// @Success 200 {object} service_admin.GetAdminOptLogRes "成功"
func (this *ApiAdminLog) get_opt_log(ctx *gin.Context) {
	var reqdata service_admin.GetAdminOptLogReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetOptLogList)
}
