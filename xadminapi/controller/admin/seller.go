package controller_admin

import (
	"net/http"
	"xadminapi/middleware"
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
	router.GET("/get_seller", this.get_seller)
	router.POST("/create_seller", middleware.Authorization("系统管理", "运营商管理", "增", "新增运营商"), this.create_seller)
	router.PATCH("/update_seller", middleware.Authorization("系统管理", "运营商管理", "改", "更新运营商"), this.update_seller)
	router.DELETE("/delete_seller", middleware.Authorization("系统管理", "运营商管理", "删", "删除运营商"), this.delete_seller)
}

// @Router /seller/get_seller [get]
// @Tags 运营商
// @Summary 获取运营商列表
// @Param x-token header string true "token"
// @Param query query service_admin.GetXSellerReq false "筛选参数"
// @Success 200 {object} []model.XSeller "成功"
func (this *ControllerAdminSeller) get_seller(ctx *gin.Context) {
	var reqdata service_admin.GetXSellerReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	total, data, merr, err := this.service.GetSellerList(&reqdata)
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
	merr, err := this.service.CreateSeller(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.Success)
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
	merr, err := this.service.UpdateSeller(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.Success)
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
	rows, merr, err := this.service.DeleteSeller(&reqdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.InternalError, err.Error()))
		return
	}
	if merr != nil {
		ctx.JSON(http.StatusBadRequest, merr)
		return
	}
	ctx.JSON(http.StatusOK, enum.MakeSucess(rows))
}
