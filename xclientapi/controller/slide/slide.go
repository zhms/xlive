package controller_slide

import (
	"xclientapi/service"
	service_slide "xclientapi/service/slide"

	"github.com/gin-gonic/gin"
)

type ControllerSlide struct {
	service *service_slide.ServiceSlide
}

func (this *ControllerSlide) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceSlide
}
