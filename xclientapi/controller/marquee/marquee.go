package controller_marquee

import (
	"xclientapi/service"
	service_marquee "xclientapi/service/marquee"

	"github.com/gin-gonic/gin"
)

type ControllerMarquee struct {
	service *service_marquee.ServiceMarquee
}

func (this *ControllerMarquee) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceMarquee
}
