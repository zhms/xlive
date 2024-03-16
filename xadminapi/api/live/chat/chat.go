package live_chat

import (
	"encoding/json"
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_chat_data", admin.Auth("直播间", "互动列表", "查", ""), get_chat_data)
	xglobal.ApiV1.POST("/update_chat_data", admin.Auth("直播间", "互动列表", "改", "聊天审核"), update_chat_data)
}

type get_chat_data_req struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	RoomId int `json:"room_id"` // 房间Id
}

type get_chat_data_res struct {
	Total int64          `json:"total"` // 总数
	Data  []*model.XChat `json:"data"`  // 数据
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

	response := new(get_chat_data_res)
	tb := xapp.DbQuery().XChat
	itb := tb.WithContext(ctx)

	itb = itb.Where(tb.SellerID.Eq(int32(token.SellerId)))
	if reqdata.RoomId > 0 {
		itb = itb.Where(tb.RoomID.Eq(int32(reqdata.RoomId)))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
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
	tb := xapp.DbQuery().XChat
	itb := tb.WithContext(ctx).Where(tb.SellerID.Eq(int32(token.SellerId)), tb.ID.Eq(int32(reqdata.Id)))
	chatdata, err := itb.First()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	if chatdata.State != 1 {
		ctx.JSON(http.StatusBadRequest, xenum.AlreadyAudited)
		return
	}
	if reqdata.State == 2 || reqdata.State == 3 {
		tb = xapp.DbQuery().XChat
		itb = tb.WithContext(ctx)
		itb = itb.Where(tb.SellerID.Eq(int32(token.SellerId)))
		itb = itb.Where(tb.ID.Eq(int32(reqdata.Id)))
		itb = itb.Where(tb.State.Eq(1))
		update_result, err := itb.Update(tb.State, int32(reqdata.State))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
			return
		}
		if reqdata.State == 2 && update_result.RowsAffected > 0 {
			bytes, _ := json.Marshal(chatdata)
			_, err := xapp.Redis().Client().RPush(ctx, "chat_audit", string(bytes)).Result()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
				return
			}
		}
	}
	if reqdata.State == 4 || reqdata.State == 5 {
		tb = xapp.DbQuery().XChat
		itb = tb.WithContext(ctx)
		itb = itb.Where(tb.SellerID.Eq(int32(token.SellerId)), tb.ID.Eq(int32(reqdata.Id)), tb.State.Eq(1))
		update_result, err := itb.Update(tb.State, 3)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
			return
		}
		if reqdata.State == 4 && update_result.RowsAffected > 0 {
			tb := xapp.DbQuery().XChatBanIP
			itb := tb.WithContext(ctx)
			itb.Create(&model.XChatBanIP{
				SellerID:     int32(token.SellerId),
				IP:           chatdata.IP,
				AdminAccount: token.Account,
			})
			_, err := xapp.Redis().Client().SAdd(ctx, "ip_ban", chatdata.IP).Result()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
				return
			}
		}
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
