package sales

import (
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_sales", admin.Auth("业务员管理", "业务员管理", "查", ""), get_sales)
	xglobal.ApiV1.POST("/create_sales", admin.Auth("业务员管理", "业务员管理", "增", "创建业务员"), create_sales)
	xglobal.ApiV1.POST("/update_sales", admin.Auth("业务员管理", "业务员管理", "改", "更新业务员"), update_sales)
	xglobal.ApiV1.POST("/delete_sales", admin.Auth("业务员管理", "业务员管理", "删", "删除业务员"), delete_sales)
}

type get_sales_req struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Account  string `json:"account"`   // 账号
	RoleName string `json:"role_name"` // 角色名
}

type get_sales_res struct {
	Total int64               `json:"total"` // 总数
	Data  []*model.XAdminUser `json:"data"`  // 数据
}

// @Router /get_sales [post]
// @Tags 业务员管理
// @Summary 获取业务员
// @Param x-token header string true "token"
// @Param body body get_sales_req true "请求参数"
// @Success 200  {object} get_sales_res "响应数据"
func get_sales(ctx *gin.Context) {
	var reqdata get_sales_req
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
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := admin.GetToken(ctx)
	response := new(get_sales_res)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if token.RoleName != "超级管理员" {
		itb = itb.Where(tb.Agent.Eq(token.Account))
	} else {
		itb = itb.Where(tb.Agent.Neq(""))
	}
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	if reqdata.RoleName != "" {
		itb = itb.Where(tb.RoleName.Eq(reqdata.RoleName))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	for _, v := range response.Data {
		v.Password = ""
		v.Token = ""
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_sales_req struct {
	Account  string `validate:"required" json:"account"`   // 账号
	Password string `validate:"required" json:"password"`  // 密码
	RoleName string `validate:"required" json:"role_name"` // 角色
	State    int32  `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
	RoomId   int32  `json:"room_id"`                       // 房间Id
}

// @Router /create_sales [post]
// @Tags 业务员管理
// @Summary 创建业务员
// @Param x-token header string true "token"
// @Param body body create_sales_req true "请求参数"
// @Success 200 "响应数据"
func create_sales(ctx *gin.Context) {
	var reqdata create_sales_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	reqdata.RoleName = "业务员"
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	user := new(model.XAdminUser)
	user.SellerID = token.SellerId
	user.Account = reqdata.Account
	user.Password = xutils.Md5(reqdata.Password)
	user.RoleName = reqdata.RoleName
	user.State = reqdata.State
	user.Memo = reqdata.Memo
	user.Agent = token.Account
	user.RoomID = reqdata.RoomId
	err := itb.Omit(tb.LoginTime).Create(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type update_sales_req struct {
	Id       int64  `validate:"required" json:"id"` // 管理员Id
	Password string `json:"password"`               // 密码
	Memo     string `json:"memo"`                   // 备注
	RoomId   int32  `json:"room_id"`                // 房间Id
}

// @Router /update_sales [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body update_sales_req true "请求参数"
// @Success 200 "响应数据"
func update_sales(ctx *gin.Context) {
	var reqdata update_sales_req
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
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.ID.Eq(reqdata.Id))
	if token.RoleName != "超级管理员" {
		itb = itb.Where(tb.Agent.Eq(token.Account))
	}
	updatedata := map[string]interface{}{}
	updatedata["memo"] = reqdata.Memo
	if reqdata.Password != "" {
		updatedata["password"] = xutils.Md5(reqdata.Password)
	}
	if reqdata.RoomId > 0 {
		updatedata["room_id"] = reqdata.RoomId
	}
	_, err := itb.Updates(updatedata)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type delete_sales_req struct {
	Id int64 `validate:"required" json:"id"` // 管理员Id
}

// @Router /delete_sales [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body delete_sales_req true "请求参数"
// @Success 200 "响应数据"
func delete_sales(ctx *gin.Context) {
	var reqdata delete_sales_req
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
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.ID.Eq(reqdata.Id))
	if token.RoleName != "超级管理员" {
		itb = itb.Where(tb.Agent.Eq(token.Account))
	}
	_, err := itb.Delete()
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
