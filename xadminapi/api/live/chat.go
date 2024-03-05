package api_live

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_live "xadminapi/service/live"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiLiveChat struct {
	service *service_live.ServiceLiveChat
}

func (this *ApiLiveChat) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveChat
	router.POST("/get_live_chat", middleware.Authorization("直播间", "互动列表", "查", ""), this.get_live_chat)
	router.POST("/audit_live_chat", middleware.Authorization("直播间", "互动列表", "改", "审核互动"), this.audit_live_chat)
}

// @Router /live_chat/get_live_chat [post]
// @Tags 直播间-互动列表
// @Summary 获取互动列表
// @Param x-token header string true "token"
// @Param body body service_live.GetChatListReq true "body参数"
// @Success 200 {object} service_live.GetChatListRes "成功"
func (this *ApiLiveChat) get_live_chat(ctx *gin.Context) {
	var reqdata service_live.GetChatListReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetChatList)
}

// @Router /live_chat/audit_live_chat [post]
// @Tags 直播间-互动列表
// @Summary 审核互动
// @Param x-token header string true "token"
// @Param body body service_live.ChatAuditReq true "body参数"
// @Success 200 "成功"
func (this *ApiLiveChat) audit_live_chat(ctx *gin.Context) {
	var reqdata service_live.ChatAuditReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.ChatAudit)
}
