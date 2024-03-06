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

type ApiAdminUser struct {
	service *service_admin.ServiceAdmin
}

// 初始化路由
func (this *ApiAdminUser) InitRouter(router *gin.RouterGroup) {
	this.service = &service.Entries().ServiceAdmin
	router.POST("/user_login", this.user_login)
	router.POST("/user_logout", this.user_logout)
	router.POST("/get_admin_user", middleware.Authorization("系统管理", "账号管理", "查", ""), this.get_admin_user)
	router.POST("/create_admin_user", middleware.Authorization("系统管理", "账号管理", "增", "新增管理员"), this.create_admin_user)
	router.POST("/update_admin_user", middleware.Authorization("系统管理", "账号管理", "改", "更新管理员"), this.update_admin_user)
	router.POST("/delete_admin_user", middleware.Authorization("系统管理", "账号管理", "删", "删除管理员"), this.delete_admin_user)
	router.POST("/set_login_googlesecret", middleware.Authorization("系统管理", "账号管理", "设置登录验证码", "设置登歌验证码"), this.set_login_googlesecret)
	router.POST("/set_opt_googlesecret", middleware.Authorization("系统管理", "账号管理", "设置操作验证码", "设置操作验证码"), this.set_opt_googlesecret)
}

// @Router /admin_user/user_login [post]
// @Tags 后台登录
// @Summary 管理员登录
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.AdminUserLoginReq true "body参数"
// @Success 200 {object} service_admin.AdminUserLoginRes "成功"
func (this *ApiAdminUser) user_login(ctx *gin.Context) {
	var reqdata service_admin.AdminUserLoginReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.AdminUserLogin)
}

// @Router /admin_user/user_logout [post]
// @Tags 后台登录
// @Summary 管理员退出
// @Param x-token header string true "token"
// @Success 200 "成功"
func (this *ApiAdminUser) user_logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("x-token")
	server.DelToken(token)
	ctx.JSON(http.StatusOK, enum.Success)
}

// @Router /admin_user/get_admin_user [post]
// @Tags 后台用户
// @Summary 获取管理员列表
// @Param x-token header string true "token"
// @Param query query service_admin.GetAdminUserReq false "筛选参数"
// @Success 200 {object} service_admin.GetAdminUserRes "成功"
func (this *ApiAdminUser) get_admin_user(ctx *gin.Context) {
	var reqdata service_admin.GetAdminUserReq
	if err := ctx.ShouldBindQuery(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.GetAdminUserList)
}

// @Router /admin_user/create_admin_user [post]
// @Tags 后台用户
// @Summary 新增管理员账号
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.CreateAdminUserReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminUser) create_admin_user(ctx *gin.Context) {
	var reqdata service_admin.CreateAdminUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.CreateAdminUser)
}

// @Router /admin_user/update_admin_user [post]
// @Tags 后台用户
// @Summary 更新管理员账号
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.UpdateAdminUserReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminUser) update_admin_user(ctx *gin.Context) {
	var reqdata service_admin.UpdateAdminUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.UpdateAdminUser)
}

// @Router /admin_user/delete_admin_user [post]
// @Tags 后台用户
// @Summary 删除管理员账号
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.DeleteAdminUserReq true "body参数"
// @Success 200 "成功"
func (this *ApiAdminUser) delete_admin_user(ctx *gin.Context) {
	var reqdata service_admin.DeleteAdminUserReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequest(ctx, reqdata, this.service.DeleteAdminUser)
}

// @Router /admin_user/set_login_googlesecret [post]
// @Tags 后台用户
// @Summary 设置登录验证码
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.SetLoginGoogleReq true "body参数"
// @Success 200 {object}  service_admin.SetLoginGoogleRes "成功"
func (this *ApiAdminUser) set_login_googlesecret(ctx *gin.Context) {
	var reqdata service_admin.SetLoginGoogleReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.SetLoginGoogle)
}

// @Router /admin_user/set_opt_googlesecret [post]
// @Tags 后台用户
// @Summary 设置操作验证码
// @Param x-token header string true "token"
// @Param VerifyCode header string true "验证码"
// @Param body body service_admin.SetOptGoogleReq true "body参数"
// @Success 200 {object}  service_admin.SetOptGoogleRes "成功"
func (this *ApiAdminUser) set_opt_googlesecret(ctx *gin.Context) {
	var reqdata service_admin.SetOptGoogleReq
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, enum.MakeError(enum.BadParams, err.Error()))
		return
	}
	server.OnRequestEx(ctx, reqdata, this.service.SetOptGoogle)
}
