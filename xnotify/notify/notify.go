package notify

import (
	"context"
	"encoding/json"
	"sync"
	"xcom/global"
	"xcom/utils"
	"xnotify/model"
	"xnotify/server"

	"github.com/beego/beego/logs"
	"github.com/gorilla/websocket"
)

var users sync.Map
var user_count int = 0

var connection_subscribes sync.Map
var subscribes sync.Map

var last_block_data []byte

type MsgData struct {
	MsgType string `json:"msg_type"`
	MsgId   string `json:"msg_id"`
	Topic   string `json:"topic"`
}

func Init() {
	server.Redis().Client().HSet(context.Background(), "online:"+global.Project, global.Id, 0)
	global.Router.GET("/ws/:id", socket_handler)
}

func on_queue_message(data []byte) error {
	msgdata := map[string]interface{}{}
	err := json.Unmarshal(data, &msgdata)
	if err != nil {
		logs.Error("on_queue_message json.Unmarshal error:", err)
		return err
	}
	msgtype := utils.ToString(msgdata["msg_type"])
	if msgtype == "topic" {
		id := utils.ToString(msgdata["id"])
		d, ok := msgdata["data"]
		if !ok {
			d = map[string]interface{}{}
		}
		last_block_data, _ = json.Marshal(map[string]interface{}{
			"msg_id": id,
			"data":   d,
		})
		if connections, ok := subscribes.Load(id); ok {
			connections.(*sync.Map).Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				conn.WriteMessage(websocket.TextMessage, last_block_data)
				return true
			})
		}
	} else {
		id := utils.ToString(msgdata["id"])
		if id == "all" {
			users.Range(func(key, value interface{}) bool {
				conn := value.(*websocket.Conn)
				d, ok := msgdata["data"]
				if !ok {
					d = msgdata
				} else {
					d = d.(map[string]interface{})
				}
				bytes, _ := json.Marshal(d)
				conn.WriteMessage(websocket.TextMessage, bytes)
				return true
			})
		} else {
			conn, ok := users.Load(id)
			if ok {
				d, ok := msgdata["data"]
				if !ok {
					d = msgdata
				} else {
					d = d.(map[string]interface{})
				}
				bytes, _ := json.Marshal(d)
				conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, bytes)
			}
		}
	}
	return nil
}

func on_block_data(data []byte) error {
	blockinfo := model.BlockInfo{}
	err := json.Unmarshal(data, &blockinfo)
	if err != nil {
		logs.Error("on_block_data json.Unmarshal err", err)
		return err
	}
	last_block_data, _ = json.Marshal(map[string]interface{}{
		"msg_id": "block_info",
		"data":   blockinfo,
	})
	if connections, ok := subscribes.Load("block_info"); ok {
		connections.(*sync.Map).Range(func(key, value interface{}) bool {
			conn := key.(*websocket.Conn)
			conn.WriteMessage(websocket.TextMessage, last_block_data)
			return true
		})
	}
	return nil
}
