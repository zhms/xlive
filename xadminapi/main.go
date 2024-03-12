package main

import (
	"embed"
	api_admin "xadminapi/api/admin"
	api_live_room "xadminapi/api/live/room"
	_ "xadminapi/docs"
	"xapp/xapp"
)

// @title          adminapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init -g main.go
//
//go:embed www
var www embed.FS

func main() {
	xapp.Init()
	xapp.Run(func() {
		api_admin.Init(&www)
		api_live_room.Init()
	})
}
