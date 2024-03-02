package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"

	"github.com/gin-gonic/gin"
)

type GetAdminLoginLogReq struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	ChannelId int    `json:"channel_id"` // 渠道商
	Account   string `json:"account"`    // 操作人
	LoginIp   string `json:"login_ip"`   // 登录Ip
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type GetAdminLoginLogRes struct {
	Total int64                  `json:"total"` // 总数
	Data  []model.XAdminLoginLog `json:"data"`  // 数据
}

// 获取登录日志
func (this *ServiceAdmin) GetLoginLogList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetAdminLoginLogReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetAdminLoginLogRes{}
	db := server.Db().Model(&model.XAdminLoginLog{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, token.SellerId, int(0))
	db = utils.DbWhere(db, edb.ChannelId+edb.EQ, reqdata.ChannelId, int(0))
	db = utils.DbWhere(db, edb.Account+edb.EQ, reqdata.Account, "")
	db = utils.DbWhere(db, edb.LoginIp+edb.EQ, reqdata.LoginIp, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.GTE, reqdata.StartTime, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.LT, reqdata.EndTime, "")
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
	// for i := range data {
	// 	rpchandler.Ip.GetLocation(&data[i].LoginIp, &data[i].IpLocation)
	// }
	return data, nil, err
}

type GetAdminOptLogReq struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	ChannelId int    `json:"channel_id"` // 渠道商
	Account   string `json:"account"`    // 操作人
	OptName   string `json:"opt_name"`   // 操作名
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type GetAdminOptLogRes struct {
	Total int64                `json:"total"` // 总数
	Data  []model.XAdminOptLog `json:"data"`  // 数据
}

// 获取操作日志
func (this *ServiceAdmin) GetOptLogList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetAdminOptLogReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetAdminOptLogRes{}
	db := server.Db().Model(&model.XAdminOptLog{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, token.SellerId, int(0))
	db = utils.DbWhere(db, edb.ChannelId+edb.EQ, reqdata.ChannelId, int(0))
	db = utils.DbWhere(db, edb.Account+edb.EQ, reqdata.Account, "")
	db = utils.DbWhere(db, edb.OptName+edb.EQ, reqdata.OptName, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.GTE, reqdata.StartTime, "")
	db = utils.DbWhere(db, edb.CreateTime+edb.LT, reqdata.EndTime, "")
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
	// for i := range data {
	// 	rpchandler.Ip.GetLocation(&data[i].ReqIp, &data[i].IpLocation)
	// }
	return data, nil, err
}
