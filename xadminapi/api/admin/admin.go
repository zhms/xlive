package admin

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"xapp/xapp"
	"xapp/xdb"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yinheli/qqwry"
	"gorm.io/gorm"
)

func Init(static *embed.FS) {
	xapp.Db().Exec("call x_init_auth()")
	if static != nil {
		xglobal.Router.Use(Serve("/", EmbedFolder(*static, "www")))
		xglobal.Router.NoRoute(func(c *gin.Context) {
			data, err := static.ReadFile("www/index.html")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		})
	}
	xglobal.ApiV1.POST("/admin_user_login", admin_user_login)
	xglobal.ApiV1.POST("/admin_user_logout", admin_user_logout)
	xglobal.ApiV1.POST("/admin_get_role", admin_get_role)
	xglobal.ApiV1.POST("/admin_create_role", Auth("系统管理", "角色管理", "增", "创建角色"), admin_create_role)
	xglobal.ApiV1.POST("/admin_update_role", Auth("系统管理", "角色管理", "改", "更新角色"), admin_update_role)
	xglobal.ApiV1.POST("/admin_delete_role", Auth("系统管理", "角色管理", "删", "删除角色"), admin_delete_role)
	xglobal.ApiV1.POST("/admin_get_user", Auth("系统管理", "管理员账号", "查", ""), admin_get_user)
	xglobal.ApiV1.POST("/admin_create_user", Auth("系统管理", "管理员账号", "查", "创建管理员"), admin_create_user)
	xglobal.ApiV1.POST("/admin_update_user", Auth("系统管理", "管理员账号", "查", "更新管理员"), admin_update_user)
	xglobal.ApiV1.POST("/admin_delete_user", Auth("系统管理", "管理员账号", "查", "删除管理员"), admin_delete_user)
	xglobal.ApiV1.POST("/admin_get_login_log", Auth("系统管理", "登录日志", "查", ""), admin_get_login_log)
	xglobal.ApiV1.POST("/admin_get_opt_log", Auth("系统管理", "操作日志", "查", ""), admin_get_opt_log)
	if xglobal.IsEnvDev() {
		xglobal.ApiV1.POST("/admin_tools", admin_tools)
	}
}

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}
type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {

		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type TokenData struct {
	SellerId     int32
	Account      string
	UserId       int64
	AuthData     string
	RoleName     string
	GoogleSecret string
}

func DelToken(token string) {
	if token == "" {
		return
	}
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, token)
	_, err := xapp.Redis().Client().Del(context.Background(), rediskey).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func SetToken(token string, value *TokenData) {
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, token)
	valuejson, _ := json.Marshal(value)
	_, err := xapp.Redis().Client().Set(context.Background(), rediskey, string(valuejson), time.Second*3600*24*7).Result()
	if err != nil {
		logs.Error("SetToken error:", err.Error())
	}
}

func GetLocation(ip string) string {
	ipdata := qqwry.NewQQwry("./config/qqwry.dat")
	if ipdata == nil {
		ipdata = qqwry.NewQQwry("./qqwry.dat")
	}
	if ipdata != nil && strings.Index(ip, ".") > 0 {
		ipdata.Find(ip)
		return fmt.Sprintf("%s %s", ipdata.Country, ipdata.City)
	}
	return ""
}

func GetToken(c *gin.Context) *TokenData {
	tokenstr := c.Request.Header.Get("x-token")
	if tokenstr == "" {
		c.JSON(http.StatusBadRequest, xenum.TokenNotFound)
		c.Abort()
		return nil
	}
	rediskey := fmt.Sprintf("%v:token:%s", xglobal.Project, tokenstr)
	value, err := xapp.Redis().Client().Get(c, rediskey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logs.Error("GetToken error:", err.Error())
		}
		c.JSON(http.StatusBadRequest, xenum.GetTokenError)
		c.Abort()
		return nil
	}
	if value == "" {
		c.JSON(http.StatusBadRequest, xenum.TokenExpired)
		c.Abort()
		return nil
	}
	tokendata := new(TokenData)
	json.Unmarshal([]byte(value), tokendata)
	return tokendata
}

