package user

import (
	"fmt"
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	excel "xapp/xexcel"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_user", admin.Auth("会员管理", "会员管理", "查", ""), get_user)
	xglobal.ApiV1.POST("/create_user", admin.Auth("会员管理", "会员管理", "增", "创建会员"), create_user)
	xglobal.ApiV1.POST("/update_user", admin.Auth("会员管理", "会员管理", "改", "更新会员"), update_user)
}

type get_user_req struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量

	Account string `json:"account"`  // 账号
	Agent   string `json:"agent"`    // 代理
	LoginIp string `json:"login_ip"` // 登录Ip
	Export  int    `json:"export"`   // 导出
}

type get_user_res struct {
	Total int64          `json:"total"` // 总数
	Data  []*model.XUser `json:"data"`  // 数据
}

// @Router /get_user [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body get_user_req true "请求参数"
// @Success 200  {object} get_user_res "响应数据"
func get_user(ctx *gin.Context) {
	var reqdata get_user_req
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
	if reqdata.Export == 1 {
		reqdata.PageSize = 100000
	}
	response := new(get_user_res)
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(int32(token.SellerId)))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	if reqdata.Agent != "" {
		itb = itb.Where(tb.Agent.Eq(reqdata.Agent))
	}
	if reqdata.LoginIp != "" {
		itb = itb.Where(tb.LoginIP.Eq(reqdata.LoginIp))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Export == 1 {
		e := excel.NewExcelBuilder("会员列表_" + fmt.Sprint(xutils.Now()))
		defer e.Close()
		e.SetTitle(tb.ID.ColumnName().String(), "Id")
		e.SetTitle(tb.Account.ColumnName().String(), "账号")
		e.SetTitle(tb.State.ColumnName().String(), "状态")
		e.SetTitle(tb.ChatState.ColumnName().String(), "聊天状态")
		e.SetTitle(tb.Agent.ColumnName().String(), "业务员")
		e.SetTitle(tb.LoginIP.ColumnName().String(), "登录Ip")
		e.SetTitle(tb.LoginIPLocation.ColumnName().String(), "登录地点")
		e.SetTitle(tb.LoginCount.ColumnName().String(), "登录时间")
		e.SetTitle(tb.LoginTime.CondError().Error(), "登录时间")
		e.SetTitle(tb.CreateTime.CondError().Error(), "注册时间")
		for i, v := range response.Data {
			e.SetValue(tb.ID.ColumnName().String(), v.ID, int64(i+2))
			e.SetValue(tb.Account.ColumnName().String(), v.Account, int64(i+2))
			if v.State == 1 {
				e.SetValue(tb.State.ColumnName().String(), "正常", int64(i+2))
			} else {
				e.SetValue(tb.State.ColumnName().String(), "禁用", int64(i+2))
			}
			if v.ChatState == 1 {
				e.SetValue(tb.ChatState.ColumnName().String(), "正常", int64(i+2))
			} else {
				e.SetValue(tb.ChatState.ColumnName().String(), "禁言", int64(i+2))
			}
			e.SetValue(tb.Agent.ColumnName().String(), v.Agent, int64(i+2))
			e.SetValue(tb.LoginIP.ColumnName().String(), v.LoginIP, int64(i+2))
			e.SetValue(tb.LoginIPLocation.ColumnName().String(), v.LoginIPLocation, int64(i+2))
			e.SetValue(tb.LoginCount.ColumnName().String(), v.LoginCount, int64(i+2))
			e.SetValue(tb.LoginTime.ColumnName().String(), v.LoginTime, int64(i+2))
			e.SetValue(tb.CreateTime.ColumnName().String(), v.CreateTime, int64(i+2))
		}
		e.Write(ctx)
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_user_req struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}

// @Router /create_user [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body create_user_req true "请求参数"
// @Success 200 "响应数据"
func create_user(ctx *gin.Context) {
	var reqdata create_user_req
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
	reqdata.Password = xutils.Md5(reqdata.Password)
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(ctx)
	itb.Create(&model.XUser{SellerID: int32(token.SellerId), Account: reqdata.Account, Password: reqdata.Password})
	ctx.JSON(http.StatusOK, xenum.Success)
}

type update_user_req struct {
	Id        int    `json:"id" validate:"required"` // Id
	Password  string `json:"password" `              // 密码
	State     int    `json:"state" `                 // 状态
	ChatState int    `json:"chat_state" `            // 聊天状态
}

// @Router /update_user [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body update_user_req true "请求参数"
// @Success 200 "响应数据"
func update_user(ctx *gin.Context) {
	var reqdata update_user_req
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
	tb := xapp.DbQuery().XUser
	itb := tb.WithContext(ctx)
	updatedata := make(map[string]interface{})
	if reqdata.Password != "" {
		updatedata[tb.Password.ColumnName().String()] = reqdata.Password
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[tb.State.ColumnName().String()] = reqdata.State
	}
	if reqdata.ChatState != 0 {
		updatedata[tb.ChatState.ColumnName().String()] = reqdata.ChatState
	}
	itb.Where(tb.SellerID.Eq(int32(token.SellerId)), tb.ID.Eq(int32(reqdata.Id))).Updates(updatedata)
	ctx.JSON(http.StatusOK, xenum.Success)
}
