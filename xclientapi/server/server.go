package server

import (
	"xcom/cuser"
	"xcom/enum"
	"xcom/global"
	"xcom/mq"
	"xcom/statistic"
	"xcom/xcom"
	"xcom/xredis"

	"xcom/xdb"

	_ "xclientapi/docs" // main 文件中导入 docs 包

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var db_game *xdb.XDb = &xdb.XDb{}
var db_user *xdb.XDb = &xdb.XDb{}
var db_order *xdb.XDb = &xdb.XDb{}
var db_statistic *xdb.XDb = &xdb.XDb{}

var redis_token *xredis.XRedis = &xredis.XRedis{}
var redis_game *xredis.XRedis = &xredis.XRedis{}
var redis_user *xredis.XRedis = &xredis.XRedis{}
var redis_order *xredis.XRedis = &xredis.XRedis{}
var redis_statistic *xredis.XRedis = &xredis.XRedis{}

var rabbitmq *mq.RabbitMq = &mq.RabbitMq{}

func Init() {
	xcom.Init()
	//初始化数据库
	db_game.Init("db_game")
	db_user.Init("db_user")
	db_order.Init("db_order")
	db_statistic.Init("db_statistic")
	xcom.DbGame = db_game
	xcom.DbUser = db_user
	xcom.DbOrder = db_order
	xcom.DbStatistic = db_statistic
	//初始化redis
	redis_token.Init("redis_token")
	redis_game.Init("redis_game")
	redis_user.Init("redis_user")
	redis_order.Init("redis_order")
	redis_statistic.Init("redis_statistic")
	xcom.RedisGame = redis_game
	xcom.RedisUser = redis_user
	xcom.RedisOrder = redis_order
	xcom.RedisStatistic = redis_statistic
	//初始化消息队列
	rabbitmq.Init("rabbitmq")
	//初始化玩家缓存
	cuser.Init()
	statistic.Init()
	//测试页面
	global.Router.GET("/game.html", func(c *gin.Context) {
		c.File("./html/game.html")
	})
}

func Run(callback func()) {
	xcom.Run(callback)
}

func Db(id int) *gorm.DB {
	switch id {
	case enum.DbGameRW:
		return db_game.Gorm()
	case enum.DbUserRW:
		return db_user.Gorm()
	case enum.DbOrderRW:
		return db_order.Gorm()
	case enum.DbStatisticRW:
		return db_statistic.Gorm()
	}
	return nil
}

func Redis(id int) *xredis.XRedis {
	switch id {
	case enum.RedisTokenRw:
		return redis_token
	case enum.RedisGameRw:
		return redis_game
	case enum.RedisUserRw:
		return redis_user
	case enum.RedisOrderRw:
		return redis_order
	case enum.RedisStatisticRw:
		return redis_statistic
	}
	return nil
}

func RabbitMq() *mq.RabbitMq {
	return rabbitmq
}
