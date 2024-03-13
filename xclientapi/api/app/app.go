package api_app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"xapp/xapp"
	"xapp/xdb"
	"xapp/xglobal"
	"xapp/xutils"
	api_user "xclientapi/api/user"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yinheli/qqwry"

	"github.com/gorilla/websocket"
)

func Init() {
	xglobal.ApiV1.GET("/ws/:id", socket_handler)
	go audit_chat()
}

type UserData struct {
	ConnTime int64
	Conn     *websocket.Conn
	Account  string
}

type ChatData struct {
	From string `json:"from"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}

var users map[string]map[string]*UserData = make(map[string]map[string]*UserData)
var maxconn map[string]int = make(map[string]int)
var locker sync.Mutex = sync.Mutex{}

func send_msg(conn *websocket.Conn, msgid string, data interface{}) {
	msgdata := map[string]interface{}{"msg_id": msgid, "msg_data": data}
	bytes, _ := json.Marshal(msgdata)
	conn.WriteMessage(websocket.TextMessage, bytes)
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MsgData struct {
	MsgId   string                 `json:"msg_id"`
	MsgData map[string]interface{} `json:"msg_data"`
}

func socket_handler(ctx *gin.Context) {
	conn, err := WsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
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
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, ids[0])
	value, err := xapp.Redis().Client().Get(context.Background(), rediskey).Result()
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
	tokendata := &api_user.TokenData{}
	json.Unmarshal([]byte(value), tokendata)
	user_come(conn, ids[1], tokendata)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if string(msg) == "ping" {
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}
		msgdata := &MsgData{}
		err = json.Unmarshal(msg, msgdata)
		if err != nil {
			break
		}
		if msgdata.MsgId == "chat" {
			chat_msg(ids[1], tokendata, msgdata.MsgData["msg"].(string))
		}
	}
	user_leave(conn, ids[1], tokendata)
}

func user_come(conn *websocket.Conn, roomid string, tokendata *api_user.TokenData) {
	locker.Lock()
	defer locker.Unlock()

	{
		v, ok := maxconn[roomid]
		if !ok {
			v = 0
			maxconn[roomid] = v
		}
	}
	{
		maxconn[roomid]++
	}
	{
		v, ok := users[roomid]
		if !ok {
			v = make(map[string]*UserData)
			users[roomid] = v
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
			send_msg(v[k].Conn, "user_come", tokendata.Account)
			send_msg(v[k].Conn, "user_count", maxconn[roomid])
		}
		send_msg(conn, "user_list", keys)
		v[tokendata.Account] = &UserData{ConnTime: time.Now().Unix(), Conn: conn, Account: tokendata.Account}
		send_msg(conn, "user_come", tokendata.Account)
		send_msg(conn, "user_count", maxconn[roomid])
	}
}

func user_leave(_ *websocket.Conn, roomid string, tokendata *api_user.TokenData) {
	locker.Lock()
	defer locker.Unlock()

	delete(users[roomid], tokendata.Account)
	maxconn[roomid]--
	for _, v := range users[roomid] {
		send_msg(v.Conn, "user_leave", tokendata.Account)
		send_msg(v.Conn, "user_count", maxconn[roomid])
	}
}

func chat_msg(roomid string, tokendata *api_user.TokenData, msgdata string) {
	locker.Lock()
	defer locker.Unlock()

	// locker := "locker:chatmsg:" + tokendata.Account
	// if !server.Redis().Lock(locker, 10) {
	// 	for _, v := range this.users[roomid] {
	// 		if v.Account == tokendata.Account {
	// 			this.SendMsg(v.Conn, "chat_limit", "")
	// 			break
	// 		}
	// 	}
	// 	return
	// }

	exists, err := xapp.Redis().Client().SIsMember(context.Background(), "ip_ban", tokendata.Ip).Result()
	if err != nil {
		return
	}
	if exists {
		for _, v := range users[roomid] {
			if v.Account == tokendata.Account {
				send_msg(v.Conn, "chat_ban", "")
				break
			}
		}
		return
	}

	exists, err = xapp.Redis().Client().SIsMember(context.Background(), "user_ban", fmt.Sprintf("%d_%s", tokendata.SellerId, tokendata.Account)).Result()
	if err != nil {
		return
	}
	if exists {
		for _, v := range users[roomid] {
			if v.Account == tokendata.Account {
				send_msg(v.Conn, "chat_ban", "")
				break
			}
		}
		return
	}

	iplocation := ""
	ipdata := qqwry.NewQQwry("./config/ipdata.dat")
	if ipdata != nil && strings.Index(tokendata.Ip, ".") > 0 {
		ipdata.Find(tokendata.Ip)
		iplocation = fmt.Sprintf("%s %s", ipdata.Country, ipdata.City)
	}

	db := xapp.Db().Model(&xdb.XChatData{}).Create(map[string]interface{}{
		xdb.SellerId:   tokendata.SellerId,
		xdb.Ip:         tokendata.Ip,
		xdb.IpLocation: iplocation,
		xdb.RoomId:     roomid,
		xdb.Account:    tokendata.Account,
		xdb.Content:    msgdata,
		xdb.State:      1,
		xdb.CreateTime: xutils.Now(),
	})
	if db.Error != nil {
		logs.Error("Create chat data error:", db.Error.Error())
		return
	}

	chatdata := &ChatData{From: tokendata.Account, Msg: msgdata, Time: xutils.Now()}
	bytes, _ := json.Marshal(chatdata)
	for _, v := range users[roomid] {
		if v.Account == tokendata.Account {
			send_msg(v.Conn, "chat", string(bytes))
			break
		}
	}
}

func audit_chat() {
	for {
		chatstr, err := xapp.Redis().Client().LPop(context.Background(), "chat_audit").Result()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		chatdata := xutils.XMap{}
		json.Unmarshal([]byte(chatstr), &chatdata.RawData)

		roomid := chatdata.String(xdb.RoomId)

		chatmsg := &ChatData{From: chatdata.String(xdb.Account), Msg: chatdata.String(xdb.Content), Time: chatdata.String(xdb.CreateTime)}

		bytes, _ := json.Marshal(chatmsg)

		xapp.Redis().Client().LPush(context.Background(), "chat_queue:"+roomid, string(bytes))
		llen, _ := xapp.Redis().Client().LLen(context.Background(), "chat_queue:"+roomid).Result()
		if llen > 200 {
			xapp.Redis().Client().RPop(context.Background(), "chat_queue:"+roomid)
		}
		for _, v := range users[roomid] {
			if chatdata.String(xdb.Account) != v.Account {
				send_msg(v.Conn, "chat", string(bytes))
			}
		}
	}
}
