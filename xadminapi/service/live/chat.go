package service_live

import (
	"encoding/json"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/xutils"

	"github.com/gin-gonic/gin"
)

type ServiceLiveChat struct {
}

func (this *ServiceLiveChat) Init() {

}

type GetChatListReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	RoomId int `json:"room_id"` // 房间Id
}

type GetChatListRes struct {
	Total int64             `json:"total"` // 总数
	Data  []model.XChatList `json:"data"`  // 数据
}

func (this *ServiceLiveChat) GetChatList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetChatListReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetChatListRes{}
	db := server.Db().Model(&model.XChatList{})
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

type ChatAuditReq struct {
	Id    int `json:"id"`    // id
	State int `json:"state"` // 状态 2通过,3拒绝,4封ip
}

func (this *ServiceLiveChat) ChatAudit(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(ChatAuditReq)
	token := server.GetToken(ctx)
	chatdata := model.XChatList{}
	db := server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
		Where(edb.Id+edb.EQ, reqdata.Id).First(&chatdata)
	if db.Error != nil {
		return nil, nil, db.Error
	}
	if chatdata.State != 1 {
		return nil, enum.AlreadyAudited, nil
	}
	if reqdata.State == 2 || reqdata.State == 3 {
		db = server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
			Where(edb.Id+edb.EQ, reqdata.Id).Where(edb.State+edb.EQ, 1).Update(edb.State, reqdata.State)
		if db.Error != nil {
			return nil, nil, db.Error
		}
		if reqdata.State == 2 && db.RowsAffected > 0 {
			bytes, _ := json.Marshal(chatdata)
			server.Redis().Client().RPush(ctx, "chat_audit", string(bytes))
		}
	}
	if reqdata.State == 4 || reqdata.State == 5 {
		db = server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
			Where(edb.Id+edb.EQ, reqdata.Id).Where(edb.State+edb.EQ, 1).Update(edb.State, 3)
		if db.Error != nil {
			return nil, nil, db.Error
		}
		if reqdata.State == 4 && db.RowsAffected > 0 {
			server.Db().Table(edb.TableChatBanIp).Create(map[string]interface{}{
				edb.Ip:           chatdata.Ip,
				edb.AdminAccount: token.Account,
			})
			server.Redis().Client().SAdd(ctx, "ip_ban", chatdata.Ip)
		}
	}
	return nil, nil, nil
}
