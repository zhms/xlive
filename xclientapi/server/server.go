package server

import (
	"net/http"
	"xcom/enum"
	"xcom/xcom"
	"xcom/xredis"

	"xcom/xdb"

	_ "xclientapi/docs" // main 文件中导入 docs 包

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var db_conn *xdb.XDb = &xdb.XDb{}
var redis_conn *xredis.XRedis = &xredis.XRedis{}

func Init() {
	xcom.Init()
	//初始化数据库
	db_conn.Init("db")
	xcom.Db = db_conn
	//初始化redis
	redis_conn.Init("redis")
	xcom.Redis = redis_conn
}

func Run(callback func()) {
	xcom.Run(callback)
}

func Db() *gorm.DB {
	return db_conn.Gorm()
}

func Redis() *xredis.XRedis {
	return redis_conn
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func OnRequest(ctx *gin.Context, reqdata interface{}, cb func(*gin.Context, interface{}) (map[string]interface{}, error)) {
	merr, err := cb(ctx, reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.Success)
}

func OnRequestEx(ctx *gin.Context, reqdata interface{}, cb func(*gin.Context, interface{}) (interface{}, map[string]interface{}, error)) {
	data, merr, err := cb(ctx, reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakeSucess(data))
}
