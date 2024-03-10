package api_user

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_user "xadminapi/service/user"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiUser struct {
	service *service_user.ServiceUser
}

func (this *ApiUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceUser
	router.POST("/get_user", middleware.Authorization("会员管理", "会员管理", "查", ""), this.get_user)
	router.POST("/add_user", middleware.Authorization("会员管理", "会员管理", "增", "添加会员"), this.add_user)
	router.POST("/update_user", middleware.Authorization("会员管理", "会员管理", "改", "更新会员"), this.update_user)
}

// @Router /user/get_user [post]
// @Tags 会员管理
// @Summary 会员列表
// @Param x-token header string true "token"
// @Param body body service_user.GetUserReq true "body参数"
// @Success 200 {object} service_user.GetUserRes "成功"
func (this *ApiUser) get_user(ctx *gin.Context) {
	var reqdata service_user.GetUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetUserList)
}

// @Router /user/add_user [post]
// @Tags 会员管理
// @Summary 添加会员
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_user.AddUserReq true "body参数"
// @Success 200 "成功"
func (this *ApiUser) add_user(ctx *gin.Context) {
	var reqdata service_user.AddUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.AddUser)
}

// @Router /user/update_user [post]
// @Tags 会员管理
// @Summary 更新会员
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_user.UpdateUserReq true "body参数"
// @Success 200 "成功"
func (this *ApiUser) update_user(ctx *gin.Context) {
	var reqdata service_user.UpdateUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateUser)
}
