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

type ApiLiveIpBan struct {
	service *service_live.ServiceLiveIpBan
}

func (this *ApiLiveIpBan) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveIpBan
	router.POST("/get_ban_ip", middleware.Authorization("直播间", "Ip封禁", "查", ""), this.get_ban_ip)
	router.POST("/delete_ban_ip", middleware.Authorization("直播间", "Ip封禁", "删", ""), this.delete_ban_ip)
}

// @Router /live_ip_ban/get_ban_ip [post]
// @Tags 直播间-ip封禁
// @Summary 获取封禁ip列表
// @Param x-token header string true "token"
// @Param body body service_live.GetBanIpReq true "body参数"
// @Success 200 {object} service_live.GetBanIpRes "成功"
func (this *ApiLiveIpBan) get_ban_ip(ctx *gin.Context) {
	var reqdata service_live.GetBanIpReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetBanIp)
}

// @Router /live_ip_ban/delete_ban_ip [post]
// @Tags 直播间-ip封禁
// @Summary 解封ip
// @Param x-token header string true "token"
// @Param body body service_live.DeleteBanIpReq true "body参数"
// @Success 200   "成功"
func (this *ApiLiveIpBan) delete_ban_ip(ctx *gin.Context) {
	var reqdata service_live.DeleteBanIpReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.DeleteBanIp)
}
