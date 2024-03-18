package main

import (
	"context"
	"embed"
	"encoding/json"
	"xadminapi/api/admin"
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

		tb := xapp.DbQuery().XTest
		itb := tb.WithContext(context.Background())
		itb = itb.Where(tb.A.Eq("1"))
		itb = itb.Where(tb.A.Eq("3"))
		a, _ := itb.First()
		b, _ := json.Marshal(a)
		println(string(b))
	})
}
