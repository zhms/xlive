package service_user

import (
	"context"
	"fmt"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/excel"
	"xcom/utils"
	"xcom/xcom"

	"github.com/gin-gonic/gin"
)

type ServiceUser struct {
}

func (this *ServiceUser) Init() {

}

type GetUserReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	UserId  int    `json:"user_id"`  // 用户Id
	Account string `json:"account"`  // 账号
	Agent   string `json:"agent"`    // 代理
	LoginIp string `json:"login_ip"` // 登录Ip
	Export  int    `json:"export"`   // 导出
}

type GetUserRes struct {
	Total int64         `json:"total"` // 总数
	Data  []model.XUser `json:"data"`  // 数据
}

// 获取会员列表
func (this *ServiceUser) GetUserList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetUserReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil, nil
	}
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	if reqdata.Export == 1 {
		reqdata.PageSize = 100000
	}
	data := GetUserRes{}
	db := server.Db().Omit(edb.Password).Model(&model.XUser{})
	db = utils.DbWhere(db, edb.SellerId, token.SellerId, int(0))
	db = utils.DbWhere(db, edb.UserId, reqdata.UserId, int(0))
	db = utils.DbWhere(db, edb.Account, reqdata.Account, "")
	db = utils.DbWhere(db, edb.Agent, reqdata.Agent, "")
	db = utils.DbWhere(db, edb.LoginIp, reqdata.LoginIp, "")
	err = db.Count(&data.Total).Error
	if err != nil {
		return err, nil, nil
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return err, nil, nil
	}
	if reqdata.Export == 1 {
		e := excel.NewExcelBuilder("会员列表_" + fmt.Sprint(utils.Now()))
		defer e.Close()
		e.SetTitle(edb.Id, "Id")
		e.SetTitle(edb.UserId, "会员Id")
		e.SetTitle(edb.Account, "账号")
		e.SetTitle(edb.State, "状态")
		e.SetTitle(edb.ChatState, "聊天状态")
		e.SetTitle(edb.Agent, "业务员")
		e.SetTitle(edb.LoginIp, "登录Ip")
		e.SetTitle(edb.LoginLocation, "登录地点")
		e.SetTitle(edb.LoginCount, "登录时间")
		e.SetTitle(edb.LoginTime, "登录时间")
		e.SetTitle(edb.CreateTime, "注册时间")
		for i, v := range data.Data {
			e.SetValue(edb.Id, v.Id, int64(i+2))
			e.SetValue(edb.UserId, v.UserId, int64(i+2))
			e.SetValue(edb.Account, v.Account, int64(i+2))
			if v.State == 1 {
				e.SetValue(edb.State, "正常", int64(i+2))
			} else {
				e.SetValue(edb.State, "禁用", int64(i+2))
			}
			if v.ChatState == 1 {
				e.SetValue(edb.ChatState, "正常", int64(i+2))
			} else {
				e.SetValue(edb.ChatState, "禁言", int64(i+2))
			}
			e.SetValue(edb.Agent, v.Agent, int64(i+2))
			e.SetValue(edb.LoginIp, v.LoginIP, int64(i+2))
			e.SetValue(edb.LoginLocation, v.LoginLocation, int64(i+2))
			e.SetValue(edb.LoginCount, v.LoginCount, int64(i+2))
			e.SetValue(edb.LoginTime, v.LoginTime, int64(i+2))
			e.SetValue(edb.CreateTime, v.CreateTime, int64(i+2))
		}
		e.Write(ctx)
		return nil, nil, nil
	}
	return data, nil, err
}

type AddUserReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}

// 添加会员
func (this *ServiceUser) AddUser(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(AddUserReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}

	UserId := xcom.NewUserId()
	if UserId <= 0 {
		return enum.NewIdError, nil
	}
	Password := utils.Md5(reqdata.Password)
	err = server.Db().Table(edb.TableUser).Create(map[string]interface{}{
		edb.SellerId:   token.SellerId,
		edb.UserId:     UserId,
		edb.Account:    reqdata.Account,
		edb.Password:   Password,
		edb.State:      1,
		edb.Agent:      token.Account,
		edb.IsVisitor:  enum.StateNo,
		edb.CreateTime: utils.Now(),
	}).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateUserReq struct {
	UserId    int    `json:"user_id" validate:"required"` // 用户Id
	Password  string `json:"password" `                   // 密码
	State     int    `json:"state" `                      // 状态
	ChatState int    `json:"chat_state" `                 // 聊天状态
}

// 更新会员
func (this *ServiceUser) UpdateUser(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateUserReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}
	reqdata.Password = utils.Md5(reqdata.Password)
	updatedata := map[string]interface{}{}
	if reqdata.Password != "" {
		updatedata[edb.Password] = reqdata.Password
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[edb.State] = reqdata.State
	}
	if reqdata.ChatState == 1 || reqdata.ChatState == 2 {
		updatedata[edb.ChatState] = reqdata.ChatState
	}
	userdata := model.XUser{}
	db := server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, token.SellerId).Where(edb.UserId+edb.EQ, reqdata.UserId).First(&userdata)
	if db.Error != nil {
		return nil, db.Error
	}
	db = server.Db().Table(edb.TableUser).Where(edb.SellerId+edb.EQ, token.SellerId).Where(edb.UserId+edb.EQ, reqdata.UserId).Updates(updatedata)
	if db.Error != nil {
		return nil, db.Error
	}
	rediskey := fmt.Sprintf("account:%v:%v", token.SellerId, userdata.Account)
	_, err = server.Redis().Client().Del(context.Background(), rediskey).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
