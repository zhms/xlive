package service_admin

import (
	"fmt"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"
)

type GetXChannelReq struct {
	Page        int    `form:"page"`         // 页码
	PageSize    int    `form:"page_size"`    // 每页数量
	SellerId    int    `form:"seller_id"`    // 运营商
	ChannelId   int    `form:"channel_id"`   // 渠道商
	ChannelName string `form:"channel_name"` // 渠道商名称
}

type GetXChannelRes struct {
	Total int              `json:"total"` // 总数
	Data  []model.XChannel `json:"data"`  // 数据
}

// 获取渠道列表
func (this *ServiceAdmin) GetChannelList(reqdata *GetXChannelReq) (total int64, data []model.XChannel, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Model(&model.XChannel{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, reqdata.SellerId, int(0))
	db = utils.DbWhere(db, edb.ChannelId+edb.EQ, reqdata.ChannelId, int(0))
	db = utils.DbWhere(db, edb.ChannelName+edb.LIKE, fmt.Sprintf("%%%s%%", reqdata.ChannelName), "")
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

type CreateXChannelReq struct {
	SellerId    int    `validate:"required" json:"seller_id"`    // 运营商
	ChannelId   int    `validate:"required" json:"channel_id"`   // 渠道商
	ChannelName string `validate:"required" json:"channel_name"` // 渠道商名称
	State       int    `validate:"required" json:"state"`        // 状态 1开启,2关闭
	Memo        string `json:"memo"`                             // 备注
}

// 创建渠道
func (this *ServiceAdmin) CreateChannel(reqdata *CreateXChannelReq) (merr map[string]interface{}, err error) {
	db := server.Db().Model(&model.XChannel{})
	db.Create(map[string]interface{}{
		edb.SellerId:    reqdata.SellerId,
		edb.ChannelId:   reqdata.ChannelId,
		edb.ChannelName: reqdata.ChannelName,
		edb.State:       reqdata.State,
		edb.Memo:        reqdata.Memo,
	})
	return nil, db.Error
}

type UpdateXChannelReq struct {
	SellerId    int    `validate:"required" json:"seller_id"`    // 运营商
	ChannelId   int    `validate:"required" json:"channel_id"`   // 渠道商
	ChannelName string `validate:"required" json:"channel_name"` // 渠道商名称
	State       int    `validate:"required" json:"state"`        // 状态 1开启,2关闭
	Memo        string `json:"memo"`                             // 备注
}

// 更新渠道
func (this *ServiceAdmin) UpdateChannel(reqdata *UpdateXChannelReq) (merr map[string]interface{}, err error) {
	updatedata := map[string]interface{}{}
	updatedata[edb.Memo] = reqdata.Memo
	utils.MapSet(&updatedata, edb.ChannelName, reqdata.ChannelName, "")
	utils.MapSetIn(&updatedata, edb.State, reqdata.State, []interface{}{int(1), int(2)})
	db := server.Db().Model(&model.XChannel{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Where(edb.ChannelId+edb.EQ, reqdata.ChannelId)
	db = db.Updates(updatedata)
	return nil, db.Error
}

type DeleteXChannelReq struct {
	SellerId  int `validate:"required" json:"seller_id"`  // 运营商
	ChannelId int `validate:"required" json:"channel_id"` // 渠道商
}

// 删除渠道
func (this *ServiceAdmin) DeleteChannel(reqdata *DeleteXChannelReq) (rows int64, merr map[string]interface{}, err error) {
	db := server.Db().Model(&model.XChannel{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Where(edb.ChannelId+edb.EQ, reqdata.ChannelId)
	db = db.Delete(&model.XChannel{})
	return db.RowsAffected, nil, db.Error
}
