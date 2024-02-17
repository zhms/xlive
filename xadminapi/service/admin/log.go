package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"
)

type GetAdminLoginLogReq struct {
	Page      int    `form:"page"`       // 页码
	PageSize  int    `form:"page_size"`  // 每页数量
	SellerId  int    `form:"seller_id"`  // 运营商
	ChannelId int    `form:"channel_id"` // 渠道商
	Account   string `form:"account"`    // 操作人
	LoginIp   string `form:"login_ip"`   // 登录Ip
	StartTime string `form:"start_time"` // 开始时间
	EndTime   string `form:"end_time"`   // 结束时间
}

type GetAdminLoginLogRes struct {
	Total int                    `json:"total"` // 总数
	Data  []model.XAdminLoginLog `json:"data"`  // 数据
}

// 获取登录日志
func (this *ServiceAdmin) GetLoginLogList(reqdata *GetAdminLoginLogReq) (total int64, data []model.XAdminLoginLog, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Model(&model.XAdminLoginLog{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, reqdata.SellerId, int(0))
	db = utils.DbWhere(db, edb.ChannelId+edb.EQ, reqdata.ChannelId, int(0))
	db = utils.DbWhere(db, edb.Account+edb.EQ, reqdata.Account, "")
	db = utils.DbWhere(db, edb.LoginIp+edb.EQ, reqdata.LoginIp, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.GTE, reqdata.StartTime, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.LT, reqdata.EndTime, "")
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
	// for i := range data {
	// 	rpchandler.Ip.GetLocation(&data[i].LoginIp, &data[i].IpLocation)
	// }
	return total, data, nil, err
}

type GetAdminOptLogReq struct {
	Page      int    `form:"page"`       // 页码
	PageSize  int    `form:"page_size"`  // 每页数量
	SellerId  int    `form:"seller_id"`  // 运营商
	ChannelId int    `form:"channel_id"` // 渠道商
	Account   string `form:"account"`    // 操作人
	OptName   string `form:"opt_name"`   // 操作名
	StartTime string `form:"start_time"` // 开始时间
	EndTime   string `form:"end_time"`   // 结束时间
}

type GetAdminOptLogRes struct {
	Total int                  `json:"total"` // 总数
	Data  []model.XAdminOptLog `json:"data"`  // 数据
}

// 获取操作日志
func (this *ServiceAdmin) GetOptLogList(reqdata *GetAdminOptLogReq) (total int64, data []model.XAdminOptLog, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Model(&model.XAdminOptLog{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, reqdata.SellerId, int(0))
	db = utils.DbWhere(db, edb.ChannelId+edb.EQ, reqdata.ChannelId, int(0))
	db = utils.DbWhere(db, edb.Account+edb.EQ, reqdata.Account, "")
	db = utils.DbWhere(db, edb.OptName+edb.EQ, reqdata.OptName, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.GTE, reqdata.StartTime, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.LT, reqdata.EndTime, "")
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
	// for i := range data {
	// 	rpchandler.Ip.GetLocation(&data[i].ReqIp, &data[i].IpLocation)
	// }
	return total, data, nil, err
}
