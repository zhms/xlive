package main

import (
	"embed"
	"xadminapi/api/admin"
	"xadminapi/api/data/inout_chart"
	"xadminapi/api/data/online_chart"
	"xadminapi/api/hongbao"
	live_ban "xadminapi/api/live/ban"
	live_chat "xadminapi/api/live/chat"
	live_room "xadminapi/api/live/room"
	"xadminapi/api/robot"
	"xadminapi/api/sales"

	"xadminapi/api/user"
	_ "xadminapi/docs"

	"xapp/xapp"
)

// @title          adminapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init --parseDependency -g main.go

//go:embed www
var www embed.FS

func main() {
	xapp.Init()
	xapp.Run(func() {
		admin.Init(&www)
		user.Init()
		live_room.Init()
		live_ban.Init()
		live_chat.Init()
		sales.Init()
		hongbao.Init()
		robot.Init()
		inout_chart.Init()
		online_chart.Init()
	})
}
