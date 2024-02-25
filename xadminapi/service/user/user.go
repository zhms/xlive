package service_user

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"
)

type ServiceUser struct {
}

func (this *ServiceUser) Init() {

}

type GetUserReq struct {
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页数量
	SellerId int    `form:"seller_id"` // 运营商
	UserId   int    `form:"user_id"`   // 用户ID
	Account  string `form:"account"`   // 账号
	Agent    string `form:"agent"`     // 代理商
	LoginIp  string `form:"login_ip"`  // 登录IP
}

type GetUserRes struct {
	Total int           `json:"total"` // 总数
	Data  []model.XUser `json:"data"`  // 数据
}

// 获取角色列表
func (this *ServiceUser) GetUserList(reqdata *GetUserReq) (total int64, data []model.XUser, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Omit(edb.Password).Model(&model.XUser{})
	db = utils.DbWhere(db, edb.SellerId, reqdata.SellerId, int(0))
	db = utils.DbWhere(db, edb.UserId, reqdata.UserId, int(0))
	db = utils.DbWhere(db, edb.Account, reqdata.Account, "")
	db = utils.DbWhere(db, edb.Agent, reqdata.Agent, "")
	db = utils.DbWhere(db, edb.LoginIp, reqdata.LoginIp, "")
	err = db.Count(&total).Error
	if err != nil {
		return 0, nil, nil, err
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data).Error
	if err != nil {
		return 0, nil, nil, err
	}
	return total, data, nil, err
}
