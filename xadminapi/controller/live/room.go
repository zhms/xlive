package controller_live

import (
	"net/http"
	"xadminapi/server"
	"xadminapi/service"
	service_live "xadminapi/service/live"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ControllerLiveRoom struct {
	service *service_live.ServiceLiveRoom
}

func (this *ControllerLiveRoom) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveRoom
	router.POST("/get_live_room", this.get_live_room)
	router.POST("/create_live_room", this.create_live_room)
	router.PATCH("/update_live_room", this.update_live_room)
}

// @Router /live_room/get_live_room [post]
// @Tags 直播间
// @Summary 直播间列表
// @Param x-token header string true "token"
// @Param body body service_live.GetLiveRoomListReq true "body参数"
// @Success 200 {object} service_live.GetLiveRoomListRes "成功"
func (this *ControllerLiveRoom) get_live_room(ctx *gin.Context) {
	var reqdata service_live.GetLiveRoomListReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetLiveRoomList)
}

func (this *ControllerLiveRoom) create_live_room(c *gin.Context) {

}

// @Router /live_room/update_live_room [patch]
// @Tags 直播间
// @Summary 更新直播间
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.UpdateLiveRoomListReq true "body参数"
// @Success 200 "成功"
func (this *ControllerLiveRoom) update_live_room(ctx *gin.Context) {
	var reqdata service_live.UpdateLiveRoomReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateLiveRoom)
}
