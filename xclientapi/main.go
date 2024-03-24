package main

import (
	"math/rand"
	"time"
	"xapp/xapp"
	"xclientapi/api/app"
	"xclientapi/api/user"
)

// @title           xclientapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init -g main.go

func main() {
	rand.Seed(time.Now().UnixNano())
	xapp.Init()
	xapp.Run(func() {
		app.Init()
		user.Init()
	})
}
