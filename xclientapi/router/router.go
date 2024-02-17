package router

import (
	"xclientapi/controller"
	"xcom/global"
)

func Init() {
	router := global.Router.Group("/api/v1")
	entries := controller.Entries()
	entries.ControllerApp.InitRouter(router.Group("app"))
	entries.ControllerUser.InitRouter(router.Group("user"))
}
