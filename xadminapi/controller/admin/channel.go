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

type ControllerChannel struct {
	service *service_admin.ServiceAdmin
}

func (this *ControllerChannel) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.GET("/get_channel", this.get_channel)
	router.POST("/create_channel", middleware.Authorization("系统管理", "渠道管理", "增", "添加渠道"), this.create_channel)
	router.PATCH("/update_channel", middleware.Authorization("系统管理", "渠道管理", "改", "更新渠道"), this.update_channel)
	router.DELETE("/delete_channel", middleware.Authorization("系统管理", "渠道管理", "删", "删除渠道"), this.delete_channel)
}

// @Router /channel/get_channel [get]
// @Tags 渠道
// @Summary 获取渠道列表
// @Param x-token header string true "token"
// @Param query query service_admin.GetXChannelReq false "筛选参数"
// @Success 200 {object} []model.XChannel "成功"
func (this *ControllerChannel) get_channel(ctx *gin.Context) {
	var reqdata service_admin.GetXChannelReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	total, data, merr, err := this.service.GetChannelList(&reqdata)
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

// @Router /channel/create_channel [post]
// @Tags 渠道
// @Summary 新增渠道
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.CreateXChannelReq true "body参数"
// @Success 200 "成功"
func (this *ControllerChannel) create_channel(ctx *gin.Context) {
	var reqdata service_admin.CreateXChannelReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	merr, err := this.service.CreateChannel(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.Success)
}

// @Router /channel/update_channel [patch]
// @Tags 渠道
// @Summary 更新渠道
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.UpdateXChannelReq true "body参数"
// @Success 200 "成功"
func (this *ControllerChannel) update_channel(ctx *gin.Context) {
	var reqdata service_admin.UpdateXChannelReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	merr, err := this.service.UpdateChannel(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.Success)
}

// @Router /channel/delete_channel [delete]
// @Tags 渠道
// @Summary 删除渠道
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.DeleteXChannelReq true "body参数"
// @Success 200 "成功"
func (this *ControllerChannel) delete_channel(ctx *gin.Context) {
	var reqdata service_admin.DeleteXChannelReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	rows, merr, err := this.service.DeleteChannel(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakeSucess(rows))
}
