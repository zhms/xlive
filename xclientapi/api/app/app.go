package app

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
	"xapp/xdb/model"
	"xapp/xglobal"
	"xapp/xutils"
	"xclientapi/api/user"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/yinheli/qqwry"
	"gorm.io/gorm"

	"github.com/gorilla/websocket"
)

var users map[string]map[string]*UserData = make(map[string]map[string]*UserData)
var maxconn map[string]int = make(map[string]int)
var locker sync.Mutex = sync.Mutex{}
var robots map[string]string = make(map[string]string)
var robot_count int = 0

func get_time_key() string {
	t := carbon.Now().TimestampMilli()
	ft := t / (5 * 60 * 1000)
	fmt.Println(t, ft*5*60*1000, carbon.CreateFromTimestampMilli(ft*5*60*1000).ToDateTimeString())
	return carbon.CreateFromTimestampMilli(ft * 5 * 60 * 1000).ToDateTimeString()
}

func Init() {
	xglobal.ApiV1.GET("/ws/:id", socket_handler)
	go audit_chat()
	go flush_robot_count()
	go flush_robot()
	go tongji()
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
	tokendata := &user.TokenData{}
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

func user_come(conn *websocket.Conn, roomid string, tokendata *user.TokenData) {
	locker.Lock()
	defer locker.Unlock()
	fmt.Println("user_come", tokendata.Account)
	{
		roomusers, ok := users[roomid]
		if !ok {
			roomusers = make(map[string]*UserData)
			users[roomid] = roomusers
		}
		keys := make([]string, 0, len(roomusers))
		for k := range roomusers {
			keys = append(keys, k)
			send_msg(roomusers[k].Conn, "user_come", tokendata.Account)
			send_msg(roomusers[k].Conn, "user_count", maxconn[roomid]+len(robots))
		}
		for k := range robots {
			if len(keys) < 100 {
				keys = append(keys, k)
			}
		}
		send_msg(conn, "user_list", keys)
		roomusers[tokendata.Account] = &UserData{ConnTime: time.Now().Unix(), Conn: conn, Account: tokendata.Account}
		send_msg(conn, "user_come", tokendata.Account)
		send_msg(conn, "user_count", maxconn[roomid]+len(robots))
	}
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
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(context.Background())
	itb = itb.Where(tb.SellerID.Eq(tokendata.SellerId))
	itb = itb.Where(tb.Account.Eq(tokendata.Account))
	itb.Update(tb.IsOnline, 1)
	allcount := 0
	for _, v := range maxconn {
		allcount += v
	}
	{
		tbs := xapp.DbQuery().XStatistic
		itbs := tbs.WithContext(context.Background())
		itbs = itbs.Where(tbs.SellerID.Eq(tokendata.SellerId))
		itbs = itbs.Where(tbs.RecordType.Eq("mx"))
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", get_time_key(), time.Local)
		itbs = itbs.Where(tbs.CreateTime.Eq(t))
		_, err := itbs.First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			xapp.DbQuery().XStatistic.Create(&model.XStatistic{
				SellerID:   tokendata.SellerId,
				RecordType: "mx",
				CreateTime: t,
				V1:         int32(allcount),
			})
		} else {
			itbs = itbs.Where(tbs.V1.Lt(int32(allcount)))
			itbs.Updates(map[string]interface{}{"v1": cast.ToString(allcount)})
		}
	}
	{
		tbs := xapp.DbQuery().XStatistic
		itbs := tbs.WithContext(context.Background())
		itbs = itbs.Where(tbs.SellerID.Eq(tokendata.SellerId))
		itbs = itbs.Where(tbs.RecordType.Eq("el"))
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", get_time_key(), time.Local)
		itbs = itbs.Where(tbs.CreateTime.Eq(t))
		_, err := itbs.First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			xapp.DbQuery().XStatistic.Create(&model.XStatistic{
				SellerID:   tokendata.SellerId,
				RecordType: "el",
				CreateTime: t,
				V1:         1,
			})
		} else {
			itbs.Updates(map[string]interface{}{"v1": gorm.Expr("v1 + 1")})
		}
	}
	chatlist, err := xapp.Redis().Client().LRange(context.Background(), "chat_queue:"+roomid, 0, 30).Result()
	if err != nil {
		return
	}
	for i := len(chatlist) - 1; i >= 0; i-- {
		send_msg(conn, "chat", chatlist[i])
	}
}

