package user

import (
	"fmt"
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb"
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
	Total int64       `json:"total"` // 总数
	Data  []xdb.XUser `json:"data"`  // 数据
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
	db := xapp.Db().Omit(xdb.Password).Model(&xdb.XUser{})
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.Account != "" {
		db = db.Where(xdb.Account+xdb.EQ, reqdata.Account)
	}
	if reqdata.Agent != "" {
		db = db.Where(xdb.Agent+xdb.EQ, reqdata.Agent)
	}
	if reqdata.LoginIp != "" {
		db = db.Where(xdb.LoginIp+xdb.EQ, reqdata.LoginIp)
	}
	err := db.Count(&response.Total).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(xdb.Id + xdb.DESC)
	err = db.Find(&response.Data).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Export == 1 {
		e := excel.NewExcelBuilder("会员列表_" + fmt.Sprint(xutils.Now()))
		defer e.Close()
		e.SetTitle(xdb.Id, "Id")
		e.SetTitle(xdb.Account, "账号")
		e.SetTitle(xdb.State, "状态")
		e.SetTitle(xdb.ChatState, "聊天状态")
		e.SetTitle(xdb.Agent, "业务员")
		e.SetTitle(xdb.LoginIp, "登录Ip")
		e.SetTitle(xdb.LoginLocation, "登录地点")
		e.SetTitle(xdb.LoginCount, "登录时间")
		e.SetTitle(xdb.LoginTime, "登录时间")
		e.SetTitle(xdb.CreateTime, "注册时间")
		for i, v := range response.Data {
			e.SetValue(xdb.Id, v.Id, int64(i+2))
			e.SetValue(xdb.Account, v.Account, int64(i+2))
			if v.State == 1 {
				e.SetValue(xdb.State, "正常", int64(i+2))
			} else {
				e.SetValue(xdb.State, "禁用", int64(i+2))
			}
			if v.ChatState == 1 {
				e.SetValue(xdb.ChatState, "正常", int64(i+2))
			} else {
				e.SetValue(xdb.ChatState, "禁言", int64(i+2))
			}
			e.SetValue(xdb.Agent, v.Agent, int64(i+2))
			e.SetValue(xdb.LoginIp, v.LoginIp, int64(i+2))
			e.SetValue(xdb.LoginLocation, v.LoginIpLocation, int64(i+2))
			e.SetValue(xdb.LoginCount, v.LoginCount, int64(i+2))
			e.SetValue(xdb.LoginTime, v.LoginTime, int64(i+2))
			e.SetValue(xdb.CreateTime, v.CreateTime, int64(i+2))
		}
		e.Write(ctx)
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_user_req struct {
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
	//response := new(create_user_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}

type update_user_req struct {
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
	//response := new(update_user_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}
