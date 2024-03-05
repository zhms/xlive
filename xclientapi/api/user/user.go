package api_user

import (
	"net/http"
	"xclientapi/server"
	"xclientapi/service"
	service_user "xclientapi/service/user"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiUser struct {
	service *service_user.ServiceUser
}

func (this *ApiUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceUser
	router.POST("/user_login", this.user_login)
}

// @Router /user/user_login [post]
// @Tags 玩家
// @Summary 玩家登录
// @Param body body service_user.UserLoginReq true "body参数"
// @Success 200 {object} service_user.UserLoginRes "成功"
func (this *ApiUser) user_login(ctx *gin.Context) {
	var reqdata service_user.UserLoginReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.UserLogin)
}
