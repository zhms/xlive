package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"
)

type GetXSellerReq struct {
	Page       int    `form:"page"`        // 页码
	PageSize   int    `form:"page_size"`   // 每页数量
	SellerId   int    `form:"seller_id"`   // 运营商
	SellerName string `form:"seller_name"` // 运营商名称
}

type GetXSellerRes struct {
	Total int             `json:"total"` // 总数
	Data  []model.XSeller `json:"data"`  // 数据
}

// 获取运营商
func (this *ServiceAdmin) GetSellerList(reqdata *GetXSellerReq) (total int64, data []model.XSeller, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Model(&model.XSeller{})
	if reqdata.SellerId > 0 {
		db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	}
	if reqdata.SellerName != "" {
		db = db.Where(edb.SellerName+edb.EQ, reqdata.SellerName)
	}
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

type CreateXSellerReq struct {
	SellerId   int    `validate:"required" json:"seller_id"`   // 运营商
	SellerName string `validate:"required" json:"seller_name"` // 运营商名称
	State      int    `validate:"required" json:"state"`       // 状态 1开启,2关闭
	Memo       string `json:"memo"`                            // 备注
}

// 创建运营商
func (this *ServiceAdmin) CreateSeller(reqdata *CreateXSellerReq) (merr map[string]interface{}, err error) {
	db := server.Db().Model(&model.XSeller{})
	db = db.Create(map[string]interface{}{
		edb.SellerId:   reqdata.SellerId,
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
func (this *ServiceAdmin) UpdateSeller(reqdata *UpdateXSellerReq) (merr map[string]interface{}, err error) {
	updatedata := map[string]interface{}{}
	updatedata[edb.Memo] = reqdata.Memo
	utils.MapSet(&updatedata, edb.SellerName, reqdata.SellerName, "")
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
func (this *ServiceAdmin) DeleteSeller(reqdata *DeleteXSellerReq) (rows int64, merr map[string]interface{}, err error) {
	db := server.Db().Model(&model.XSeller{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Delete(&model.XSeller{})
	return db.RowsAffected, nil, db.Error
}
