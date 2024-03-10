package service_hongbao

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/xutils"

	"github.com/gin-gonic/gin"
)

type ServiceHongbao struct {
}

func (this *ServiceHongbao) Init() {

}

type GetHongbaoListReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	RoomId int `json:"room_id"` // 房间Id
}

type GetHongbaoListRes struct {
	Total int64            `json:"total"` // 总数
	Data  []model.XHongbao `json:"data"`  // 数据
}

func (this *ServiceHongbao) GetHongbaoList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetHongbaoListReq)
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetHongbaoListRes{}
	db := server.Db().Model(&model.XHongbao{})
	db = xutils.DbWhere(db, edb.SellerId, token.SellerId, int(0))
	db = xutils.DbWhere(db, edb.RoomId, reqdata.RoomId, int(0))
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

type SendHongbaoReq struct {
	TotalAmount int    `validate:"required" json:"total_amount"` // 总金额
	TotalCount  int    `validate:"required" json:"total_count"`  // 总数量
	RoomId      int    `validate:"required" json:"room_id"`      // 房间Id
	Memo        string `json:"memo"`                             // 备注
}

func (this *ServiceHongbao) SendHongbao(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(SendHongbaoReq)
	token := server.GetToken(ctx)
	err = server.Db().Table(edb.TableHongbao).Create(map[string]interface{}{
		edb.SellerId:    token.SellerId,
		edb.RoomId:      reqdata.RoomId,
		edb.TotalCount:  reqdata.TotalCount,
		edb.TotalAmount: reqdata.TotalAmount,
		edb.Memo:        reqdata.Memo,
		edb.Sender:      token.Account,
		edb.CreateTime:  xutils.Now(),
	}).Error
	return nil, err
}
