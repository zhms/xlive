package server

import (
	"net/http"
	"xcom/xcom"

	"xcom/xredis"

	"github.com/gorilla/websocket"
)

var redis *xredis.XRedis = &xredis.XRedis{}

func Init() {
	xcom.Init()
	redis.Init("redis")
	xcom.Redis = redis
}

func Run(callback func()) {
	xcom.Run(callback)
}

func Redis() *xredis.XRedis {
	return redis
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
