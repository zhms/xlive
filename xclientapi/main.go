package main

import (
	"embed"
	"math/rand"
	"time"
	"xapp/xapp"
	"xclientapi/api/app"
	"xclientapi/api/page"
	"xclientapi/api/user"
)

// @title           xclientapi
// @version         1.0

// @BasePath  /api/v1

// @securityDefinitions.apiKey
// @in header
// @name x-token

// swag init -g main.go

//go:embed www
var www embed.FS

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	xapp.Init()
	xapp.Run(func() {
		page.Init(&www)
		app.Init()
		user.Init()
	})
}
