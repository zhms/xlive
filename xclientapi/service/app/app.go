package service_app

import (
	"encoding/json"
	"sync"
	"time"
	"xclientapi/server"

	"github.com/gorilla/websocket"
)

type UserData struct {
	ConnTime int64
	Conn     *websocket.Conn
}

type ServiceApp struct {
	users   map[string]map[string]*UserData
	maxconn map[string]int
	locker  sync.Mutex
}

func (this *ServiceApp) Init() {
	this.users = make(map[string]map[string]*UserData)
	this.maxconn = make(map[string]int)
}

func (this *ServiceApp) SendMsg(conn *websocket.Conn, msgid string, data interface{}) {
	msgdata := map[string]interface{}{"msg_id": msgid, "msg_data": data}
	bytes, _ := json.Marshal(msgdata)
	conn.WriteMessage(websocket.TextMessage, bytes)
}

func (this *ServiceApp) UserCome(conn *websocket.Conn, appid string, tokendata *server.TokenData) {
	this.locker.Lock()
	defer this.locker.Unlock()

	{
		v, ok := this.maxconn[appid]
		if !ok {
			v = 0
			this.maxconn[appid] = v
		}
	}
	{
		this.maxconn[appid]++
	}
	{
		v, ok := this.users[appid]
		if !ok {
			v = make(map[string]*UserData)
			this.users[appid] = v
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
			this.SendMsg(v[k].Conn, "user_come", tokendata.Account)
			this.SendMsg(v[k].Conn, "user_count", this.maxconn[appid])
		}
		this.SendMsg(conn, "user_list", keys)
		v[tokendata.Account] = &UserData{ConnTime: time.Now().Unix(), Conn: conn}
		this.SendMsg(conn, "user_come", tokendata.Account)
		this.SendMsg(conn, "user_count", this.maxconn[appid])
	}
}

func (this *ServiceApp) UserLeave(conn *websocket.Conn, appid string, tokendata *server.TokenData) {
	this.locker.Lock()
	defer this.locker.Unlock()

	delete(this.users[appid], tokendata.Account)
	this.maxconn[appid]--
	for _, v := range this.users[appid] {
		this.SendMsg(v.Conn, "user_leave", tokendata.Account)
		this.SendMsg(v.Conn, "user_count", this.maxconn[appid])
	}
}
