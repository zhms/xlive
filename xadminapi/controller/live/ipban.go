package controller_live

import (
	"xadminapi/service"
	service_live "xadminapi/service/live"

	"github.com/gin-gonic/gin"
)

type ControllerLiveIpBan struct {
	service *service_live.ServiceLiveIpBan
}

func (this *ControllerLiveIpBan) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveIpBan
}
