package robot

import (
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

func Init() {
	xglobal.ApiV1.POST("/get_robot", admin.Auth("机器人管理", "机器人管理", "查", ""), get_robot)
	xglobal.ApiV1.POST("/create_robot", admin.Auth("机器人管理", "机器人管理", "增", "创建机器人"), create_robot)
	xglobal.ApiV1.POST("/update_robot", admin.Auth("机器人管理", "机器人管理", "改", "更新机器人"), update_robot)
	xglobal.ApiV1.POST("/delete_robot", admin.Auth("机器人管理", "机器人管理", "删", "删除机器人"), delete_robot)
	xglobal.ApiV1.POST("/get_robot_count", admin.Auth("机器人管理", "机器人管理", "查", ""), get_robot_count)
	xglobal.ApiV1.POST("/update_robot_count", admin.Auth("机器人管理", "机器人管理", "改", "更新机器人数量"), update_robot_count)
}

type get_robot_req struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Account  string `json:"account"`   // 账号
}

type get_robot_res struct {
	Total int64           `json:"total"` // 总数
	Data  []*model.XRobot `json:"data"`  // 数据
}

// @Router /get_robot [post]
// @Tags 机器人管理
// @Summary 获取机器人
// @Param x-token header string true "token"
// @Param body body get_robot_req true "请求参数"
// @Success 200  {object} get_robot_res "响应数据"
func get_robot(ctx *gin.Context) {
	var reqdata get_robot_req
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
	response := new(get_robot_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_robot_req struct {
	Account string `json:"account" validate:"required"` // 账号
}

// @Router /create_robot [post]
// @Tags 机器人管理
// @Summary 创建机器人
// @Param x-token header string true "token"
// @Param body body create_robot_req true "请求参数"
// @Success 200 "响应数据"
func create_robot(ctx *gin.Context) {
	var reqdata create_robot_req
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
	item := new(model.XRobot)
	item.SellerID = token.SellerId
	{
		item.Account = reqdata.Account
	}
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx)
	err := itb.Omit(tb.CreateTime).Create(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type update_robot_req struct {
	Id      int32  `json:"id" validate:"required"` // id
	Account string `json:"account"`                // 账号
}

// @Router /update_robot [post]
// @Tags 机器人管理
// @Summary 更新机器人
// @Param x-token header string true "token"
// @Param body body update_robot_req true "请求参数"
// @Success 200 "响应数据"
func update_robot(ctx *gin.Context) {
	var reqdata update_robot_req
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
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
		itb = itb.Where(tb.ID.Eq(reqdata.Id))
	}
	item := map[string]interface{}{}
	{
		if reqdata.Account != "" {
			item[tb.Account.ColumnName().String()] = reqdata.Account
		}
	}
	_, err := itb.Updates(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type delete_robot_req struct {
	Id int32 `json:"id" validate:"required"` // id
}

// @Router /delete_robot [post]
// @Tags 机器人管理
// @Summary 删除机器人
// @Param x-token header string true "token"
// @Param body body delete_robot_req true "请求参数"
// @Success 200 "响应数据"
func delete_robot(ctx *gin.Context) {
	var reqdata delete_robot_req
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
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
		itb = itb.Where(tb.ID.Eq(reqdata.Id))
	}
	_, err := itb.Delete()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type get_robot_count_req struct {
}

type get_robot_count_res struct {
	Count int32 `json:"count"` // 数量
}

// @Router /get_robot_count [post]
// @Tags 机器人管理
// @Summary 获取机器人数量
// @Param x-token header string true "token"
// @Param body body get_robot_count_req true "请求参数"
// @Success 200  {object} get_robot_count_res "响应数据"
func get_robot_count(ctx *gin.Context) {
	var reqdata get_robot_count_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	response := new(get_robot_count_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XKv
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
		itb = itb.Where(tb.K.Eq("robot_count"))
	}
	pkv, err := itb.First()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	response.Count = cast.ToInt32(pkv.V)
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type update_robot_count_req struct {
	Count int32 `json:"count"` // 数量
}

// @Router /update_robot_count [post]
// @Tags 机器人管理
// @Summary 更新机器人数量
// @Param x-token header string true "token"
// @Param body body update_robot_count_req true "请求参数"
// @Success 200 "响应数据"
func update_robot_count(ctx *gin.Context) {
	var reqdata update_robot_count_req
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
	tb := xapp.DbQuery().XKv
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
		itb = itb.Where(tb.K.Eq("robot_count"))
	}
	item := map[string]interface{}{}
	{
		item[tb.V.ColumnName().String()] = reqdata.Count
	}
	_, err := itb.Updates(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
