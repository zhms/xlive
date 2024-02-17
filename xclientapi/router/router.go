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
	entries.ControllerMail.InitRouter(router.Group("mail"))
	entries.ControllerMarquee.InitRouter(router.Group("marquee"))
	entries.ControllerSlide.InitRouter(router.Group("slide"))
	entries.ControllerHash.InitRouter(router.Group("hash"))
}
