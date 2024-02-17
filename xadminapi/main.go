package main

import (
	"xadminapi/router"
	"xadminapi/server"
	"xadminapi/service"
)

// @title          adminapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init -g main.go

func main() {
	server.Init()
	server.Run(func() {
		service.Init()
		router.Init()
	})
}
