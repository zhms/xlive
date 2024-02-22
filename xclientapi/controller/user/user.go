package controller_user

import (
	"net/http"
	"strings"
	"xclientapi/service"
	service_user "xclientapi/service/user"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ControllerUser struct {
	service *service_user.ServiceUser
}

func (this *ControllerUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceUser
	router.POST("/user_login", this.user_login)
}

// @Router /user/user_login [post]
// @Tags 玩家
// @Summary 玩家登录
// @Param body body service_user.UserLoginReq true "body参数"
// @Success 200 {object} service_user.UserLoginRes "成功"
func (this *ControllerUser) user_login(ctx *gin.Context) {
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
	host := ctx.Request.Host
	host = strings.Replace(host, "www.", "", -1)
	host = strings.Split(host, ":")[0]
	appid := ctx.GetHeader("appid")
	reponse, merr, err := this.service.UserLogin(appid, host, ctx.ClientIP(), &reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakeSucess(reponse))
}
