package service_user

import (
	"context"
	"fmt"
	"time"
	"xclientapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/utils"
	"xcom/xcom"

	"github.com/beego/beego/logs"
	"github.com/google/uuid"
)

type ServiceUser struct {
}

func (this *ServiceUser) Init() {
}

type UserRegisterReq struct {
	SellerId int    `json:"-"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserRegisterRes struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (this *ServiceUser) UserRegister(host string, ip string, reqdata *UserRegisterReq) (response *UserRegisterRes, merr map[string]interface{}, err error) {
	locker := enum.Lock_UserRegister + reqdata.Account
	if !server.Redis().Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}
	return response, nil, err
}

type UserLoginReq struct {
	SellerId  int    `json:"-"`
	Account   string `validate:"required" json:"account"`
	Password  string `json:"password"`
	IsVisitor int    `json:"is_visitor"`
}

type UserLoginRes struct {
	Account string `json:"account"`
	UserId  int    `json:"user_id"`
	Token   string `json:"token"`
}

func (this *ServiceUser) UserLogin(host string, ip string, reqdata *UserLoginReq) (response *UserLoginRes, merr map[string]interface{}, err error) {
	locker := enum.Lock_UserLogin + reqdata.Account
	if !server.Redis().Lock(locker, 1) {
		return nil, enum.TooManyRequest, nil
	}
	reqdata.SellerId = xcom.GetSellerId(host)
	if reqdata.SellerId == 0 {
		return nil, enum.SellerNotFound, nil
	}

	if reqdata.IsVisitor == enum.StateYes {
		reqdata.Password = utils.Md5(reqdata.Account)
	} else {
		reqdata.Password = utils.Md5(reqdata.Password)
	}

	rediskey := fmt.Sprintf("account:%v:%v", reqdata.SellerId, reqdata.Account)
	accountdata, err := server.Redis().GetCacheMap(rediskey, func() (*utils.XMap, error) {
		rows, err := server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, reqdata.SellerId).Where(edb.Account+edb.EQ, reqdata.Account).Rows()
		if err != nil {
			return nil, err
		}
		data := utils.DbFirst(rows)
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
		UserId := xcom.NewUserId()
		err := server.Db().Table(edb.TableUser).Create(map[string]interface{}{
			edb.SellerId:  reqdata.SellerId,
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
		rows, err := server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, reqdata.SellerId).Where(edb.Account+edb.EQ, reqdata.Account).Rows()
		if err != nil {
			logs.Error("UserLogin:", err)
			return nil, nil, err
		}
		accountdata = utils.DbFirst(rows)
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
	tokendata.SellerId = reqdata.SellerId
	tokendata.UserId = accountdata.Int(edb.UserId)
	tokendata.IsVisitor = accountdata.Int(edb.IsVisitor)
	server.SetToken(accountdata.String(edb.Token), &tokendata)
	response = &UserLoginRes{}
	response.Account = accountdata.String(edb.Account)
	response.UserId = accountdata.Int(edb.UserId)
	response.Token = accountdata.String(edb.Token)
	return response, nil, err
}

type UserTestReq struct {
	SellerId int `json:"-"`
	UserId   int `json:"user_id"`
}

type UserTestRes struct {
	UserId int `json:"user_id"`
}

func (this *ServiceUser) UserTest(reqdata *UserTestReq) (response *UserTestRes, merr map[string]interface{}, err error) {
	return nil, nil, nil
}
