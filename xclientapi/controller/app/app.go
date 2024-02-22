package controller_app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"xclientapi/server"
	"xclientapi/service"
	service_app "xclientapi/service/app"
	"xcom/global"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type ControllerApp struct {
	service *service_app.ServiceApp
}

func (this *ControllerApp) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceApp
	router.GET("/ws/:id", this.socket_handler)
}

type MsgData struct {
	MsgId   string                 `json:"msg_id"`
	MsgData map[string]interface{} `json:"msg_data"`
}

func (this *ControllerApp) socket_handler(ctx *gin.Context) {
	conn, err := server.WsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logs.Error("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	id := ctx.Param("id")
	if id == "" {
		conn.Close()
		return
	}
	ids := strings.Split(id, "_")
	if len(ids) != 2 {
		conn.Close()
		return
	}
	host := ctx.Request.Host
	host = strings.Replace(host, "www.", "", -1)
	host = strings.Split(host, ":")[0]
	rediskey := fmt.Sprintf("%v:token:%s", global.Project, ids[0])
	value, err := server.Redis().Client().Get(context.Background(), rediskey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetToken error:", err.Error())
		conn.Close()
		return
	}
	if errors.Is(err, redis.Nil) {
		conn.Close()
		return
	}
	if value == "" {
		conn.Close()
		return
	}
	tokendata := &server.TokenData{}
	json.Unmarshal([]byte(value), tokendata)
	this.service.UserCome(conn, ids[1], tokendata)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if string(msg) == "ping" {
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}
		break
	}
}
