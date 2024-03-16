package main

import (
	"context"
	"embed"
	"encoding/json"
	"xadminapi/api/admin"
	live_ban "xadminapi/api/live/ban"
	live_chat "xadminapi/api/live/chat"
	live_room "xadminapi/api/live/room"

	"xadminapi/api/user"
	_ "xadminapi/docs"

	"xapp/xapp"

	"xapp/xdb/query"

	"github.com/beego/beego/logs"
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
var daoXHashGame *query.Query

func main() {
	xapp.Init()
	daoXHashGame = query.Use(xapp.Db())
	xapp.Run(func() {
		admin.Init(&www)
		user.Init()
		live_room.Init()
		live_ban.Init()
		live_chat.Init()
		tb := daoXHashGame.XKv
		itb := tb.WithContext(context.Background())
		data, count, _ := itb.FindByPage(0, 2)
		bytes, _ := json.Marshal(data)
		logs.Debug("data", count, string(bytes))
	})
}