func Auth(mainmenu string, submenu string, opt string, optname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokendata := GetToken(c)
		if tokendata == nil {
			return
		}
		if xglobal.IsEnvPrd() && optname != "" {
			VerifyCode := c.Request.Header.Get("VerifyCode")
			if VerifyCode == "" {
				c.JSON(http.StatusBadRequest, xenum.VerifyCodeNotFound)
				c.Abort()
				return
			}
			if tokendata.GoogleSecret == "" {
				c.JSON(http.StatusBadRequest, xenum.VerifySecretNotFound)
				c.Abort()
				return
			}

			if !xutils.VerifyGoogleCode(tokendata.GoogleSecret, VerifyCode) {
				c.JSON(http.StatusBadRequest, xenum.VerifyCodeError)
				c.Abort()
				return
			}
		}
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		mapdata := make(map[string]interface{})
		json.Unmarshal([]byte(tokendata.AuthData), &mapdata)
		m, ok := mapdata[mainmenu]
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized1)
			c.Abort()
			return
		}
		mm, ok := m.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized2)
			c.Abort()
			return
		}
		s, ok := mm[submenu]
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized3)
			c.Abort()
			return
		}
		ms, ok := s.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized4)
			c.Abort()
			return
		}
		o, ok := ms[opt]
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized5)
			c.Abort()
			return
		}
		fo, ok := o.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized6)
			c.Abort()
			return
		}
		if fo != 1 {
			c.JSON(http.StatusUnauthorized, xenum.Unauthorized7)
			c.Abort()
			return
		}

		reqbytes, _ := io.ReadAll(c.Request.Body)
		req := string(reqbytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqbytes))

		c.Next()

		if optname == "" {
			return
		}
		optlog := new(model.XAdminOptLog)
		optlog.SellerID = tokendata.SellerId
		optlog.Account = tokendata.Account
		optlog.ReqData = req
		optlog.ReqPath = c.Request.URL.Path
		optlog.OptName = optname
		optlog.ReqIP = c.ClientIP()
		optlog.ReqIPLocation = GetLocation(optlog.ReqIP)
		tb := xapp.DbQuery().XAdminOptLog
		itb := tb.WithContext(c)
		err := itb.Omit(tb.CreateTime).Create(optlog)
		if err != nil {
			logs.Error("optlog error:", err.Error())
		}
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 用户登录
type admin_user_login_req struct {
	Account  string `validate:"required" json:"account"`  // 账号
	Password string `validate:"required" json:"password"` // 密码
}

type admin_user_login_res struct {
	SellerId   int32  `json:"seller_id"`   // 运营商
	Account    string `json:"account"`     // 账号
	Token      string `json:"token"`       // token
	LoginCount int32  `json:"login_count"` // 登录次数
	AuthData   string `json:"auth_data"`   // 权限数据
	UtcOffset  int    `json:"utc_offset"`  // 当地时区与utc的偏移量
	LoginIp    string `json:"login_ip"`    // 登录Ip
	LoginTime  string `json:"login_time"`  // 登录时间
	Env        string `json:"env"`         // 环境
	LiveUrl    string `json:"live_url"`    // 直播地址
}

