package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yinheli/qqwry"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

var locker sync.Mutex

func Init() {
	xglobal.ApiV1.POST("/user_login", user_login)
	xglobal.ApiV1.POST("/open_hongbao", open_hongbao)
}

type TokenData struct {
	SellerId  int32
	Account   string
	IsVisitor int32
	Token     string
	Ip        string
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
	SaleId    string `json:"sale_id"`
}

type user_login_res struct {
	Account   string `json:"account"`
	Token     string `json:"token"`
	IsVisitor int32  `json:"is_visitor"`
	LiveData  string `json:"live_data"`
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
	SellerId := int32(1)
	livingdata := xapp.Redis().Client().HGet(context.Background(), "living", fmt.Sprintf("%v_%v", SellerId, roomid)).Val()
	if len(livingdata) == 0 {
		ctx.JSON(http.StatusBadRequest, xenum.LiveNotAvailable)
		return
	}
	if reqdata.IsVisitor == 1 {
		reqdata.Password = xutils.Md5(reqdata.Account)
	} else {
		reqdata.Password = xutils.Md5(reqdata.Password)
	}
	var userdata *model.XUser
	for {
		tb := xapp.DbQuery().XUser
		itb := tb.WithContext(ctx)
		itb = itb.Where(tb.SellerID.Eq(SellerId), tb.Account.Eq(reqdata.Account))
		ud, err := itb.First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if reqdata.IsVisitor != 1 {
				ctx.JSON(http.StatusBadRequest, xenum.UserNotFound)
				return
			}
			errx := xapp.DbQuery().XUser.Create(&model.XUser{
				SellerID:   SellerId,
				Account:    reqdata.Account,
				Password:   reqdata.Password,
				IsVisitor:  1,
				State:      1,
				Agent:      reqdata.SaleId,
				LoginTime:  carbon.Now().StdTime(),
				CreateTime: carbon.Now().StdTime(),
			})
			if errx != nil {
				ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, errx.Error()))
				return
			}
			continue
		}
		userdata = ud
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
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.ID.Eq(userdata.ID))
	_, err := itb.Updates(map[string]interface{}{
		tb.LoginIP.ColumnName().String():         tokendata.Ip,
		tb.LoginTime.ColumnName().String():       carbon.Now().StdTime(),
		tb.Token.ColumnName().String():           userdata.Token,
		tb.LoginCount.ColumnName().String():      gorm.Expr(tb.LoginCount.ColumnName().String() + "+1"),
		tb.LoginIPLocation.ColumnName().String(): GetLocation(tokendata.Ip),
	})
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

type open_hongbao_req struct {
	Id int32 `validate:"required" json:"id"`
}

type open_hongbao_res struct {
	Amount float64 `json:"amount"`
}

// @Router /open_hongbao [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body open_hongbao_req true "请求参数"
// @Success 200  {object} open_hongbao_res "响应数据"
func open_hongbao(ctx *gin.Context) {
	locker.Lock()
	defer locker.Unlock()

	var reqdata open_hongbao_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	response := new(open_hongbao_res)
	{
		tb := xapp.DbQuery().XHongbaoex
		itb := tb.WithContext(ctx)
		itb1 := itb.Where(tb.SellerID.Eq(token.SellerId))
		itb1 = itb1.Where(tb.HongbaoID.Eq(reqdata.Id))
		itb1 = itb1.Where(tb.Account.Eq(token.Account))
		recorddata, err := itb1.First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
			return
		}
		if recorddata != nil {
			response.Amount = -1
			ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
			return
		}
	}
	var hongbao *model.XHongbao
	var err error
	{
		tb := xapp.DbQuery().XHongbao
		itb := tb.WithContext(ctx)
		itb1 := itb.Where(tb.SellerID.Eq(token.SellerId))
		itb1 = itb1.Where(tb.ID.Eq(reqdata.Id))
		hongbao, err = itb1.First()
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			response.Amount = -2
			ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
			return
		}
		now := carbon.Now()
		timedif := carbon.CreateFromStdTime(hongbao.CreateTime).DiffInMinutes(now)
		if timedif > 30000 {
			response.Amount = -3
			ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
			return
		}
	}
	if hongbao.UsedCount >= hongbao.TotalCount || hongbao.UsedAmount >= hongbao.TotalAmount {
		response.Amount = -4
		ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
		return
	}
	avage := hongbao.TotalAmount / float64(hongbao.TotalCount)
	max_amount := avage * 1.3
	min_amount := avage * 0.5
	response.Amount = float64(rand.Float64()*(max_amount-min_amount) + min_amount)
	_, err = xapp.DbQuery().XHongbao.Where(xapp.DbQuery().XHongbao.ID.Eq(reqdata.Id)).Updates(map[string]interface{}{
		xapp.DbQuery().XHongbao.UsedCount.ColumnName().String():  hongbao.UsedCount + 1,
		xapp.DbQuery().XHongbao.UsedAmount.ColumnName().String(): hongbao.UsedAmount + response.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	err = xapp.DbQuery().XHongbaoex.Create(&model.XHongbaoex{
		SellerID:   token.SellerId,
		HongbaoID:  reqdata.Id,
		Account:    token.Account,
		Amount:     response.Amount,
		CreateTime: carbon.Now().StdTime(),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}
