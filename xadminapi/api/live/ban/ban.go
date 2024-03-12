package api_live_ban

import (
	"net/http"
	"xapp/xenum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {

}

type get_ip_ban_req struct {
}

type get_ip_ban_res struct {
}

// @Router /get_ip_ban [post]
// @Tags 直播间 - Ip封禁
// @Summary 获取封禁Ip
// @Param x-token header string true "token"
// @Param body body get_ip_ban_req true "请求参数"
// @Success 200 {object} get_ip_ban_res "响应数据"
func get_ip_ban(ctx *gin.Context) {
	var reqdata get_ip_ban_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(get_ip_ban_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}

type create_ip_ban_req struct {
}

// @Router /create_ip_ban [post]
// @Tags 直播间 - Ip封禁
// @Summary 封禁Ip
// @Param x-token header string true "token"
// @Param body body create_ip_ban_req true "请求参数"
// @Success 200 "响应数据"
func create_ip_ban(ctx *gin.Context) {
	var reqdata create_ip_ban_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(create_ip_ban_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}

type delete_ip_ban_req struct {
}

// @Router /delete_ip_ban [post]
// @Tags 直播间 - Ip封禁
// @Summary 解封Ip
// @Param x-token header string true "token"
// @Param body body delete_ip_ban_req true "请求参数"
// @Success 200 "响应数据"
func delete_ip_ban(ctx *gin.Context) {
	var reqdata delete_ip_ban_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(delete_ip_ban_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}
