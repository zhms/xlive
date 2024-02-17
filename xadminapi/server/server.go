package server

import (
	"xcom/xcom"
	"xcom/xredis"

	"xcom/xdb"

	_ "xadminapi/docs" // main 文件中导入 docs 包

	_ "github.com/go-sql-driver/mysql"
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
