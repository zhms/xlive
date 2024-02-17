package service_user

import (
	"xclientapi/server"
	"xcom/enum"

	"github.com/shopspring/decimal"
)

type ServiceUser struct {
}

func (this *ServiceUser) Init() {
}

type UserRegisterReq struct {
	SellerId int    `json:"-"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Agent    int    `json:"agent"`
}

type UserRegisterRes struct {
	UserId int                        `json:"user_id"`
	Token  string                     `json:"token"`
	Amount map[string]decimal.Decimal `json:"amount"`
}

func (this *ServiceUser) UserRegister(host string, ip string, reqdata *UserRegisterReq) (response *UserRegisterRes, merr map[string]interface{}, err error) {
	locker := enum.Lock_UserRegister + reqdata.Account
	if !server.Redis(enum.RedisUserRw).Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}

	// reqdata.SellerId = 1
	// ChannelId := 1

	// //检查账号是否存在
	// reqdata.Account = strings.ToLower(reqdata.Account)
	// {
	// 	count := int64(0)
	// 	db := server.Db().Model(&model.XUser{})
	// 	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	// 	db = db.Where(edb.ChannelId+edb.EQ, ChannelId)
	// 	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	// 	db = db.Count(&count)
	// 	if db.Error != nil {
	// 		return nil, nil, db.Error
	// 	}
	// 	if count > 0 {
	// 		return nil, enum.UserExist, nil
	// 	}
	// }
	// //检查代理是否存在
	// {
	// 	if reqdata.Agent > 0 {
	// 		agentdata := &model.XUser{}
	// 		db := server.Db().Model(agentdata).Select(edb.State)
	// 		db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	// 		db = db.Where(edb.ChannelId+edb.EQ, ChannelId)
	// 		db = db.Where(edb.UserId+edb.EQ, reqdata.Agent)
	// 		db = db.First(agentdata)
	// 		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
	// 			return nil, nil, db.Error
	// 		}
	// 		if db.Error == gorm.ErrRecordNotFound {
	// 			return nil, enum.AgentNotFound, nil
	// 		}
	// 		if agentdata.State != enum.StateYes {
	// 			return nil, enum.AgentStateError, nil
	// 		}
	// 	}
	// }
	// //分配UserId
	// UserId := cuser.NewUserId()
	// if UserId == 0 {
	// 	return nil, nil, errors.New("分配UserId失败")
	// }
	// reqdata.Password = utils.Md5(reqdata.Password)
	// token := uuid.New().String()
	// server.Db(enum.DbUserRW).Transaction(func(tx *gorm.DB) error {
	// 	//创建账号
	// 	createdata := map[string]interface{}{
	// 		edb.SellerId:  reqdata.SellerId,
	// 		edb.ChannelId: ChannelId,
	// 		edb.UserId:    UserId,
	// 		edb.Account:   reqdata.Account,
	// 		edb.Password:  reqdata.Password,
	// 		edb.NickName:  fmt.Sprint(UserId),
	// 		edb.RegIp:     ip,
	// 		edb.LoginIp:   ip,
	// 		edb.LoginTime: time.Now(),
	// 		edb.State:     1,
	// 		edb.Token:     token,
	// 	}
	// 	if reqdata.Agent > 0 {
	// 		createdata[edb.Agent] = reqdata.Agent
	// 	}
	// 	//创建用户数据
	// 	err := tx.Table(edb.TableUser).Create(createdata).Error
	// 	if err != nil {
	// 		logs.Error("UserRegister:", err)
	// 		return err
	// 	}
	// 	//初始化金额
	// 	err = tx.Table(edb.TableUserScore).Create(map[string]interface{}{
	// 		edb.SellerId:  reqdata.SellerId,
	// 		edb.ChannelId: ChannelId,
	// 		edb.UserId:    UserId,
	// 		edb.Symbol:    edb.Usdt,
	// 	}).Error
	// 	if err != nil {
	// 		logs.Error("UserRegister:", err)
	// 		return err
	// 	}
	// 	//生成新token
	// 	tokendata := server.TokenData{}
	// 	tokendata.SellerId = reqdata.SellerId
	// 	tokendata.ChannelId = ChannelId
	// 	tokendata.UserId = UserId
	// 	server.SetToken(token, &tokendata)
	// 	//填充返回数据
	// 	response = &UserRegisterRes{}
	// 	response.UserId = UserId
	// 	response.Token = token
	// 	response.Amount = map[string]decimal.Decimal{
	// 		edb.Usdt:       decimal.Zero,
	// 		edb.UsdtFrozen: decimal.Zero,
	// 		edb.UsdtBank:   decimal.Zero,
	// 	}
	// 	return nil
	// })
	// server.Redis(enum.RedisUserRw).Client().LPush(context.Background(), global.Project+":newuser", UserId)
	return response, nil, err
}

type UserLoginReq struct {
	SellerId int    `json:"-"`
	Account  string `validate:"required" json:"account"`
	Password string `validate:"required" json:"password"`
}

type UserLoginRes struct {
	UserId int                        `json:"user_id"`
	Token  string                     `json:"token"`
	Amount map[string]decimal.Decimal `json:"amount"`
}

func (this *ServiceUser) UserLogin(host string, ip string, reqdata *UserLoginReq) (response *UserLoginRes, merr map[string]interface{}, err error) {
	locker := enum.Lock_UserLogin + reqdata.Account
	if !server.Redis(enum.RedisUserRw).Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}
	reqdata.SellerId = 1
	//ChannelId := 1

	// userdata := &model.XUser{}
	// db := server.Db(enum.DbUserRW).Model(userdata)
	// db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	// db = db.Where(edb.ChannelId+edb.EQ, ChannelId)
	// db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	// db = db.First(userdata)
	// if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
	// 	logs.Error("UserLogin:", db.Error)
	// 	return nil, enum.InternalError, nil
	// }
	// if db.Error == gorm.ErrRecordNotFound {
	// 	return nil, enum.UserNotFound, nil
	// }
	// if userdata.State != 1 {
	// 	return nil, enum.UserStateError, nil
	// }
	// if userdata.Password != utils.Md5(reqdata.Password) {
	// 	return nil, enum.UserPasswordError, nil
	// }
	// //删除旧的token
	// server.DelToken(userdata.Token)
	// //生成新的token
	// token := uuid.New().String()
	// tokendata := server.TokenData{}
	// tokendata.SellerId = reqdata.SellerId
	// tokendata.ChannelId = ChannelId
	// tokendata.UserId = userdata.UserId
	// server.SetToken(token, &tokendata)
	// response = &UserLoginRes{}
	// response.UserId = userdata.UserId
	// response.Token = token
	// cuserdata := cuser.GetCacheData(userdata.UserId)
	// response.Amount = cuserdata.RawData
	// //更新登录信息
	// db = server.Db(enum.DbUserRW).Table(edb.TableUser)
	// db = db.Where(edb.UserId+edb.EQ, userdata.UserId)
	// db.Updates(map[string]interface{}{
	// 	edb.Token:     token,
	// 	edb.LoginIp:   ip,
	// 	edb.LoginTime: time.Now(),
	// })
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
	// response = &UserTestRes{}
	// response.UserId = reqdata.UserId
	// cuser.AlterUserAmount(reqdata.UserId, edb.Usdt, decimal.NewFromInt(100), decimal.Zero, decimal.Zero, enum.AmountChangeType_AdminAlter, "测试")
	// return response, nil, err
	// r, err := rpcclient.ThirdApi.PgLogin(123)
	// fmt.Println(r)
	return nil, nil, nil
}
