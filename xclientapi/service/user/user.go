package service_user

import (
	"context"
	"fmt"
	"strings"
	"time"
	"xclientapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/xcom"
	"xcom/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceUser struct {
}

func (this *ServiceUser) Init() {
}

// 用户注册
type UserRegisterReq struct {
	SellerId int    `json:"-"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserRegisterRes struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (this *ServiceUser) UserRegister(reqdata *UserRegisterReq) (response *UserRegisterRes, merr map[string]interface{}, err error) {
	locker := enum.Lock_UserRegister + reqdata.Account
	if !server.Redis().Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}
	return response, nil, err
}

// 用户登录
type UserLoginReq struct {
	Account   string `validate:"required" json:"account"` // 账号
	Password  string `json:"password"`                    // 密码
	IsVisitor int    `json:"is_visitor"`                  // 是否游客
	SalesId   int    `json:"sales_id"`                    //业务员
}

type UserLoginRes struct {
	Account   string `json:"account"`    // 账号
	UserId    int    `json:"user_id"`    // 用户Id
	Token     string `json:"token"`      // token
	IsVisitor int    `json:"is_visitor"` // 是否游客
	LiveData  string `json:"live_data"`  // 直播数据
}

func (this *ServiceUser) UserLogin(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(*UserLoginReq)
	locker := enum.Lock_UserLogin + reqdata.Account
	if !server.Redis().Lock(locker, 1) {
		return nil, enum.TooManyRequest, nil
	}
	host := ctx.Request.Host
	host = strings.Replace(host, "www.", "", -1)
	host = strings.Split(host, ":")[0]

	SellerId := xcom.GetSellerId(host)
	if SellerId == 0 {
		SellerId = 1
	}
	RoomId := ctx.Request.Header.Get("RoomId")
	livingdata := server.Redis().Client().HGet(context.Background(), "living", fmt.Sprintf("%v_%v", SellerId, RoomId)).Val()
	if len(livingdata) == 0 {
		return nil, enum.LiveNotAvailable, nil
	}

	if reqdata.IsVisitor == enum.StateYes {
		reqdata.Password = xutils.Md5(reqdata.Account)
	} else {
		reqdata.Password = xutils.Md5(reqdata.Password)
	}
	rediskey := fmt.Sprintf("account:%v:%v", SellerId, reqdata.Account)
	accountdata, err := server.Redis().GetCacheMap(rediskey, func() (*xutils.XMap, error) {
		rows, err := server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, SellerId).Where(edb.Account+edb.EQ, reqdata.Account).Rows()
		if err != nil {
			return nil, err
		}
		data := xutils.DbFirst(rows)
		if data == nil {
			return nil, nil
		}
		_, err = server.Redis().Client().Set(context.Background(), rediskey, data.ToString(), time.Second*60*60*24).Result()
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		logs.Error("UserLogin:", err)
		return nil, nil, err
	}
	if accountdata == nil {
		if reqdata.IsVisitor == enum.StateNo {
			return nil, enum.UserNotFound, nil
		}
		UserId := xcom.NewUserId()
		err := server.Db().Table(edb.TableUser).Create(map[string]interface{}{
			edb.SellerId:  SellerId,
			edb.UserId:    UserId,
			edb.Account:   reqdata.Account,
			edb.Password:  reqdata.Password,
			edb.State:     enum.StateYes,
			edb.IsVisitor: reqdata.IsVisitor,
		}).Error
		if err != nil {
			logs.Error("UserLogin:", err)
			return nil, nil, err
		}
		rows, err := server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, SellerId).Where(edb.Account+edb.EQ, reqdata.Account).Rows()
		if err != nil {
			logs.Error("UserLogin:", err)
			return nil, nil, err
		}
		accountdata = xutils.DbFirst(rows)
		_, err = server.Redis().Client().Set(context.Background(), rediskey, accountdata.ToString(), time.Second*60*60*24).Result()
		if err != nil {
			logs.Error("UserLogin:", err)
			return nil, nil, err
		}
	}

	if accountdata.String(edb.Password) != reqdata.Password {
		return nil, enum.UserPasswordError, nil
	}
	server.DelToken(accountdata.String(edb.Token))
	accountdata.Set(edb.Token, uuid.New().String())
	server.Redis().Client().Set(context.Background(), rediskey, accountdata.ToString(), time.Second*60*60*24).Result()

	tokendata := server.TokenData{}
	tokendata.SellerId = SellerId
	tokendata.Account = accountdata.String(edb.Account)
	tokendata.UserId = accountdata.Int(edb.UserId)
	tokendata.IsVisitor = accountdata.Int(edb.IsVisitor)
	tokendata.Token = accountdata.String(edb.Token)
	tokendata.Ip = ctx.ClientIP()
	server.SetToken(accountdata.String(edb.Token), &tokendata)
	response := &UserLoginRes{}
	response.Account = accountdata.String(edb.Account)
	response.UserId = accountdata.Int(edb.UserId)
	response.Token = accountdata.String(edb.Token)
	response.IsVisitor = accountdata.Int(edb.IsVisitor)
	response.LiveData = livingdata
	return response, nil, err
}