func user_leave(_ *websocket.Conn, roomid string, tokendata *user.TokenData) {
	locker.Lock()
	defer locker.Unlock()
	fmt.Println("user_leave", tokendata.Account)
	delete(users[roomid], tokendata.Account)
	maxconn[roomid]--
	for _, v := range users[roomid] {
		send_msg(v.Conn, "user_leave", tokendata.Account)
		send_msg(v.Conn, "user_count", maxconn[roomid]+len(robots))
	}
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(context.Background())
	itb = itb.Where(tb.SellerID.Eq(tokendata.SellerId))
	itb = itb.Where(tb.Account.Eq(tokendata.Account))
	itb.Update(tb.IsOnline, 2)
	{
		tbs := xapp.DbQuery().XStatistic
		itbs := tbs.WithContext(context.Background())
		itbs = itbs.Where(tbs.SellerID.Eq(tokendata.SellerId))
		itbs = itbs.Where(tbs.RecordType.Eq("el"))
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", get_time_key(), time.Local)
		itbs = itbs.Where(tbs.CreateTime.Eq(t))
		itbs.Updates(map[string]interface{}{"v2": gorm.Expr("v2 + 1")})
	}
}

func chat_msg(roomid string, tokendata *user.TokenData, msgdata string) {
	locker.Lock()
	defer locker.Unlock()
	if strings.Index(msgdata, "__hongbao__") >= 0 {
		return
	}

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
	tb := xapp.DbQuery().XChat
	itb := tb.WithContext(context.Background())
	err = itb.Create(&model.XChat{
		SellerID:   tokendata.SellerId,
		IP:         tokendata.Ip,
		IPLocation: iplocation,
		RoomID:     cast.ToInt32(roomid),
		Account:    tokendata.Account,
		Content:    msgdata,
		State:      1,
		CreateTime: carbon.Now().StdTime(),
	})
	if err != nil {
		logs.Error("Create chat data error:", err.Error())
		return
	}
	chatdata := &ChatData{From: tokendata.Account, Msg: msgdata, Time: carbon.Now().ToDateTimeString()}
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

func flush_robot_count() {
	for {
		tb := xapp.DbQuery().XKv
		itb := tb.WithContext(context.Background())
		itb = itb.Where(tb.SellerID.Eq(1))
		itb = itb.Where(tb.K.Eq("robot_count"))
		data, err := itb.First()
		if err != nil {
			logs.Error("Get robot count error:", err.Error())
			return
		}
		robot_count = cast.ToInt(data.V)
		time.Sleep(time.Second * 1)
	}
}

func flush_robot() {
	for {
		locker.Lock()
		if len(robots) >= robot_count {
			locker.Unlock()
			time.Sleep(time.Second * 1)
			continue
		}
		locker.Unlock()

		account := ""
		xapp.Db().Raw("SELECT account FROM x_robot WHERE seller_id = 1 ORDER BY RAND() LIMIT 1").Scan(&account)

		locker.Lock()
		robots[account] = account

		for k := range users["1"] {
			send_msg(users["1"][k].Conn, "user_come", account)
			send_msg(users["1"][k].Conn, "user_count", maxconn["1"]+len(robots))
		}

		locker.Unlock()
	}
}

func tongji() {
	for {
		allcount := 0
		for _, v := range maxconn {
			allcount += v
		}
		{
			tbs := xapp.DbQuery().XStatistic
			itbs := tbs.WithContext(context.Background())
			itbs = itbs.Where(tbs.SellerID.Eq(1))
			itbs = itbs.Where(tbs.RecordType.Eq("mx"))
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", get_time_key(), time.Local)
			itbs = itbs.Where(tbs.CreateTime.Eq(t))
			_, err := itbs.First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				xapp.DbQuery().XStatistic.Create(&model.XStatistic{
					SellerID:   1,
					RecordType: "mx",
					CreateTime: t,
					V1:         int32(allcount),
					V2:         0,
				})
			}
		}
		{
			tbs := xapp.DbQuery().XStatistic
			itbs := tbs.WithContext(context.Background())
			itbs = itbs.Where(tbs.SellerID.Eq(1))
			itbs = itbs.Where(tbs.RecordType.Eq("el"))
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", get_time_key(), time.Local)
			itbs = itbs.Where(tbs.CreateTime.Eq(t))
			_, err := itbs.First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				xapp.DbQuery().XStatistic.Create(&model.XStatistic{
					SellerID:   1,
					RecordType: "el",
					CreateTime: t,
					V1:         0,
					V2:         0,
				})
			}
		}
		time.Sleep(time.Second * 1)
	}
}
