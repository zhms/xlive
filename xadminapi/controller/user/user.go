package controller_user

import (
	"net/http"
	"xadminapi/server"
	"xadminapi/service"
	service_user "xadminapi/service/user"
	"xcom/enum"

	"github.com/gin-gonic/gin"
)

type ControllerUser struct {
	service *service_user.ServiceUser
}

func (this *ControllerUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceUser
	router.GET("/get_user", this.get_user)
}

// @Router /user/get_user [get]
// @Tags 会员管理
// @Summary 会员管理
// @Param x-token header string true "token"
// @Param query query service_user.GetUserReq false "筛选参数"
// @Success 200 {object} []model.XUser "成功"
func (this *ControllerUser) get_user(ctx *gin.Context) {
	var reqdata service_user.GetUserReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	token := server.GetToken(ctx)
	if token == nil {
		return
	}
	reqdata.SellerId = token.SellerId
	total, data, merr, err := this.service.GetUserList(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakePageSucess(total, data))
}
