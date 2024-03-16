package live_ban

import (
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
	xglobal.ApiV1.POST("/get_ip_ban", admin.Auth("直播间", "Ip封禁", "查", ""), get_ip_ban)
	xglobal.ApiV1.POST("/delete_ip_ban", admin.Auth("直播间", "Ip封禁", "改", "解封Ip"), delete_ip_ban)
}

type get_ip_ban_req struct {
	Page     int `json:"page"`
	PageSize int `json:"size"`
}

type get_ip_ban_res struct {
	Total int64          `json:"total"`
	Data  []*model.XChat `json:"data"`
}

// @Router /get_ip_ban [post]
// @Tags 直播间 - Ip封禁
// @Summary 获取封禁Ip
// @Param x-token header string true "token"
// @Param body body get_ip_ban_req true "请求参数"
// @Success 200 {object} get_ip_ban_res "响应数据"
func get_ip_ban(ctx *gin.Context) {
	var reqdata get_ip_ban_req
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

	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	response := new(get_ip_ban_res)
	tb := xapp.DbQuery().XChat
	itb := tb.WithContext(ctx)
	itb.Where(tb.SellerID.Eq(int32(token.SellerId)))
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type delete_ip_ban_req struct {
	Id int32 `json:"id"`
}

// @Router /delete_ip_ban [post]
// @Tags 直播间 - Ip封禁
// @Summary 解封Ip
// @Param x-token header string true "token"
// @Param body body delete_ip_ban_req true "请求参数"
// @Success 200 "响应数据"
func delete_ip_ban(ctx *gin.Context) {
	var reqdata delete_ip_ban_req
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
	tb := xapp.DbQuery().XChatBanIP
	itb := tb.WithContext(ctx)
	itb.Where(tb.SellerID.Eq(int32(token.SellerId)))
	banip, err := itb.First()
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	_, err = itb.Delete(&model.XChatBanIP{ID: reqdata.Id})
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	xapp.Redis().Client().SRem(ctx, "ip_ban", banip.IP)
	ctx.JSON(http.StatusOK, xenum.Success)
}
