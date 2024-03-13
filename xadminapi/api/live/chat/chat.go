package live_chat

import (
	"encoding/json"
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_chat_data", get_chat_data)
	xglobal.ApiV1.POST("/update_chat_data", update_chat_data)
}

type get_chat_data_req struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	RoomId int `json:"room_id"` // 房间Id
}

type get_chat_data_res struct {
	Total int64           `json:"total"` // 总数
	Data  []xdb.XChatData `json:"data"`  // 数据
}

// @Router /get_chat_data [post]
// @Tags 直播间 - 互动列表
// @Summary 获取互动列表
// @Param x-token header string true "token"
// @Param body body get_chat_data_req true "请求参数"
// @Success 200  {object} get_chat_data_res "响应数据"
func get_chat_data(ctx *gin.Context) {
	var reqdata get_chat_data_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := admin.GetToken(ctx)
	if token == nil {
		return
	}
	response := new(get_chat_data_res)
	db := xapp.Db().Model(&xdb.XChatData{})
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.RoomId > 0 {
		db = db.Where(xdb.RoomId+xdb.EQ, reqdata.RoomId)
	}
	db = db.Count(&response.Total)
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(xdb.Id + xdb.DESC)
	err := db.Find(&response.Data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type update_chat_data_req struct {
	Id    int `json:"id"`    // id
	State int `json:"state"` // 状态 2通过,3拒绝,4封ip
}

// @Router /update_chat_data [post]
// @Tags 直播间 - 互动列表
// @Summary 审核互动列表
// @Param x-token header string true "token"
// @Param body body update_chat_data_req true "请求参数"
// @Success 200 "响应数据"
func update_chat_data(ctx *gin.Context) {
	var reqdata update_chat_data_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := admin.GetToken(ctx)
	chatdata := xdb.XChatData{}
	db := xapp.Db().Model(&xdb.XChatData{}).Where(xdb.SellerId+xdb.EQ, token.SellerId).Where(xdb.Id+xdb.EQ, reqdata.Id).First(&chatdata)
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	if chatdata.State != 1 {
		ctx.JSON(http.StatusBadRequest, xenum.AlreadyAudited)
		return
	}
	if reqdata.State == 2 || reqdata.State == 3 {
		db = xapp.Db().Model(&xdb.XChatData{}).Where(xdb.SellerId+xdb.EQ, token.SellerId).
			Where(xdb.Id+xdb.EQ, reqdata.Id).Where(xdb.State+xdb.EQ, 1).Update(xdb.State, reqdata.State)
		if db.Error != nil {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
			return
		}
		if reqdata.State == 2 && db.RowsAffected > 0 {
			bytes, _ := json.Marshal(chatdata)
			_, err := xapp.Redis().Client().RPush(ctx, "chat_audit", string(bytes)).Result()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
				return
			}
		}
	}
	if reqdata.State == 4 || reqdata.State == 5 {
		db = xapp.Db().Model(&xdb.XChatData{}).Where(xdb.SellerId+xdb.EQ, token.SellerId).Where(xdb.Id+xdb.EQ, reqdata.Id).Where(xdb.State+xdb.EQ, 1).Update(xdb.State, 3)
		if db.Error != nil {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
			return
		}
		if reqdata.State == 4 && db.RowsAffected > 0 {
			xapp.Db().Model(&xdb.XChatBanIp{}).Create(map[string]interface{}{
				xdb.SellerId:     token.SellerId,
				xdb.Ip:           chatdata.Ip,
				xdb.AdminAccount: token.Account,
			})
			_, err := xapp.Redis().Client().SAdd(ctx, "ip_ban", chatdata.Ip).Result()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
				return
			}
		}
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
