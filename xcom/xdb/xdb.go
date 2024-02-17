package xdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"xcom/global"

	"github.com/beego/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type XDb struct {
	user            string
	password        string
	host            string
	port            int
	connmaxlifetime int
	database        string
	db              *gorm.DB
	connmaxidletime int
	connmaxidle     int
	connmaxopen     int
}

// 初始化db
func (this *XDb) Init(cfgname string) {
	this.user = viper.GetString(fmt.Sprint(cfgname, ".user"))
	this.password = viper.GetString(fmt.Sprint(cfgname, ".password"))
	this.host = viper.GetString(fmt.Sprint(cfgname, ".host"))
	this.database = viper.GetString(fmt.Sprint(cfgname, ".database"))
	this.port = viper.GetInt(fmt.Sprint(cfgname, ".port"))

	this.connmaxlifetime = 60
	this.connmaxidletime = 600
	this.connmaxidle = 2
	this.connmaxopen = 5

	if strings.Contains(global.Env, "prd") {
		this.connmaxlifetime = 600
		this.connmaxidletime = 6000
		this.connmaxidle = 20
		this.connmaxopen = 100
	}

	connmaxlifetime := viper.GetInt(fmt.Sprint(cfgname, ".connmaxlifetime"))
	if connmaxlifetime > 0 {
		this.connmaxlifetime = connmaxlifetime
	}
	connmaxidletime := viper.GetInt(fmt.Sprint(cfgname, ".connmaxidletime"))
	if connmaxidletime > 0 {
		this.connmaxidletime = connmaxidletime
	}
	connmaxidle := viper.GetInt(fmt.Sprint(cfgname, ".connmaxidle"))
	if connmaxidle > 0 {
		this.connmaxidle = connmaxidle
	}
	connmaxopen := viper.GetInt(fmt.Sprint(cfgname, ".connmaxopen"))
	if connmaxopen > 0 {
		this.connmaxopen = connmaxopen
	}

	conurl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", this.user, this.password, this.host, this.port, this.database)
	f, _ := os.OpenFile(fmt.Sprintf("_log/gorm_%s.log", cfgname), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	gormcfg := &gorm.Config{
		Logger: logger.New(log.New(f, "", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	}
	db, err := gorm.Open(mysql.Open(conurl), gormcfg)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	gdb, _ := db.DB()
	gdb.SetConnMaxLifetime(time.Second * time.Duration(this.connmaxlifetime))
	gdb.SetConnMaxIdleTime(time.Second * time.Duration(this.connmaxidletime))
	gdb.SetMaxIdleConns(this.connmaxidle)
	gdb.SetMaxOpenConns(this.connmaxopen)
	this.db = db
	logs.Debug("连接数据库成功:", this.host, this.port, this.database)
}

func (this *XDb) Gorm() *gorm.DB {
	return this.db
}

// 获取连接的数据库
func (this *XDb) Database() string {
	return this.database
}
