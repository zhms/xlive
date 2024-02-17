package service_active

import (
	"fmt"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/utils"
)

type ServiceActive struct {
}

func (this *ServiceActive) Init() {

}

type GetXActiveReq struct {
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页数量
	SellerId int    `form:"seller_id"` // 运营商
	Title    string `form:"title"`     // 标题模糊查询
}

type GetXActiveRes struct {
	Total int             `json:"total"` // 总数
	Data  []model.XActive `json:"data"`  // 数据
}

// 获取活动列表
func (this *ServiceActive) GetActiveList(reqdata *GetXActiveReq) (total int64, data []model.XActive, merr map[string]interface{}, err error) {
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	db := server.Db().Model(&model.XActive{})
	db = utils.DbWhere(db, edb.SellerId+edb.EQ, reqdata.SellerId, int(0))
	db = utils.DbWhere(db, edb.Title+edb.LIKE, fmt.Sprintf("%%%s%%", reqdata.Title), "")
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

type CreateXActiveReq struct {
	SellerId  int    `validate:"required" json:"seller_id"` // 运营商
	ChannelId string `json:"channel_id"`                    // 适用渠道 json格式
	ActiveId  int    `json:"active_id"`                     // 活动id
	Picture   string `json:"picture"`                       // 图片路径
	Title     string `json:"title"`                         // 标题
	Content   string `json:"content"`                       // 文字内容
	Memo      string `json:"memo"`                          // 备注
	State     int    `json:"state"`                         // 状态
	Sort      int    `json:"sort"`                          // 排序
}

// 创建活动
func (this *ServiceActive) CreateActive(reqdata *CreateXActiveReq) (merr map[string]interface{}, err error) {
	if reqdata.State != enum.StateYes {
		reqdata.State = enum.StateNo
	}
	db := server.Db().Model(&model.XActive{})
	db.Create(map[string]interface{}{
		edb.SellerId:  reqdata.SellerId,
		edb.ChannelId: reqdata.ChannelId,
		edb.ActiveId:  reqdata.ActiveId,
		edb.Picture:   reqdata.Picture,
		edb.Title:     reqdata.Title,
		edb.Content:   reqdata.Content,
		edb.Memo:      reqdata.Memo,
		edb.State:     reqdata.State,
		edb.Sort:      reqdata.Sort,
	})
	return nil, db.Error
}

type UpdateXActiveReq struct {
	SellerId  int    `validate:"required" json:"seller_id"` // 运营商
	Id        int    `validate:"required" json:"id"`
	ChannelId string `json:"channel_id"` // 适用渠道 json格式
	ActiveId  int    `json:"active_id"`  // 活动id
	Picture   string `json:"picture"`    // 图片路径
	Title     string `json:"title"`      // 标题
	Content   string `json:"content"`    // 文字内容
	Memo      string `json:"memo"`       // 备注
	State     int    `json:"state"`      // 状态
	Sort      int    `json:"sort"`       // 排序
}

// 更新活动
func (this *ServiceActive) UpdateActive(reqdata *UpdateXActiveReq) (merr map[string]interface{}, err error) {
	updatedata := map[string]interface{}{}
	utils.MapSet(&updatedata, edb.ChannelId, reqdata.ChannelId, "")
	utils.MapSet(&updatedata, edb.ActiveId, reqdata.ActiveId, 0)
	utils.MapSet(&updatedata, edb.Picture, reqdata.Picture, "")
	utils.MapSet(&updatedata, edb.Title, reqdata.Title, "")
	utils.MapSet(&updatedata, edb.Content, reqdata.Content, "")
	utils.MapSet(&updatedata, edb.Memo, reqdata.Memo, "")
	if reqdata.State == enum.StateYes || reqdata.State == enum.StateNo {
		updatedata[edb.State] = reqdata.State
	}
	utils.MapSet(&updatedata, edb.Sort, reqdata.Sort, 0)
	db := server.Db().Model(&model.XActive{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Where(edb.Id+edb.EQ, reqdata.Id)
	db = db.Updates(updatedata)
	return nil, db.Error
}

type DeleteXActiveReq struct {
	SellerId int `validate:"required" json:"seller_id"` // 运营商
	Id       int `validate:"required" json:"id"`        // id
}

// 删除活动
func (this *ServiceActive) DeleteActive(reqdata *DeleteXActiveReq) (rows int64, merr map[string]interface{}, err error) {
	db := server.Db().Model(&model.XActive{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Where(edb.Id+edb.EQ, reqdata.Id)
	db = db.Delete(&model.XActive{})
	return db.RowsAffected, nil, db.Error
}
