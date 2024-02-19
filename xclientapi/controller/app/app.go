package controller_app

import (
	"net/http"
	"strings"
	"xclientapi/service"
	service_app "xclientapi/service/app"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ControllerApp struct {
	service *service_app.ServiceApp
}

func (this *ControllerApp) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceApp
	router.GET("/get_live_info", this.get_live_info)
	router.GET("/get_online_info", this.get_online_info)
}

// @Router /app/get_live_info [get]
// @Tags 应用
// @Summary 获取直播信息
// @Param x-token header string true "token"
// @Param body body service_app.AppGetLiveInfoReq false "筛选参数"
// @Success 200 {object} service_app.AppGetLiveInfoRes "成功"
func (this *ControllerApp) get_live_info(ctx *gin.Context) {
	var reqdata service_app.AppGetLiveInfoReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
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
	reponse, merr, err := this.service.GetLiveInfo(host, &reqdata)
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

// @Router /app/get_online_info [get]
// @Tags 应用
// @Summary 获取在线人数
// @Param x-token header string true "token"
// @Param body body service_app.AppGetOnlineInfoReq false "筛选参数"
// @Success 200 {object} service_app.AppGetOnlineInfoRes "成功"
func (this *ControllerApp) get_online_info(ctx *gin.Context) {
	var reqdata service_app.AppGetOnlineInfoReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
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
	reponse, merr, err := this.service.GetOnlineInfo(host, &reqdata)
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
