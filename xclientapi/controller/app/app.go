package controller_app

import (
	"xclientapi/service"
	service_app "xclientapi/service/app"

	"github.com/gin-gonic/gin"
)

type ControllerApp struct {
	service *service_app.ServiceApp
}

func (this *ControllerApp) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceApp

}
