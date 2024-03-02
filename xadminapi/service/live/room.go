package service_live

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"

	"github.com/gin-gonic/gin"
)

type ServiceLiveRoom struct {
}

func (this *ServiceLiveRoom) Init() {

}

type GetLiveRoomListReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量
}

type GetLiveRoomListRes struct {
	Total int64             `json:"total"` // 总数
	Data  []model.XLiveRoom `json:"data"`  // 数据
}

func (this *ServiceLiveRoom) GetLiveRoomList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetLiveRoomListReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil, nil
	}
	data := GetLiveRoomListRes{}
	db := server.Db().Model(&model.XLiveRoom{})
	db = utils.DbWhere(db, "seller_id", token.SellerId, int(0))
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

func (this *ServiceLiveRoom) get_stream_url(app string, name string) (string, string) {
	return "", ""
}

func (this *ServiceLiveRoom) CreateLiveRoom(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	return nil, nil
}

type UpdateLiveRoomReq struct {
	Id      int    `json:"id"`      // 直播间Id
	Name    string `json:"name"`    // 直播间名称
	Account string `json:"account"` // 直播间账号
	State   int    `json:"state"`   // 直播间状态
}

func (this *ServiceLiveRoom) UpdateLiveRoom(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateLiveRoomReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}
	updatedata := map[string]interface{}{}
	if reqdata.Name != "" {
		updatedata["name"] = reqdata.Name
	}
	if reqdata.Account != "" {
		updatedata["account"] = reqdata.Account
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata["state"] = reqdata.State
	}
	db := server.Db().Model(&model.XLiveRoom{}).Where("id = ? and seller_id = ?", reqdata.Id, token.SellerId).Updates(updatedata)
	return nil, db.Error
}
