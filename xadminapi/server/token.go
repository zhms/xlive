package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"xcom/enum"
	"xcom/global"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TokenData struct {
	SellerId     int
	Account      string
	UserId       int
	AuthData     string
	GoogleSecret string
}

func DelToken(token string) {
	if token == "" {
		return
	}
	rediskey := fmt.Sprintf("%v:token:%s", global.Project, token)
	_, err := redis_conn.Client().Del(context.Background(), rediskey).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func SetToken(token string, value *TokenData) {
	rediskey := fmt.Sprintf("%v:token:%s", global.Project, token)
	valuejson, _ := json.Marshal(value)
	_, err := redis_conn.Client().Set(context.Background(), rediskey, string(valuejson), time.Second*3600*24*7).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func GetToken(c *gin.Context) *TokenData {
	tokenstr := c.Request.Header.Get("x-token")
	if tokenstr == "" {
		c.JSON(http.StatusBadRequest, enum.AuthTokenNotFound)
		c.Abort()
		return nil
	}
	rediskey := fmt.Sprintf("%v:token:%s", global.Project, tokenstr)
	value, err := redis_conn.Client().Get(context.Background(), rediskey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logs.Error("GetToken error:", err.Error())
		}
		c.JSON(http.StatusBadRequest, enum.AuthGetTokenError)
		c.Abort()
		return nil
	}
	if value == "" {
		c.JSON(http.StatusBadRequest, enum.AuthTokenExpired)
		c.Abort()
		return nil
	}
	tokendata := &TokenData{}
	json.Unmarshal([]byte(value), tokendata)
	return tokendata
}
