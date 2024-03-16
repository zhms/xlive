package api_user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"xapp/xapp"
	"xapp/xdb"
	"xapp/xenum"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yinheli/qqwry"
	"gorm.io/gorm"
)

func Init() {
	xglobal.ApiV1.POST("/user_login", user_login)
}

type TokenData struct {
	SellerId  int
	Account   string
	IsVisitor int
	Token     string
	Ip        string
}

func DelToken(token string) {
	if token == "" {
		return
	}
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, token)
	_, err := xapp.Redis().Client().Del(context.Background(), rediskey).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func SetToken(token string, value *TokenData) {
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, token)
	valuejson, _ := json.Marshal(value)
	_, err := xapp.Redis().Client().Set(context.Background(), rediskey, string(valuejson), time.Second*3600*24*7).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func GetToken(c *gin.Context) *TokenData {
	tokenstr := c.Request.Header.Get("x-token")
	if tokenstr == "" {
		c.JSON(http.StatusBadRequest, xenum.AuthTokenNotFound)
		c.Abort()
		return nil
	}
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, tokenstr)
	value, err := xapp.Redis().Client().Get(context.Background(), rediskey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetToken error:", err.Error())
		c.JSON(http.StatusBadRequest, xenum.AuthGetTokenError)
		c.Abort()
		return nil
	}
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusBadRequest, xenum.AuthTokenNotFound)
		c.Abort()
		return nil
	}
	if value == "" {
		c.JSON(http.StatusBadRequest, xenum.AuthTokenExpired)
		c.Abort()
		return nil
	}
	tokendata := &TokenData{}
	json.Unmarshal([]byte(value), tokendata)
	return tokendata
}

type user_login_req struct {
	Account   string `validate:"required" json:"account"`
	Password  string `json:"password"`
	IsVisitor int    `json:"is_visitor"`
}

type user_login_res struct {
	Account   string `json:"account"`
	Token     string `json:"token"`
	IsVisitor int    `json:"is_visitor"`
	LiveData  string `json:"live_data"`
}

func GetLocation(ip string) string {
	ipdata := qqwry.NewQQwry("./config/qqwry.dat")
	if ipdata == nil {
		ipdata = qqwry.NewQQwry("./qqwry.dat")
	}
	if ipdata != nil && strings.Index(ip, ".") > 0 {
		ipdata.Find(ip)
		return fmt.Sprintf("%s %s", ipdata.Country, ipdata.City)
	}
	return ""
}

// @Router /user_login [post]
// @Tags 用户
// @Summary 登录
// @Param x-token header string true "token"
// @Param body body user_login_req true "请求参数"
// @Success 200  {object} user_login_res "响应数据"
func user_login(ctx *gin.Context) {
	var reqdata user_login_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	response := new(user_login_res)
	host := ctx.Request.Host
	host = strings.Replace(host, "www.", "", -1)
	host = strings.Split(host, ":")[0]
	roomid := ctx.GetHeader("roomid")
	SellerId := 1
	livingdata := xapp.Redis().Client().HGet(context.Background(), "living", fmt.Sprintf("%v_%v", SellerId, roomid)).Val()
	if len(livingdata) == 0 {
		ctx.JSON(http.StatusBadRequest, xenum.LiveNotAvailable)
		return
	}
	if reqdata.IsVisitor == xdb.StateYes {
		reqdata.Password = xutils.Md5(reqdata.Account)
	} else {
		reqdata.Password = xutils.Md5(reqdata.Password)
	}
	userdata := xdb.XUser{}
	for {
		err := xapp.Db().Where(xdb.SellerId+xdb.EQ, SellerId).Where(xdb.Account+xdb.EQ, reqdata.Account).First(&userdata).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if reqdata.IsVisitor == xdb.StateNo {
				ctx.JSON(http.StatusBadRequest, xenum.UserNotFound)
				return
			}
			xapp.Db().Model(&xdb.XUser{}).Create(map[string]interface{}{
				xdb.SellerId:   SellerId,
				xdb.Account:    reqdata.Account,
				xdb.Password:   reqdata.Password,
				xdb.IsVisitor:  reqdata.IsVisitor,
				xdb.State:      xdb.StateYes,
				xdb.LoginTime:  xutils.Now(),
				xdb.CreateTime: xutils.Now(),
			})
			continue
		}
		break
	}
	if userdata.Password != reqdata.Password {
		ctx.JSON(http.StatusBadRequest, xenum.UserPasswordError)
		return
	}
	DelToken(userdata.Token)
	userdata.Token = uuid.New().String()

	tokendata := TokenData{}
	tokendata.SellerId = SellerId
	tokendata.Account = userdata.Account
	tokendata.IsVisitor = userdata.IsVisitor
	tokendata.Token = userdata.Token
	tokendata.Ip = ctx.ClientIP()
	SetToken(tokendata.Token, &tokendata)

	err := xapp.Db().Model(&xdb.XUser{}).Where(xdb.Id+xdb.EQ, userdata.Id).Updates(map[string]interface{}{
		xdb.LoginIp:         tokendata.Ip,
		xdb.LoginTime:       xutils.Now(),
		xdb.Token:           userdata.Token,
		xdb.LoginCount:      gorm.Expr(xdb.LoginCount+xdb.PLUS, 1),
		xdb.LoginIpLocation: GetLocation(tokendata.Ip),
	}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}

	response.Account = userdata.Account
	response.Token = userdata.Token
	response.IsVisitor = userdata.IsVisitor
	response.LiveData = livingdata

	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}
