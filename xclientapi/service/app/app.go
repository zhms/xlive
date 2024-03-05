package service_app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"xclientapi/server"
	"xcom/edb"
	"xcom/utils"

	"github.com/gorilla/websocket"
	"github.com/yinheli/qqwry"
)

type UserData struct {
	Conn     *websocket.Conn
	ConnTime int64
	Account  string
}

type ChatData struct {
	From string `json:"from"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
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

func (this *ServiceApp) UserCome(conn *websocket.Conn, roomid string, tokendata *server.TokenData) {
	this.locker.Lock()
	defer this.locker.Unlock()

	{
		v, ok := this.maxconn[roomid]
		if !ok {
			v = 0
			this.maxconn[roomid] = v
		}
	}
	{
		this.maxconn[roomid]++
	}
	{
		v, ok := this.users[roomid]
		if !ok {
			v = make(map[string]*UserData)
			this.users[roomid] = v
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
			this.SendMsg(v[k].Conn, "user_come", tokendata.Account)
			this.SendMsg(v[k].Conn, "user_count", this.maxconn[roomid])
		}
		if len(keys) > 200 {
			this.SendMsg(conn, "user_list", keys[:200])
		} else {
			this.SendMsg(conn, "user_list", keys)
		}
		v[tokendata.Account] = &UserData{ConnTime: time.Now().Unix(), Conn: conn, Account: tokendata.Account}
		this.SendMsg(conn, "user_come", tokendata.Account)
		this.SendMsg(conn, "user_count", this.maxconn[roomid])
	}
	chardata, _ := server.Redis().Client().LRange(context.Background(), "chat_queue:"+roomid, 0, 30).Result()
	for i, j := 0, len(chardata)-1; i < j; i, j = i+1, j-1 {
		chardata[i], chardata[j] = chardata[j], chardata[i]
	}
	for _, v := range chardata {
		this.SendMsg(conn, "chat", v)
	}
}

func (this *ServiceApp) UserLeave(conn *websocket.Conn, roomid string, tokendata *server.TokenData) {
	this.locker.Lock()
	defer this.locker.Unlock()

	delete(this.users[roomid], tokendata.Account)
	this.maxconn[roomid]--
	for _, v := range this.users[roomid] {
		this.SendMsg(v.Conn, "user_leave", tokendata.Account)
		this.SendMsg(v.Conn, "user_count", this.maxconn[roomid])
	}
}

func (this *ServiceApp) ChatMsg(roomid string, tokendata *server.TokenData, msgdata string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	locker := "locker:chatmsg:" + tokendata.Account
	if !server.Redis().Lock(locker, 10) {
		for _, v := range this.users[roomid] {
			if v.Account == tokendata.Account {
				this.SendMsg(v.Conn, "chat_limit", "")
				break
			}
		}
		return
	}

	exists, err := server.Redis().Client().SIsMember(context.Background(), "ip_ban:", tokendata.Ip).Result()
	if err != nil {
		return
	}
	if exists {
		for _, v := range this.users[roomid] {
			if v.Account == tokendata.Account {
				this.SendMsg(v.Conn, "chat_ban", "")
				break
			}
		}
		return
	}

	exists, err = server.Redis().Client().SIsMember(context.Background(), "user_ban:", tokendata.UserId).Result()
	if err != nil {
		return
	}
	if exists {
		for _, v := range this.users[roomid] {
			if v.Account == tokendata.Account {
				this.SendMsg(v.Conn, "chat_ban", "")
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

	server.Db().Table(edb.TableChatList).Create(map[string]interface{}{
		edb.SellerId:   tokendata.SellerId,
		edb.UserId:     tokendata.UserId,
		edb.Ip:         tokendata.Ip,
		edb.IpLocation: iplocation,
		edb.RoomId:     roomid,
		edb.Account:    tokendata.Account,
		edb.Content:    msgdata,
		edb.State:      1,
		edb.CreateTime: utils.Now(),
	})

	chatdata := &ChatData{From: tokendata.Account, Msg: msgdata, Time: utils.Now()}
	bytes, _ := json.Marshal(chatdata)
	for _, v := range this.users[roomid] {
		if v.Account == tokendata.Account {
			this.SendMsg(v.Conn, "chat", string(bytes))
			break
		}
	}

	// server.Redis().Client().LPush(context.Background(), "chat_queue:"+roomid, string(bytes))

	// llen, _ := server.Redis().Client().LLen(context.Background(), "chat_queue:"+roomid).Result()
	// if llen > 500 {
	// 	server.Redis().Client().RPop(context.Background(), "chat_queue:"+roomid)
	// }

	// for _, v := range this.users[roomid] {

	// 	this.SendMsg(v.Conn, "chat", string(bytes))
	// }
}
