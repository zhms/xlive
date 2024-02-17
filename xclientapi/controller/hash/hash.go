package controller_hash

import (
	"net/http"
	"strings"
	"xclientapi/middleware"
	"xclientapi/server"
	"xclientapi/service"
	service_hash "xclientapi/service/hash"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ControllerHash struct {
	service *service_hash.ServiceHash
}

func (this *ControllerHash) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceHash
	router.GET("/lottery_next", this.lottery_next)
	router.POST("/hash_bet", middleware.Login(), this.hash_bet)
	router.POST("/lottery_bet", middleware.Login(), this.lottery_bet)
}

// @Router /hash/hash_bet [post]
// @Tags 哈希
// @Summary 哈希下注
// @Param body body service_hash.HashBetReq true "body参数"
// @Success 200 {object} service_hash.HashBetRes "成功"
func (this *ControllerHash) hash_bet(ctx *gin.Context) {
	var reqdata service_hash.HashBetReq
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

	token := server.GetToken(ctx)
	reponse, merr, err := this.service.HashBet(token, host, &reqdata)
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

// @Router /hash/lottery_next [get]
// @Tags 哈希
// @Summary 彩票下一期
// @Param body body service_hash.LotteryNextReq true "body参数"
// @Success 200 {object} service_hash.LotteryNextRes "成功"
func (this *ControllerHash) lottery_next(ctx *gin.Context) {
	var reqdata service_hash.LotteryNextReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	reponse, merr, err := this.service.LotteryNext(&reqdata)
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

// @Router /hash/lottery_bet [post]
// @Tags 哈希
// @Summary 彩票下注
// @Param body body service_hash.LotteryBetReq true "body参数"
// @Success 200 {object} service_hash.LotteryBetRes "成功"
func (this *ControllerHash) lottery_bet(ctx *gin.Context) {
	var reqdata service_hash.LotteryBetReq
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

	token := server.GetToken(ctx)
	reponse, merr, err := this.service.LotteryBet(token, host, &reqdata)
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
