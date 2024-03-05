package controller_admin

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

type ControllerConfig struct {
	service *service_admin.ServiceAdmin
}

func (this *ControllerConfig) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.POST("/get_config", middleware.Authorization("系统管理", "系统设置", "查", ""), this.get_config)
	router.PATCH("/update_config", middleware.Authorization("系统管理", "系统设置", "改", "更新配置"), this.update_config)
}

// @Router /config/get_config [post]
// @Tags 系统设置
// @Summary 获取配置
// @Param x-token header string true "token"
// @Param body body service_admin.GetXConfigReq false "筛选参数"
// @Success 200 {object} []model.XConfig "成功"
func (this *ControllerConfig) get_config(ctx *gin.Context) {
	var reqdata service_admin.GetXConfigReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetXConfig)
}

// @Router /config/update_config [patch]
// @Tags 系统设置
// @Summary 更新配置
// @Param body body service_admin.GetXConfigReq false "筛选参数"
// @Success 200  "成功"
func (this *ControllerConfig) update_config(ctx *gin.Context) {
	var reqdata service_admin.UpdateXConfigReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateXConfig)
}
