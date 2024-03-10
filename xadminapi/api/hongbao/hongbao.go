package api_hongbao

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_hongbao "xadminapi/service/hongbao"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiHongbao struct {
	service *service_hongbao.ServiceHongbao
}

func (this *ApiHongbao) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceHongbao
	router.POST("/get_hongbao_list", middleware.Authorization("红包管理", "红包管理", "查", ""), this.get_hongbao_list)
	router.POST("/send_hongbao", middleware.Authorization("红包管理", "红包管理", "发红包", "发红包"), this.send_hongbao)
}

func (this *ApiHongbao) get_hongbao_list(ctx *gin.Context) {
	var reqdata service_hongbao.GetHongbaoListReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetHongbaoList)
}

func (this *ApiHongbao) send_hongbao(ctx *gin.Context) {
	var reqdata service_hongbao.SendHongbaoReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.SendHongbao)
}
