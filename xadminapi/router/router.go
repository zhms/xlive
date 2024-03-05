package router

import (
	"xadminapi/api"
	"xcom/global"
)

func Init() {
	router := global.Router.Group("/api/v1")
	entries := api.Entries()

	entries.ApiConfig.InitRouter(router.Group("/config"))
	entries.ApiAdminSeller.InitRouter(router.Group("/seller"))

	entries.ApiAdminLog.InitRouter(router.Group("/admin_log"))
	entries.ApiAdminUser.InitRouter(router.Group("/admin_user"))
	entries.ApiAdminRole.InitRouter(router.Group("/admin_role"))

	entries.ApiUser.InitRouter(router.Group("/user_list"))

	entries.ApiLiveRoom.InitRouter(router.Group("/live_room"))
	entries.ApiLiveChat.InitRouter(router.Group("/live_chat"))
	entries.ApiLiveIpBan.InitRouter(router.Group("/live_ip_ban"))
}
