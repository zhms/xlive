package api_live

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_live "xadminapi/service/live"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ApiLiveRoom struct {
	service *service_live.ServiceLiveRoom
}

func (this *ApiLiveRoom) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceLiveRoom
	router.POST("/get_live_room_id", this.get_live_room_id)
	router.POST("/get_live_room", middleware.Authorization("直播间", "直播间列表", "查", ""), this.get_live_room)
	router.POST("/create_live_room", middleware.Authorization("直播间", "直播间列表", "增", "创建直播间"), this.create_live_room)
	router.POST("/update_live_room", middleware.Authorization("直播间", "直播间列表", "改", "修改直播间"), this.update_live_room)
	router.POST("/delete_live_room", middleware.Authorization("直播间", "直播间列表", "删", "删除直播间"), this.delete_live_room)
}

func (this *ApiLiveRoom) get_live_room_id(ctx *gin.Context) {
	var reqdata service_live.GetLiveRoomIdReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetLiveRoomId)
}

// @Router /live_room/get_live_room [post]
// @Tags 直播间
// @Summary 直播间列表
// @Param x-token header string true "token"
// @Param body body service_live.GetLiveRoomListReq true "body参数"
// @Success 200 {object} service_live.GetLiveRoomListRes "成功"
func (this *ApiLiveRoom) get_live_room(ctx *gin.Context) {
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

// @Router /live_room/create_live_room [post]
// @Tags 直播间
// @Summary 创建直播间
// @Param x-token header string true "token"
// @Param body body service_live.CreateLiveRoomReq true "body参数"
// @Success 200 "成功"
func (this *ApiLiveRoom) create_live_room(ctx *gin.Context) {
	var reqdata service_live.CreateLiveRoomReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.CreateLiveRoom)
}

// @Router /live_room/update_live_room [post]
// @Tags 直播间
// @Summary 更新直播间
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_live.UpdateLiveRoomReq true "body参数"
// @Success 200 "成功"
func (this *ApiLiveRoom) update_live_room(ctx *gin.Context) {
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

// @Router /live_room/delete_live_room [post]
// @Tags 直播间
// @Summary 删除直播间
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_live.DeleteLiveRoomReq true "body参数"
// @Success 200 "成功"
func (this *ApiLiveRoom) delete_live_room(ctx *gin.Context) {
	var reqdata service_live.DeleteLiveRoomReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.DeleteLiveRoom)
}
