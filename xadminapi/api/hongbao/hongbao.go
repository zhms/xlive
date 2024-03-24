package hongbao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_hongbao", admin.Auth("红包管理", "红包管理", "查", ""), get_hongbao)
	xglobal.ApiV1.POST("/get_hongbao_detail", admin.Auth("红包管理", "红包管理", "查", ""), get_hongbao_detail)
	xglobal.ApiV1.POST("/create_hongbao", admin.Auth("红包管理", "红包管理", "发红包", "发红包"), create_hongbao)
}

type get_hongbao_req struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	RoomId int32 `json:"room_id"` // 房间Id
}

type get_hongbao_res struct {
	Total int64             `json:"total"` // 总数
	Data  []*model.XHongbao `json:"data"`  // 数据
}

// @Router /get_hongbao [post]
// @Tags 红包管理
// @Summary 获取红包
// @Param x-token header string true "token"
// @Param body body get_hongbao_req true "请求参数"
// @Success 200  {object} get_hongbao_res "响应数据"
func get_hongbao(ctx *gin.Context) {
	var reqdata get_hongbao_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 15
	}
	response := new(get_hongbao_res)
	tb := xapp.DbQuery().XHongbao
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	if reqdata.RoomId > 0 {
		itb = itb.Where(tb.RoomID.Eq(reqdata.RoomId))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_hongbao_req struct {
	TotalAmount float64 `validate:"required" json:"total_amount"` // 总金额
	TotalCount  int32   `validate:"required" json:"total_count"`  // 总数量
	RoomId      int32   `validate:"required" json:"room_id"`      // 房间Id
	Memo        string  `json:"memo"`                             // 备注
}

// @Router /create_hongbao [post]
// @Tags 红包管理
// @Summary 发红包
// @Param x-token header string true "token"
// @Param body body create_hongbao_req true "请求参数"
// @Success 200 "响应数据"
func create_hongbao(ctx *gin.Context) {
	var reqdata create_hongbao_req
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
	tb := xapp.DbQuery().XHongbao
	itb := tb.WithContext(ctx)
	item := &model.XHongbao{
		SellerID:    token.SellerId,
		RoomID:      reqdata.RoomId,
		TotalCount:  reqdata.TotalCount,
		TotalAmount: reqdata.TotalAmount,
		Memo:        reqdata.Memo,
		Sender:      token.Account,
	}
	err := itb.Omit(tb.CreateTime).Create(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	chatdata := new(model.XChat)
	chatdata.SellerID = token.SellerId
	chatdata.RoomID = reqdata.RoomId
	chatdata.Content = fmt.Sprintf("__hongbao__%v", item.ID)
	chatdata.CreateTime = time.Now()
	bytes, _ := json.Marshal(chatdata)
	xapp.Redis().Client().RPush(ctx, "chat_audit", string(bytes)).Result()
	ctx.JSON(http.StatusOK, xenum.Success)
}

type get_hongbao_detail_req struct {
	Id int32 `validate:"required" json:"id"`
}

type get_hongbao_detail_res struct {
	Total int64               `json:"total"` // 总数
	Data  []*model.XHongbaoex `json:"data"`  // 数据
}

// @Router /get_hongbao_detail [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body get_hongbao_detail_req true "请求参数"
// @Success 200  {object} get_hongbao_detail_res "响应数据"
func get_hongbao_detail(ctx *gin.Context) {
	var reqdata get_hongbao_detail_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	response := new(get_hongbao_detail_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XHongbaoex
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
		itb = itb.Where(tb.HongbaoID.Eq(reqdata.Id))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage(0, 1000)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}
