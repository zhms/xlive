package main

import (
	"xclientapi/eventhandler"
	"xclientapi/router"
	"xclientapi/server"
	"xclientapi/service"
	"xcom/rpcclient"
)

// @title           xclientapi
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
		rpcclient.Init()
		eventhandler.Init()
		router.Init()
	})
}
