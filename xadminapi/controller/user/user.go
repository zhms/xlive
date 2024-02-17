package controller_user

import (
	"xadminapi/service"
	service_user "xadminapi/service/user"

	"github.com/gin-gonic/gin"
)

type ControllerUser struct {
	service *service_user.ServiceUser
}

func (this *ControllerUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceUser
}
