package service_live

import (
	"encoding/json"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/utils"

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
	db = utils.DbWhere(db, edb.SellerId, token.SellerId, int(0))
	db = utils.DbWhere(db, edb.RoomId, reqdata.RoomId, int(0))
	db = db.Count(&data.Total)
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.ASC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

type ChatAuditReq struct {
	Id    int `json:"id"`    // id
	State int `json:"state"` // 状态 2通过,3拒绝,4封ip,5封号
}

func (this *ServiceLiveChat) ChatAudit(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(ChatAuditReq)
	token := server.GetToken(ctx)
	chatdata := model.XChatList{}
	err = server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
		Where(edb.Id+edb.EQ, reqdata.Id).First(&chatdata).Error
	if err != nil {
		return nil, nil, err
	}
	if reqdata.State == 2 || reqdata.State == 3 {
		err = server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
			Where(edb.Id+edb.EQ, reqdata.Id).Update(edb.State, reqdata.State).Error
		if err != nil {
			return nil, nil, err
		}
		if reqdata.State == 2 {
			bytes, _ := json.Marshal(chatdata)
			server.Redis().Client().RPush(ctx, "chat_audit", string(bytes))
		}
	}
	if reqdata.State == 4 || reqdata.State == 5 {
		err = server.Db().Model(&model.XChatList{}).Where(edb.SellerId+edb.EQ, token.SellerId).
			Where(edb.Id+edb.EQ, reqdata.Id).Update(edb.State, 3).Error
		if err != nil {
			return nil, nil, err
		}
		if reqdata.State == 4 {
			server.Db().Table(edb.TableChatBanIp).Create(map[string]interface{}{
				edb.Ip:           chatdata.Ip,
				edb.AdminAccount: token.Account,
			})
			server.Redis().Client().SAdd(ctx, "ip_ban:", chatdata.Ip)
		}
		if reqdata.State == 5 {
			server.Db().Table(edb.TableChatBanIp).Create(map[string]interface{}{
				edb.SellerId: token.SellerId,
			})
			server.Redis().Client().SAdd(ctx, "user_ban:", chatdata.UserId)
		}
	}

	return nil, nil, nil
}
