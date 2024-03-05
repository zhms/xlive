package router

import (
	"xclientapi/api"
	"xcom/global"
)

func Init() {
	router := global.Router.Group("/api/v1")
	entries := api.Entries()
	entries.ApiApp.InitRouter(router.Group("app"))
	entries.ApiUser.InitRouter(router.Group("user"))
}
