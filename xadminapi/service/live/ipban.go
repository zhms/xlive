package service_live

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/xutils"

	"github.com/gin-gonic/gin"
)

type ServiceLiveIpBan struct {
}

func (this *ServiceLiveIpBan) Init() {

}

type GetBanIpReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	Ip           string `json:"ip"`            // ip
	AdminAccount string `json:"admin_account"` // 管理员账号
}

type GetBanIpRes struct {
	Total int64              `json:"total"` // 总数
	Data  []model.XChatBanIP `json:"data"`  // 数据
}

func (this *ServiceLiveIpBan) GetBanIp(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetBanIpReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	data := GetBanIpRes{}
	db := server.Db().Model(&model.XChatBanIP{})
	db = xutils.DbWhere(db, edb.Ip, reqdata.Ip, "")
	db = xutils.DbWhere(db, edb.AdminAccount, reqdata.AdminAccount, "")
	db = db.Count(&data.Total)
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

type DeleteBanIpReq struct {
	Id int `json:"id"` // id
}

func (this *ServiceLiveIpBan) DeleteBanIp(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(DeleteBanIpReq)
	err = server.Db().Where(edb.Id+edb.EQ, reqdata.Id).Delete(&model.XChatBanIP{}).Error
	if err != nil {
		return nil, nil, err
	}
	return nil, nil, nil
}
