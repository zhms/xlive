package main

import (
	"xapp/xapp"
	api_app "xclientapi/api/app"
	api_user "xclientapi/api/user"
)

// @title           xclientapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init -g main.go

func main() {
	xapp.Init()
	xapp.Run(func() {
		api_app.Init()
		api_user.Init()
	})
}
