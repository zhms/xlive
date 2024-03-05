package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/xutils"

	"github.com/gin-gonic/gin"
)

type GetXSellerReq struct {
	Page       int    `json:"page"`        // 页码
	PageSize   int    `json:"page_size"`   // 每页数量
	SellerName string `json:"seller_name"` // 运营商名称
}

type GetXSellerRes struct {
	Total int64           `json:"total"` // 总数
	Data  []model.XSeller `json:"data"`  // 数据
}

// 获取运营商
func (this *ServiceAdmin) GetSellerList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetXSellerReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	db := server.Db().Model(&model.XSeller{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	if reqdata.SellerName != "" {
		db = db.Where(edb.SellerName+edb.EQ, reqdata.SellerName)
	}
	data := GetXSellerRes{}
	err = db.Count(&data.Total).Error
	if err != nil {
		return nil, nil, err
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return err, nil, nil
	}
	return data, nil, err
}

type CreateXSellerReq struct {
	SellerName string `validate:"required" json:"seller_name"` // 运营商名称
	State      int    `validate:"required" json:"state"`       // 状态 1开启,2关闭
	Memo       string `json:"memo"`                            // 备注
}

// 创建运营商
func (this *ServiceAdmin) CreateSeller(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(CreateXSellerReq)
	db := server.Db().Model(&model.XSeller{})
	token := server.GetToken(ctx)
	db = db.Create(map[string]interface{}{
		edb.SellerId:   token.SellerId,
		edb.SellerName: reqdata.SellerName,
		edb.State:      reqdata.State,
		edb.Memo:       reqdata.Memo,
	})
	return nil, db.Error
}

type UpdateXSellerReq struct {
	SellerId   int    `validate:"required" json:"seller_id"` // 运营商
	SellerName string `json:"seller_name"`                   // 运营商名称
	State      int    `json:"state"`                         // 状态 1开启,2关闭
	Memo       string `json:"memo"`                          // 备注
}

// 更新运营商
func (this *ServiceAdmin) UpdateSeller(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateXSellerReq)
	updatedata := map[string]interface{}{}
	updatedata[edb.Memo] = reqdata.Memo
	xutils.MapSet(&updatedata, edb.SellerName, reqdata.SellerName, "")
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[edb.State] = reqdata.State
	}
	db := server.Db().Model(&model.XSeller{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Updates(updatedata)
	return nil, db.Error
}

type DeleteXSellerReq struct {
	SellerId int `validate:"required" json:"seller_id"` // 运营商
}

// 删除运营商
func (this *ServiceAdmin) DeleteSeller(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(DeleteXSellerReq)
	db := server.Db().Model(&model.XSeller{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Delete(&model.XSeller{})
	return nil, db.Error
}
