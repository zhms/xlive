package router

import (
	"xadminapi/controller"
	"xcom/global"
)

func Init() {
	router := global.Router.Group("/api/v1")
	entries := controller.Entries()

	entries.ControllerConfig.InitRouter(router.Group("/config"))
	entries.ControllerAdminSeller.InitRouter(router.Group("/seller"))

	entries.ControllerAdminLog.InitRouter(router.Group("/admin_log"))
	entries.ControllerAdminUser.InitRouter(router.Group("/admin_user"))
	entries.ControllerAdminRole.InitRouter(router.Group("/admin_role"))

	entries.ControllerUser.InitRouter(router.Group("/user_list"))

	entries.ControllerLiveRoom.InitRouter(router.Group("/live_room"))
	entries.ControllerLiveChat.InitRouter(router.Group("/live_chat"))
	entries.ControllerLiveIpBan.InitRouter(router.Group("/live_ip_ban"))
}
