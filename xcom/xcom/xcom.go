package xcom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	mrand "math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"xcom/edb"
	"xcom/global"
	"xcom/utils"
	"xcom/xdb"
	"xcom/xredis"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var Db *xdb.XDb
var Redis *xredis.XRedis

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
		logdata := map[string]any{}
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
	global.Id = viper.GetString("server.id")
	global.Project = viper.GetString("server.project")
	global.Env = viper.GetString("server.env")
	//设置随机数种子
	mrand.NewSource(time.Now().UnixNano())
	//初始化日志
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"_log/%s.log","maxsize":10485760}`, global.Project))
	logs.SetLogger(logs.AdapterConsole, `{"color":true}`)
	// selfhost := viper.GetString("server.host")
	// if selfhost == "" {
	// 	//获取本机ip
	// 	ips := []string{}
	// 	addrs, err := net.InterfaceAddrs()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	for _, address := range addrs {
	// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 			if ipnet.IP.To4() != nil {
	// 				ips = append(ips, ipnet.IP.String())
	// 			}
	// 		}
	// 	}
	// 	selfhost = ips[0]
	// }
	// //注册到consul
	// config := consul.DefaultConfig()
	// config.Address = viper.GetString("consul.host") + ":" + viper.GetString("consul.port")
	// consul_client, err := consul.NewClient(config)
	// if err != nil {
	// 	logs.Error("create consul client error : ", err)
	// 	return
	// }
	// registration := &consul.AgentServiceRegistration{
	// 	ID:      fmt.Sprintf("%s_%s", global.Project, viper.GetString("server.id")),
	// 	Name:    global.Project,
	// 	Port:    viper.GetInt("http.port"),
	// 	Tags:    []string{viper.GetString("rpc.port")},
	// 	Address: selfhost,
	// }
	// check := &consul.AgentServiceCheck{}
	// check.HTTP = fmt.Sprintf("http://%s:%s/consul", viper.GetString("server.host"), viper.GetString("http.port"))
	// check.Timeout = "1s"
	// check.Interval = "2s"
	// check.DeregisterCriticalServiceAfter = "1s"
	// registration.Check = check
	// if err := consul_client.Agent().ServiceRegister(registration); err != nil {
	// 	logs.Error("register to consul error: ", err.Error())
	// 	return
	// }
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.OpenFile("_log/http.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	gin.DefaultWriter = io.MultiWriter(f)
	global.Router = gin.New()
	global.Router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	global.Router.Use(access_log())
	global.Router.Use(cors())
	// consul健康检查
	// global.Router.GET("/consul", func(c *gin.Context) {
	// 	c.String(200, "ok")
	// })
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
			global.Router.Run(":" + port)
		}()
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	global.WorkGroup.Add(1)
	go func() {
		for {
			select {
			case <-sig:
				global.Running = false
				global.WorkGroup.Done()
			}
		}
	}()
	callback()
	global.WorkGroup.Wait()
	logs.Debug("****************server exit****************")
}

func GetSellerId(host string) int {
	seller, err := Redis.Client().HGet(context.Background(), "host_seller", host).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("get seller id error: ", err)
		return 0
	}
	seller_id := 0
	if seller != "" {
		seller_id = utils.ToInt(seller)
		return seller_id
	}
	err = Db.Gorm().Table(edb.TableHostSeller).Where("host = ?", host).Select("seller_id").Row().Scan(&seller_id)
	if err != nil {
		return 0
	}
	if seller_id > 0 {
		Redis.Client().HSet(context.Background(), "host_seller", host, seller_id)
	}
	return seller_id
}

func NewUserId() int {
	UserId := 0
	for i := 0; i < 100; i++ {
		id := 10000000 + rand.Intn(99999999-10000000)
		rdb, _ := Db.Gorm().DB()
		//用原生sql,免得user_id重复,一大堆错误日志
		sql := fmt.Sprintf("insert into %v (user_id) values (?)", edb.TableUserIdPool)
		_, err := rdb.Exec(sql, id)
		if err == nil {
			UserId = id
			break
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return i
		}
	}
	return UserId
}