// @Router /admin_user_login [post]
// @Tags 管理员
// @Summary 登录
// @Param VerifyCode header string true "验证码"
// @Param body body admin_user_login_req true "请求参数"
// @Success 200 {object} admin_user_login_res "成功"
func admin_user_login(ctx *gin.Context) {
	var reqdata admin_user_login_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	verifycode := ctx.Request.Header.Get("VerifyCode")
	if verifycode == "" {
		ctx.JSON(http.StatusBadRequest, xenum.VerifyCodeNotFound)
		return
	}

	login_locker := "locker:admin_login:" + reqdata.Account
	if !xapp.Redis().Lock(login_locker, 5) {
		ctx.JSON(http.StatusBadRequest, xenum.TooManyRequest)
		return
	}

	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.Account.Eq(reqdata.Account))
	adminuser, err := itb.First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, xenum.UserNotFound)
			return
		}
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}

	if xglobal.IsEnvPrd() && adminuser.LoginGoogle != "" && !xutils.VerifyGoogleCode(adminuser.LoginGoogle, verifycode) {
		ctx.JSON(http.StatusBadRequest, xenum.VerifyCodeError)
		return
	}

	password := xutils.Md5(reqdata.Password)
	if password != adminuser.Password {
		ctx.JSON(http.StatusBadRequest, xenum.UserPasswordError)
		return
	}

	if adminuser.State != 1 {
		ctx.JSON(http.StatusBadRequest, xenum.UserBanded)
		return
	}

	trole := xapp.DbQuery().XAdminRole
	itrole := trole.WithContext(ctx)
	itrole = itrole.Where(trole.SellerID.Eq(adminuser.SellerID), trole.RoleName.Eq(adminuser.RoleName))
	roledata, err := itrole.First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, xenum.RoleNotFound)
			return
		}
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	if roledata.State != 1 {
		ctx.JSON(http.StatusBadRequest, xenum.RoleBaned)
		return
	}
	DelToken(adminuser.Token)
	tokendata := new(TokenData)
	tokendata.SellerId = adminuser.SellerID
	tokendata.UserId = adminuser.ID
	tokendata.Account = reqdata.Account
	tokendata.AuthData = roledata.RoleData
	tokendata.GoogleSecret = adminuser.OptGoogle
	tokendata.RoleName = adminuser.RoleName
	token := uuid.New().String()
	SetToken(token, tokendata)
	response := new(admin_user_login_res)
	response.Account = reqdata.Account
	response.Token = token
	response.AuthData = roledata.RoleData
	response.UtcOffset = xutils.UtcOffset()
	response.LoginCount = adminuser.LoginCount + 1
	response.LoginIp = ctx.ClientIP()
	response.LoginTime = adminuser.LoginTime.Format("2006-01-02 15:04:05")
	response.Env = xglobal.Env
	tkv := xapp.DbQuery().XKv
	itkv := tkv.WithContext(ctx)
	itkv.Select(tkv.V).Where(tkv.K.Eq("client_url")).Scan(&response.LiveUrl)

	tb = xapp.DbQuery().XAdminUser
	itb = tb.WithContext(ctx)
	itb = itb.Where(tb.ID.Eq(adminuser.ID))

	_, err = itb.Updates(map[string]interface{}{
		tb.Token.ColumnName().String():      token,
		tb.LoginIP.ColumnName().String():    ctx.ClientIP(),
		tb.LoginTime.ColumnName().String():  carbon.Now().ToDateTimeString(),
		tb.LoginCount.ColumnName().String(): gorm.Expr(tb.LoginCount.ColumnName().String() + "+1"),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	loginlog := new(model.XAdminLoginLog)
	loginlog.SellerID = tokendata.SellerId
	loginlog.Account = tokendata.Account
	loginlog.Token = token
	loginlog.LoginIP = ctx.ClientIP()
	loginlog.LoginIPLocation = GetLocation(loginlog.LoginIP)
	err = xapp.DbQuery().XAdminLoginLog.Omit(xapp.DbQuery().XAdminLoginLog.CreateTime).Create(loginlog)
	if err != nil {
		logs.Error("loginlog error:", err.Error())
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 用户登出
// @Router /admin_user_logout [post]
// @Tags 管理员
// @Summary 退出
// @Param x-token header string true "token"
// @Success 200 "成功"
func admin_user_logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("x-token")
	DelToken(token)
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 请求角色
type admin_get_role_req struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	RoleName string `json:"role_name"` // 角色名
}

type admin_get_role_res struct {
	Total int64               `json:"total"` // 总数
	Data  []*model.XAdminRole `json:"data"`  // 数据
}

// @Router /admin_get_role [post]
// @Tags 系统管理 - 角色管理
// @Summary 获取角色
// @Param x-token header string true "token"
// @Param body body admin_get_role_req true "请求参数"
// @Success 200 {object} admin_get_role_res "成功"
func admin_get_role(ctx *gin.Context) {
	var reqdata admin_get_role_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := GetToken(ctx)
	response := new(admin_get_role_res)
	tb := xapp.DbQuery().XAdminRole
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if reqdata.RoleName != "" {
		itb = itb.Where(tb.RoleName.Eq(reqdata.RoleName))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 创建角色
type admin_create_role_req struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
	Parent   string `validate:"required" json:"parent"`    // 上级角色
	RoleData string `validate:"required" json:"role_data"` // 权限数据
	State    int32  `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// @Router /admin_create_role [post]
// @Tags 系统管理 - 角色管理
// @Summary 创建角色
// @Param x-token header string true "token"
// @Param body body admin_create_role_req true "请求参数"
// @Success 200  "成功"
func admin_create_role(ctx *gin.Context) {
	var reqdata admin_create_role_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminRole
	itb := tb.WithContext(ctx)
	role := new(model.XAdminRole)
	role.SellerID = token.SellerId
	role.RoleName = reqdata.RoleName
	role.Parent = reqdata.Parent
	role.RoleData = reqdata.RoleData
	role.State = reqdata.State
	role.Memo = reqdata.Memo
	err := itb.Omit(tb.CreateTime).Create(role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 更新角色
type admin_update_role_req struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
	Parent   string `validate:"required" json:"parent"`    // 上级角色
	RoleData string `validate:"required" json:"role_data"` // 权限数据
	State    int    `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// @Router /admin_update_role [post]
// @Tags 系统管理 - 角色管理
// @Summary 更新角色
// @Param x-token header string true "token"
// @Param body body admin_update_role_req true "请求参数"
// @Success 200  "成功"
func admin_update_role(ctx *gin.Context) {
	var reqdata admin_update_role_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminRole
	itb := tb.WithContext(ctx)
	updatedata := make(map[string]interface{})
	if reqdata.Parent != "" {
		updatedata[tb.Parent.ColumnName().String()] = reqdata.Parent
	}
	if reqdata.RoleData != "" {
		updatedata[tb.RoleData.ColumnName().String()] = reqdata.RoleData
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[tb.State.ColumnName().String()] = reqdata.State
	}
	updatedata[tb.Memo.ColumnName().String()] = reqdata.Memo
	itb = itb.Where(tb.SellerID.Eq(token.SellerId), tb.RoleName.Eq(reqdata.RoleName))
	_, err := itb.Updates(updatedata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 删除角色
type admin_delete_role_req struct {
	Id int64 `validate:"required" json:"id"` // 角色Id
}

// @Router /admin_delete_role [post]
// @Tags 系统管理 - 角色管理
// @Summary 删除角色
// @Param x-token header string true "token"
// @Param body body admin_delete_role_req true "请求参数"
// @Success 200  "成功"
func admin_delete_role(ctx *gin.Context) {
	var reqdata admin_delete_role_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminRole
	itb := tb.WithContext(ctx)
	_, err := itb.Where(tb.SellerID.Eq(token.SellerId), tb.ID.Eq(reqdata.Id)).Delete()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 请求管理员
type admin_get_user_req struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Account  string `json:"account"`   // 账号
	RoleName string `json:"role_name"` // 角色名
}

type admin_get_user_res struct {
	Total int64               `json:"total"` // 总数
	Data  []*model.XAdminUser `json:"data"`  // 数据
}

// @Router /admin_get_user [post]
// @Tags 系统管理 - 账号管理
// @Summary 获取账号
// @Param x-token header string true "token"
// @Param body body admin_get_user_req true "请求参数"
// @Success 200 {object} admin_get_user_res "成功"
func admin_get_user(ctx *gin.Context) {
	var reqdata admin_get_user_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := GetToken(ctx)
	response := new(admin_get_user_res)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx).Debug()
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	itb = itb.Where(tb.Agent.Eq(""))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	if reqdata.RoleName != "" {
		itb = itb.Where(tb.RoleName.Eq(reqdata.RoleName))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	for _, v := range response.Data {
		v.Password = ""
		v.Token = ""
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 创建管理员
type admin_create_user_req struct {
	Account  string `validate:"required" json:"account"`   // 账号
	Password string `validate:"required" json:"password"`  // 密码
	RoleName string `validate:"required" json:"role_name"` // 角色
	State    int32  `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// @Router /admin_create_user [post]
// @Tags 系统管理 - 账号管理
// @Summary 创建账号
// @Param x-token header string true "token"
// @Param body body admin_create_user_req true "请求参数"
// @Success 200  "成功"
func admin_create_user(ctx *gin.Context) {
	var reqdata admin_create_user_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	reqdata.RoleName = "超级管理员"
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	user := new(model.XAdminUser)
	user.SellerID = token.SellerId
	user.Account = reqdata.Account
	user.Password = xutils.Md5(reqdata.Password)
	user.RoleName = reqdata.RoleName
	user.State = reqdata.State
	user.Memo = reqdata.Memo
	err := itb.Omit(tb.LoginTime).Create(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 更新管理员
type admin_update_user_req struct {
	Id       int64  `validate:"required" json:"id"` // 管理员Id
	Password string `json:"password"`               // 密码
	RoleName string `json:"role_name"`              // 角色
	State    int    `json:"state"`                  // 状态 1开启,2关闭
	Memo     string `json:"memo"`                   // 备注
}

// @Router /admin_update_user [post]
// @Tags 系统管理 - 账号管理
// @Summary 更新账号
// @Param x-token header string true "token"
// @Param body body admin_update_user_req true "请求参数"
// @Success 200  "成功"
func admin_update_user(ctx *gin.Context) {
	var reqdata admin_update_user_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	updatedata := make(map[string]interface{})
	if reqdata.Password != "" {
		updatedata[tb.Password.ColumnName().String()] = xutils.Md5(reqdata.Password)
	}
	if reqdata.RoleName != "" {
		updatedata[tb.RoleName.CondError().Error()] = reqdata.RoleName
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[tb.State.ColumnName().String()] = reqdata.State
	}
	updatedata[tb.Memo.ColumnName().String()] = reqdata.Memo
	itb = itb.Where(tb.SellerID.Eq(token.SellerId), tb.ID.Eq(reqdata.Id))
	_, err := itb.Updates(updatedata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 删除管理员
type admin_delete_user_req struct {
	Id int64 `validate:"required" json:"id"` // 管理员Id
}

// @Router /admin_delete_user [post]
// @Tags 系统管理 - 账号管理
// @Summary 删除账号
// @Param x-token header string true "token"
// @Param body body admin_delete_user_req true "请求参数"
// @Success 200  "成功"
func admin_delete_user(ctx *gin.Context) {
	var reqdata admin_delete_user_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := GetToken(ctx)
	tb := xapp.DbQuery().XAdminUser
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId), tb.ID.Eq(reqdata.Id))
	_, err := itb.Delete()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 获取登录日志
type admin_get_login_log_req struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	Account   string `json:"account"`    // 操作人
	LoginIp   string `json:"login_ip"`   // 登录Ip
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type admin_get_login_log_res struct {
	Total int64                   `json:"total"` // 总数
	Data  []*model.XAdminLoginLog `json:"data"`  // 数据
}

// @Router /admin_get_login_log [post]
// @Tags 系统管理 - 登录日志
// @Summary 获取登录日志
// @Param x-token header string true "token"
// @Param body body admin_get_opt_log_req true "请求参数"
// @Success 200 {object} admin_get_opt_log_res "成功"
func admin_get_login_log(ctx *gin.Context) {
	var reqdata admin_get_login_log_req
	if err := ctx.ShouldBind(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := GetToken(ctx)
	response := new(admin_get_login_log_res)
	tb := xapp.DbQuery().XAdminLoginLog
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	if reqdata.LoginIp != "" {
		itb = itb.Where(tb.LoginIP.Eq(reqdata.LoginIp))
	}
	if reqdata.StartTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", reqdata.StartTime)
		itb = itb.Where(tb.CreateTime.Gte(t))
	}
	if reqdata.EndTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", reqdata.EndTime)
		itb = itb.Where(tb.CreateTime.Lt(t))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 获取操作日志
type admin_get_opt_log_req struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	Account   string `json:"account"`    // 操作人
	OptName   string `json:"opt_name"`   // 操作名
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type admin_get_opt_log_res struct {
	Total int64                 `json:"total"` // 总数
	Data  []*model.XAdminOptLog `json:"data"`  // 数据
}

// @Router /admin_get_opt_log [post]
// @Tags 系统管理 - 操作日志
// @Summary 获取操作日志
// @Param x-token header string true "token"
// @Param body body admin_get_opt_log_req true "请求参数"
// @Success 200 {object} admin_get_opt_log_res "成功"
func admin_get_opt_log(ctx *gin.Context) {
	var reqdata admin_get_opt_log_req
	if err := ctx.ShouldBind(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := GetToken(ctx)
	response := new(admin_get_opt_log_res)
	tb := xapp.DbQuery().XAdminOptLog
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	if reqdata.OptName != "" {
		itb = itb.Where(tb.OptName.Eq(reqdata.OptName))
	}
	if reqdata.StartTime != "" {
		itb = itb.Where(tb.CreateTime.Gte(carbon.Parse(reqdata.StartTime).StdTime()))
	}
	if reqdata.EndTime != "" {
		itb = itb.Where(tb.CreateTime.Lt(carbon.Parse(reqdata.EndTime).StdTime()))
	}
	var err error
	itb = itb.Order(tb.ID.Desc())
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 系统工具
type admin_tools_req struct {
	QueryType string `validate:"required" json:"query_type"` // 查询类型
}

type admin_tools_res struct {
	Data string `json:"data"` // 数据
}

func get_db_field(db *gorm.DB, dbname string, result *map[string]int) {
	db = db.Raw("show tables")
	d, _ := xdb.Find(db)
	d.ForEach(func(item *xutils.XMap) bool {
		item.ForEach(func(k string, v interface{}) bool {
			t := v.(string)
			db = xapp.Db().Raw("select COLUMN_NAME from information_schema.COLUMNS where table_name = ? and table_schema = ?", t, dbname)
			tb, _ := xdb.Find(db)
			tb.ForEach(func(tbitem *xutils.XMap) bool {
				tbitem.ForEach(func(tk string, tv interface{}) bool {
					(*result)[tv.(string)] = 1
					return true
				})
				return true
			})
			return true
		})
		return true
	})
}

func admin_tools(ctx *gin.Context) {
	var reqdata admin_tools_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	response := new(admin_tools_res)
	if reqdata.QueryType == "db_table" {
		db := xapp.Db().Raw("show tables")
		d, _ := xdb.Find(db)
		response.Data = "const (\r\n"
		d.ForEach(func(item *xutils.XMap) bool {
			item.ForEach(func(k string, v interface{}) bool {
				t := v.(string)
				t = strings.Replace(t, "_", " ", -1)
				t = strings.Title(t)
				t = strings.Replace(t, " ", "", -1)
				response.Data += fmt.Sprintf("\tTable%s=\"%s\"\r\n", t, v)
				return true
			})
			return true
		})
		response.Data += ")\r\n"
	}
	if reqdata.QueryType == "db_field" {
		response.Data = "const (\r\n"
		fields := make(map[string]int)
		get_db_field(xapp.Db(), "x_live", &fields)
		for k := range fields {
			kx := strings.Replace(k, "_", " ", -1)
			kx = strings.Title(kx)
			kx = strings.Replace(kx, " ", "", -1)
			response.Data += fmt.Sprintf("\t%s=\"%s\"\r\n", kx, k)
		}
		response.Data += ")\r\n"
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}
