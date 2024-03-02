package controller_live

import (
	"xadminapi/service"
	service_live "xadminapi/service/live"

	"github.com/gin-gonic/gin"
)

type ControllerLiveChat struct {
	service *service_live.ServiceLiveChat
}

func (this *ControllerLiveChat) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveChat
}

func (this *ControllerLiveChat) GetChatList(c *gin.Context) {

}
