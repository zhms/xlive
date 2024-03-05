package api_admin

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

type ApiAdminRole struct {
	service *service_admin.ServiceAdmin
}

func (this *ApiAdminRole) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.POST("/get_admin_role", this.get_admin_role)
	router.POST("/create_admin_role", middleware.Authorization("系统管理", "角色管理", "增", "新增角色"), this.create_admin_role)
	router.PATCH("/update_admin_role", middleware.Authorization("系统管理", "角色管理", "改", "更新角色"), this.update_admin_role)
	router.DELETE("/delete_admin_role", middleware.Authorization("系统管理", "角色管理", "删", "删除角色"), this.delete_admin_role)
}

// @Router /admin_role/get_admin_role [get]
// @Tags 后台角色
// @Summary 获取角色列表
// @Param x-token header string true "token"
// @Param query query service_admin.GetAdminRoleReq false "筛选参数"
// @Success 200 {object} service_admin.GetAdminRoleRes "成功"
func (this *ApiAdminRole) get_admin_role(ctx *gin.Context) {
	var reqdata service_admin.GetAdminRoleReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetRoleList)
}

// @Router /admin_role/create_admin_role [post]
// @Tags 后台角色
// @Summary 新增角色
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.CreateAdminRoleReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminRole) create_admin_role(ctx *gin.Context) {
	var reqdata service_admin.CreateAdminRoleReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.CreateRole)
}

// @Router /admin_role/update_admin_role [patch]
// @Tags 后台角色
// @Summary 更新角色
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.UpdateAdminRoleReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminRole) update_admin_role(ctx *gin.Context) {
	var reqdata service_admin.UpdateAdminRoleReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateRole)
}

// @Router /admin_role/delete_admin_role [delete]
// @Tags 后台角色
// @Summary 删除角色
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.DeleteAdminRoleReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminRole) delete_admin_role(ctx *gin.Context) {
	var reqdata service_admin.DeleteAdminRoleReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.DeleteRole)
}
