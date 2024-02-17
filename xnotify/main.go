package main

import (
	"xnotify/notify"
	"xnotify/server"
)

func main() {
	server.Init()
	server.Run(func() {
		notify.Init()
	})
}
