package router

import (
	"xadminapi/controller"
	"xcom/global"
)

func Init() {
	router := global.Router.Group("/api/v1")
	entries := controller.Entries()
	entries.ControllerAdminUser.InitRouter(router.Group("/admin_user"))
	entries.ControllerAdminRole.InitRouter(router.Group("/admin_role"))
	entries.ControllerAdminLog.InitRouter(router.Group("/admin_log"))
	entries.ControllerAdminSeller.InitRouter(router.Group("/seller"))
	entries.ControllerConfig.InitRouter(router.Group("/config"))
	entries.ControllerUser.InitRouter(router.Group("/user"))
}
