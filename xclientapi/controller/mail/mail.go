package controller_mail

import (
	"xclientapi/service"
	service_mail "xclientapi/service/mail"

	"github.com/gin-gonic/gin"
)

type ControllerMail struct {
	service *service_mail.ServiceMail
}

func (this *ControllerMail) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceMail
}
