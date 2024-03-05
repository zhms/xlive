package controller_admin

import (
	"net/http"
	"xadminapi/middleware"
	"xadminapi/server"
	"xadminapi/service"
	service_admin "xadminapi/service/admin"
	"xcom/enum"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

type ControllerAdminSeller struct {
	service *service_admin.ServiceAdmin
}

// 初始化路由
func (this *ControllerAdminSeller) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.POST("/get_seller", this.get_seller)
	router.POST("/create_seller", middleware.Authorization("系统管理", "运营商管理", "增", "新增运营商"), this.create_seller)
	router.PATCH("/update_seller", middleware.Authorization("系统管理", "运营商管理", "改", "更新运营商"), this.update_seller)
	router.DELETE("/delete_seller", middleware.Authorization("系统管理", "运营商管理", "删", "删除运营商"), this.delete_seller)
}

// @Router /seller/get_seller [get]
// @Tags 运营商
// @Summary 获取运营商列表
// @Param x-token header string true "token"
// @Param query query service_admin.GetXSellerReq false "筛选参数"
// @Success 200 {object} service_admin.GetXSellerRes "成功"
func (this *ControllerAdminSeller) get_seller(ctx *gin.Context) {
	var reqdata service_admin.GetXSellerReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetSellerList)
}

// @Router /seller/create_seller [post]
// @Tags 运营商
// @Summary 新增运营商
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.CreateXSellerReq true "body参数"
// @Success 200 "成功"
func (this *ControllerAdminSeller) create_seller(ctx *gin.Context) {
	var reqdata service_admin.CreateXSellerReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.CreateSeller)
}

// @Router /seller/update_seller [patch]
// @Tags 运营商
// @Summary 更新运营商
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.UpdateXSellerReq true "body参数"
// @Success 200 "成功"
func (this *ControllerAdminSeller) update_seller(ctx *gin.Context) {
	var reqdata service_admin.UpdateXSellerReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateSeller)
}

// @Router /seller/delete_seller [delete]
// @Tags 运营商
// @Summary 删除运营商
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.DeleteXSellerReq true "body参数"
// @Success 200 "成功"
func (this *ControllerAdminSeller) delete_seller(ctx *gin.Context) {
	var reqdata service_admin.DeleteXSellerReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.DeleteSeller)
}
