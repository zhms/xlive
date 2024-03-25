package inout_chart

import (
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
)

func Init() {
	xglobal.ApiV1.POST("/get_inout_data", admin.Auth("数据分析", "在线图表", "查", ""), get_inout_data)
}

type get_inout_data_req struct {
	StartTime string `json:"start_time" validate:"required"` // 开始时间
	EndTime   string `json:"end_time" validate:"required"`   // 结束时间
}

type get_inout_data_res struct {
	Total int64               `json:"total"` // 总数
	Data  []*model.XStatistic `json:"data"`  // 数据
}

// @Router /get_inout_data [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body get_inout_data_req true "请求参数"
// @Success 200  {object} get_inout_data_res "响应数据"
func get_inout_data(ctx *gin.Context) {
	var reqdata get_inout_data_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	response := new(get_inout_data_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XStatistic
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	itb = itb.Where(tb.RecordType.Eq("el"))
	{
		itb = itb.Where(tb.CreateTime.Gte(carbon.Parse(reqdata.StartTime).StdTime()))
		itb = itb.Where(tb.CreateTime.Lt(carbon.Parse(reqdata.EndTime).StdTime()))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage(0, 10000)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}
