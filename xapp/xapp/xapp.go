package xapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	mrand "math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"xapp/xdb"
	"xapp/xdb/query"
	"xapp/xglobal"
	"xapp/xredis"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var db *xdb.XDb = new(xdb.XDb)

var redis *xredis.XRedis = new(xredis.XRedis)

var db_query *query.Query

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func access_log() gin.HandlerFunc {
	return func(c *gin.Context) {
		//自定义writer
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		startTime := time.Now()

		c.Next()

		//过滤不需要记录的请求
		pattern := `(/consul|/upload|/download|/swagger)`
		re := regexp.MustCompile(pattern)
		match := re.MatchString(c.Request.URL.Path)
		if match {
			return
		}

		endTime := time.Now()
		//计算耗时
		spend := endTime.Sub(startTime)
		//获取请求数据
		reqbytes, _ := io.ReadAll(c.Request.Body)
		req := string(reqbytes)
		//响应数据
		//res := bodyLogWriter.body.String()
		//构造日志对象
		logdata := make(map[string]interface{})
		logdata["path"] = c.Request.URL.Path
		logdata["method"] = c.Request.Method
		logdata["ip"] = c.ClientIP()
		logdata["time"] = startTime.Format("2006-01-02 15:04:05")
		logdata["spend"] = spend
		//logdata["res"] = res
		logdata["status"] = c.Writer.Status()
		logdata["req"] = req
		//写日志
		bytes, _ := json.Marshal(logdata)
		bytes = append(bytes, []byte("\r\n")...)
		gin.DefaultWriter.Write(bytes)
	}
}

func cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type, x-token, Content-Length, X-Requested-With")
		context.Header("Access-Control-Allow-Methods", "GET,POST")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func Init() {
	//初始化配置
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	//读取配置
	err := viper.ReadInConfig()
	if err != nil {
		logs.Error("读取配置文件失败", err)
		return
	}
	xglobal.Id = viper.GetString("server.id")
	xglobal.Project = viper.GetString("server.project")
	xglobal.Env = viper.GetString("server.env")
	//设置随机数种子
	mrand.NewSource(time.Now().UnixNano())
	//初始化日志
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"_log/%s.log","maxsize":10485760}`, xglobal.Project))
	logs.SetLogger(logs.AdapterConsole, `{"color":true}`)

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.OpenFile("_log/http.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	gin.DefaultWriter = io.MultiWriter(f)
	xglobal.Router = gin.New()
	xglobal.Router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	xglobal.Router.Use(access_log())
	xglobal.Router.Use(cors())
	xglobal.ApiV1 = xglobal.Router.Group("/api/v1")
	xglobal.ApiV2 = xglobal.Router.Group("/api/v2")
	xglobal.ApiV3 = xglobal.Router.Group("/api/v3")
	xglobal.ApiV4 = xglobal.Router.Group("/api/v4")
	xglobal.ApiV5 = xglobal.Router.Group("/api/v5")
	xglobal.ApiV6 = xglobal.Router.Group("/api/v6")
	xglobal.ApiV7 = xglobal.Router.Group("/api/v7")
	xglobal.ApiV8 = xglobal.Router.Group("/api/v8")
	xglobal.ApiV9 = xglobal.Router.Group("/api/v9")
	//注册到consul
	if viper.GetString("consul.host") != "" {
		selfhost := viper.GetString("server.host")
		config := consul.DefaultConfig()
		config.Address = viper.GetString("consul.host") + ":" + viper.GetString("consul.port")
		consul_client, err := consul.NewClient(config)
		if err != nil {
			logs.Error("create consul client error : ", err)
			return
		}
		registration := &consul.AgentServiceRegistration{
			ID:      fmt.Sprintf("%s_%s", xglobal.Project, viper.GetString("server.id")),
			Name:    xglobal.Project,
			Port:    viper.GetInt("http.port"),
			Tags:    []string{viper.GetString("rpc.port")},
			Address: selfhost,
		}
		check := new(consul.AgentServiceCheck)
		check.HTTP = fmt.Sprintf("http://%s:%s/consul", viper.GetString("server.host"), viper.GetString("http.port"))
		check.Timeout = "1s"
		check.Interval = "2s"
		check.DeregisterCriticalServiceAfter = "1s"
		registration.Check = check
		if err := consul_client.Agent().ServiceRegister(registration); err != nil {
			logs.Error("register to consul error: ", err.Error())
			return
		}
		xglobal.Router.GET("/consul", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}
	if viper.GetString("db.user") != "" {
		db.Init("db")
		db_query = query.Use(db.Gorm())
	}
	if viper.GetString("redis.host") != "" {
		redis.Init("redis")
	}
}

func Run(callback func()) {
	{
		prcport := viper.GetInt("rpc.port")
		if prcport > 0 {
			go func() {
				listener, err := net.Listen("tcp", fmt.Sprintf(":%d", prcport))
				if err != nil {
					logs.Error("Error creating rpc listener:", err.Error())
					return
				}
				logs.Debug("start rpc server at port: ", prcport)
				for {
					conn, err := listener.Accept()
					if err != nil {
						logs.Error("Error accepting rpc connection:", err.Error())
						continue
					}
					go rpc.ServeConn(conn)
				}
			}()
		}
	}
	time.Sleep(time.Microsecond * 100)
	{
		go func() {
			port := viper.GetString("http.port")
			logs.Debug("start server at port: ", port)
			xglobal.Router.Run(":" + port)
		}()
	}
	time.Sleep(time.Microsecond * 100)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT) // ctrl+c
	signal.Notify(sig, syscall.SIGTERM)
	xglobal.Working.Add(1)
	go func() {
		for {
			select {
			case <-sig:
				fmt.Println("server exit")
				xglobal.Running = false
				xglobal.Working.Done()
			}
		}
	}()
	callback()
	xglobal.Working.Wait()
	logs.Debug("****************server exit****************")
}

func Db() *gorm.DB {
	return db.Gorm()
}

func DbQuery() *query.Query {
	return db_query
}

func Redis() *xredis.XRedis {
	return redis
}
