package api_live_chat

import (
	"net/http"
	"xapp/xenum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {

}

type get_chat_req struct {
}

type get_chat_res struct {
}

// @Router /get_chat [post]
// @Tags 直播间 - 互动列表
// @Summary 获取互动列表
// @Param x-token header string true "token"
// @Param body body get_chat_req true "请求参数"
// @Success 200  {object} get_chat_res "响应数据"
func get_chat(ctx *gin.Context) {
	var reqdata get_chat_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(get_chat_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}

type update_chat_req struct {
	Id    int `json:"id" validate:"required"`
	State int `json:"state" validate:"required"` // 1:未审核 2:通过 3:拒绝 4:封Ip
}

// @Router /update_chat [post]
// @Tags 直播间 - 互动列表
// @Summary 审核互动列表
// @Param x-token header string true "token"
// @Param body body update_chat_req true "请求参数"
// @Success 200 "响应数据"
func update_chat(ctx *gin.Context) {
	var reqdata update_chat_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(update_chat_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}
