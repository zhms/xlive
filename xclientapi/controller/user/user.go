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
	//router.POST("/user_register", this.user_register)
	//router.POST("/user_test", this.user_test)
}

// @Router /user/user_register [post]
// @Tags 玩家
// @Summary 玩家注册
// @Param body body service_user.UserRegisterReq true "body参数"
// @Success 200 {object} service_user.UserRegisterRes "成功"
func (this *ControllerUser) user_register(ctx *gin.Context) {
	var reqdata service_user.UserRegisterReq
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
	reponse, merr, err := this.service.UserRegister(host, ctx.ClientIP(), &reqdata)
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
	reponse, merr, err := this.service.UserLogin(host, ctx.ClientIP(), &reqdata)
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

// @Router /user/user_test [post]
// @Tags 玩家
// @Summary 玩家测试
// @Param body body service_user.UserTestReq true "body参数"
// @Success 200 {object} service_user.UserTestRes "成功"
func (this *ControllerUser) user_test(ctx *gin.Context) {
	var reqdata service_user.UserTestReq
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
	reponse, merr, err := this.service.UserTest(&reqdata)
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
