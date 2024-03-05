package api_live

import (
	"xadminapi/service"
	service_live "xadminapi/service/live"

	"github.com/gin-gonic/gin"
)

type ApiLiveIpBan struct {
	service *service_live.ServiceLiveIpBan
}

func (this *ApiLiveIpBan) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveIpBan
}
