package xdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"xapp/xutils"

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
	env             string
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

	this.env = viper.GetString("server.env")

	if strings.Contains(this.env, "prd") {
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

func (this *XDb) Database() string {
	return this.database
}

func dbgetone(rows *sql.Rows) *map[string]interface{} {
	data := make(map[string]interface{})
	fields, _ := rows.Columns()
	scans := make([]interface{}, len(fields))
	for i := range scans {
		scans[i] = &scans[i]
	}
	err := rows.Scan(scans...)
	if err != nil {
		logs.Error(err)
		return nil
	}
	for i := range fields {
		if scans[i] != nil {
			data[fields[i]] = xutils.ToString(scans[i])
		} else {
			data[fields[i]] = nil
		}
	}
	return &data
}

func First(db *gorm.DB) (*xutils.XMap, error) {
	rows, err := db.Rows()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	data := new(xutils.XMap)
	if rows.Next() {
		data.RawData = *dbgetone(rows)
	} else {
		return nil, nil
	}
	rows.Close()
	return data, nil
}

func Find(db *gorm.DB) (*xutils.XMaps, error) {
	rows, err := db.Rows()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	data := new(xutils.XMaps)
	for rows.Next() {
		data.RawData = append(data.RawData, *dbgetone(rows))
	}
	rows.Close()
	return data, nil
}

const (
	DESC      = " desc "
	ASC       = " asc "
	EQ        = " = ? "
	NEQ       = " <> ? "
	AND       = " and "
	OR        = " or "
	ISNULL    = " is null "
	ISNOTNULL = " is not null "
	LIKE      = " like ? "
	IN        = " in (?) "
	NOTIN     = " not in (?) "
	GT        = " > ? "
	GTE       = " >= ? "
	LT        = " < ? "
	LTE       = " <= ? "
	BETWEEN   = " between ? and ? "
	PLUS      = " + ? "
	MINUS     = " - ? "
)
