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
	"xapp/xenum"
	"xapp/xglobal"
	"xapp/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
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
	xglobal.ApiV1.POST("/admin_get_role", Auth("系统管理", "角色管理", "查", ""), admin_get_role)
	xglobal.ApiV1.POST("/admin_create_role", Auth("系统管理", "角色管理", "增", "创建角色"), admin_create_role)
	xglobal.ApiV1.POST("/admin_update_role", Auth("系统管理", "角色管理", "改", "更新角色"), admin_update_role)
	xglobal.ApiV1.POST("/admin_delete_role", Auth("系统管理", "角色管理", "删", "删除角色"), admin_delete_role)
	xglobal.ApiV1.POST("/admin_get_user", Auth("系统管理", "账号管理", "查", ""), admin_get_user)
	xglobal.ApiV1.POST("/admin_create_user", Auth("系统管理", "账号管理", "查", "创建管理员"), admin_create_user)
	xglobal.ApiV1.POST("/admin_update_user", Auth("系统管理", "账号管理", "查", "更新管理员"), admin_update_user)
	xglobal.ApiV1.POST("/admin_delete_user", Auth("系统管理", "账号管理", "查", "删除管理员"), admin_delete_user)
	xglobal.ApiV1.POST("/admin_get_login_log", Auth("系统管理", "登录日志", "查", ""), admin_get_login_log)
	xglobal.ApiV1.POST("/admin_get_opt_log", Auth("系统管理", "操作日志", "查", ""), admin_get_opt_log)
	if !xglobal.IsEnvPrd() {
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

type XAdminLoginLog struct {
	Id              uint64 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增Id'" json:"id"`                           // 自增Id
	SellerId        int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                       // 运营商
	Account         string `gorm:"column:account;type:varchar(32);comment:'账号'" json:"account"`                           // 账号
	Token           string `gorm:"column:token;type:varchar(64);comment:'登录的token'" json:"-"`                             // 登录的token
	LoginIp         string `gorm:"column:login_ip;type:varchar(32);comment:'最近一次登录Ip'" json:"login_ip"`                   // 登录Ip
	LoginIpLocation string `gorm:"column:login_ip_location;type:varchar(64);comment:'登录Ip地理位置'" json:"login_ip_location"` // 登录Ip地理位置
	Memo            string `gorm:"column:memo;type:varchar(256);comment:'备注'" json:"memo"`                                // 备注
	CreateTime      string `gorm:"column:create_time" json:"create_time"`                                                 // 创建时间
}

func (XAdminLoginLog) TableName() string {
	return "x_admin_login_log"
}

type XAdminOptLog struct {
	Id            uint64 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增Id'" json:"id"`                                       // 自增Id
	SellerId      int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                                   // 运营商
	Account       string `gorm:"column:account;type:varchar(32);charset:utf8mb4;comment:'账号'" json:"account"`                       // 账号
	ReqPath       string `gorm:"column:req_path;type:varchar(256);charset:utf8mb4;comment:'请求路径'" json:"req_path"`                  // 请求路径
	OptName       string `gorm:"column:opt_name;type:varchar(64);charset:utf8mb4;comment:'操作名称'" json:"opt_name"`                   // 请求路径
	ReqData       string `gorm:"column:req_data;type:varchar(256);charset:utf8mb4;comment:'请求参数'" json:"req_data"`                  // 请求参数
	ReqIp         string `gorm:"column:req_ip;type:varchar(32);charset:utf8mb4;comment:'请求的Ip'" json:"req_ip"`                      // 请求的Ip
	ReqIpLocation string `gorm:"column:req_ip_location;type:varchar(64);charset:utf8mb4;comment:'请求Ip地理位置'" json:"req_ip_location"` // 请求Ip地理位置
	CreateTime    string `gorm:"column:create_time" json:"create_time"`                                                             // 创建时间
}

func (XAdminOptLog) TableName() string {
	return "x_admin_opt_log"
}

type XAdminUser struct {
	Id          uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT;comment:'自增Id'" json:"id"`                     // 自增Id
	SellerId    int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                   // 运营商
	Account     string `gorm:"column:account;type:varchar(32);charset:utf8mb4;comment:'账号'" json:"account"`       // 账号
	Password    string `gorm:"column:password;type:varchar(64);charset:utf8mb4;comment:'登录密码'" json:"-"`          // 登录密码
	RoleName    string `gorm:"column:role_name;type:varchar(32);charset:utf8mb4;comment:'角色'" json:"role_name"`   // 角色
	LoginGoogle string `gorm:"column:login_google;type:varchar(32);charset:utf8mb4;comment:'登录验证码'" json:"-"`     // 登录验证码
	OptGoogle   string `gorm:"column:opt_google;type:varchar(32);charset:utf8mb4;comment:'渠道商'" json:"-"`         // 渠道商
	State       int    `gorm:"column:state;comment:'状态 1开启,2关闭'" json:"state"`                                    // 状态 1开启,2关闭
	Token       string `gorm:"column:token;type:varchar(255);charset:utf8mb4;comment:'最后登录的token'" json:"-"`      // 最后登录的token
	LoginCount  int    `gorm:"column:login_count;comment:'登录次数'" json:"login_count"`                              // 登录次数
	LoginTime   string `gorm:"column:login_time;default:CURRENT_TIMESTAMP;comment:'最后登录时间'" json:"login_time"`    // 最后登录时间
	LoginIp     string `gorm:"column:login_ip;type:varchar(32);charset:utf8mb4;comment:'最后登录Ip'" json:"login_ip"` // 最后登录Ip
	Memo        string `gorm:"column:memo;type:varchar(256);charset:utf8mb4;comment:'备注'" json:"memo"`            // 备注
	CreateTime  string `gorm:"column:create_time" json:"create_time"`                                             // 创建时间
}

func (XAdminUser) TableName() string {
	return "x_admin_user"
}

type XAdminRole struct {
	Id         uint64 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增Id'" json:"id"`      // 自增Id
	SellerId   int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                  // 运营商
	RoleName   string `gorm:"column:role_name;type:varchar(32);comment:'角色名'" json:"role_name"` // 角色名
	Parent     string `gorm:"column:parent;type:varchar(32);comment:'上级角色'" json:"parent"`      // 上级角色
	RoleData   string `gorm:"column:role_data;type:text;comment:'权限数据'" json:"role_data"`       // 权限数据
	State      int    `gorm:"column:state;comment:'状态 1开启,2关闭'" json:"state"`                   // 状态 1开启,2关闭
	Memo       string `gorm:"column:memo;type:varchar(256);comment:'备注'" json:"memo"`           // 备注
	CreateTime string `gorm:"column:create_time" json:"create_time"`                            // 创建时间
}

func (XAdminRole) TableName() string {
	return "x_admin_role"
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
	SellerId     int
	Account      string
	UserId       int
	AuthData     string
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
	value, err := xapp.Redis().Client().Get(context.Background(), rediskey).Result()
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

		optlog := new(XAdminOptLog)
		optlog.SellerId = tokendata.SellerId
		optlog.Account = tokendata.Account
		optlog.ReqData = req
		optlog.ReqPath = c.Request.URL.Path
		optlog.OptName = optname
		optlog.ReqIp = c.ClientIP()
		optlog.ReqIpLocation = GetLocation(optlog.ReqIp)
		err := xapp.Db().Omit(xdb.CreateTime).Create(&optlog).Error
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
	SellerId   int    `json:"seller_id"`   // 运营商
	Account    string `json:"account"`     // 账号
	Token      string `json:"token"`       // token
	LoginCount int    `json:"login_count"` // 登录次数
	AuthData   string `json:"auth_data"`   // 权限数据
	UtcOffset  int    `json:"utc_offset"`  // 当地时区与utc的偏移量
	LoginIp    string `json:"login_ip"`    // 登录Ip
	LoginTime  string `json:"login_time"`  // 登录时间
	Env        string `json:"env"`         // 环境
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

	var adminuser XAdminUser
	db := xapp.Db()
	db = db.Where(xdb.Account+xdb.EQ, reqdata.Account)
	err := db.First(&adminuser).Error
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

	if adminuser.State != xdb.StateYes {
		ctx.JSON(http.StatusBadRequest, xenum.UserBanded)
		return
	}

	roledata := new(XAdminRole)
	db = xapp.Db()
	db = db.Where(xdb.SellerId+xdb.EQ, adminuser.SellerId)
	db = db.Where(xdb.RoleName+xdb.EQ, adminuser.RoleName)
	err = db.First(&roledata).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, xenum.RoleNotFound)
			return
		}
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	if roledata.State != xdb.StateYes {
		ctx.JSON(http.StatusBadRequest, xenum.RoleBaned)
		return
	}
	DelToken(adminuser.Token)
	tokendata := new(TokenData)
	tokendata.SellerId = adminuser.SellerId
	tokendata.UserId = int(adminuser.Id)
	tokendata.Account = reqdata.Account
	tokendata.AuthData = roledata.RoleData
	tokendata.GoogleSecret = adminuser.OptGoogle
	token := uuid.New().String()
	SetToken(token, tokendata)
	response := new(admin_user_login_res)
	response.Account = reqdata.Account
	response.Token = token
	response.AuthData = roledata.RoleData
	response.UtcOffset = xutils.UtcOffset()
	response.LoginCount = adminuser.LoginCount + 1
	response.LoginIp = ctx.ClientIP()
	response.LoginTime = adminuser.LoginTime
	response.Env = xglobal.Env

	db = xapp.Db().Model(new(XAdminUser))
	db = db.Where(xdb.Id+xdb.EQ, adminuser.Id)
	err = db.Updates(map[string]interface{}{
		xdb.Token:      token,
		xdb.LoginIp:    ctx.ClientIP(),
		xdb.LoginTime:  xutils.Now(),
		xdb.LoginCount: gorm.Expr(xdb.LoginCount+xdb.PLUS, 1),
	}).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	loginlog := new(XAdminLoginLog)
	loginlog.SellerId = tokendata.SellerId
	loginlog.Account = tokendata.Account
	loginlog.Token = token
	loginlog.LoginIp = ctx.ClientIP()
	loginlog.CreateTime = xutils.Now()
	loginlog.LoginIpLocation = GetLocation(loginlog.LoginIp)
	db = xapp.Db().Create(&loginlog)
	if db.Error != nil {
		logs.Error("loginlog error:", db.Error.Error())
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
	Total int64        `json:"total"` // 总数
	Data  []XAdminRole `json:"data"`  // 数据
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
	db := xapp.Db().Model(new(XAdminRole))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.RoleName != "" {
		db = db.Where(xdb.RoleName+xdb.EQ, reqdata.RoleName)
	}
	db.Count(&response.Total)
	db = db.Limit(reqdata.PageSize).Offset((reqdata.Page - 1) * reqdata.PageSize)
	db.Find(&response.Data)
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 创建角色
type admin_create_role_req struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
	Parent   string `validate:"required" json:"parent"`    // 上级角色
	RoleData string `validate:"required" json:"role_data"` // 权限数据
	State    int    `validate:"required" json:"state"`     // 状态 1开启,2关闭
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
	role := new(XAdminRole)
	role.SellerId = token.SellerId
	role.RoleName = reqdata.RoleName
	role.Parent = reqdata.Parent
	role.RoleData = reqdata.RoleData
	role.State = reqdata.State
	role.Memo = reqdata.Memo
	role.CreateTime = xutils.Now()
	db := xapp.Db().Create(&role)
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
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
	updatedata := make(map[string]interface{})
	if reqdata.Parent != "" {
		updatedata[xdb.Parent] = reqdata.Parent
	}
	if reqdata.RoleData != "" {
		updatedata[xdb.RoleData] = reqdata.RoleData
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[xdb.State] = reqdata.State
	}
	updatedata[xdb.Memo] = reqdata.Memo
	db := xapp.Db().Model(new(XAdminRole))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	db = db.Where(xdb.RoleName+xdb.EQ, reqdata.RoleName)
	db = db.Updates(updatedata)
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 删除角色
type admin_delete_role_req struct {
	Id int `validate:"required" json:"id"` // 角色Id
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
	db := xapp.Db().Model(new(XAdminRole))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	db = db.Where(xdb.Id+xdb.EQ, reqdata.Id)
	db = db.Delete(&XAdminRole{})
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
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
	Total int64        `json:"total"` // 总数
	Data  []XAdminUser `json:"data"`  // 数据
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
	db := xapp.Db().Model(new(XAdminUser))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.Account != "" {
		db = db.Where(xdb.Account+xdb.EQ, reqdata.Account)
	}
	if reqdata.RoleName != "" {
		db = db.Where(xdb.RoleName+xdb.EQ, reqdata.RoleName)
	}
	db.Count(&response.Total)
	db = db.Limit(reqdata.PageSize).Offset((reqdata.Page - 1) * reqdata.PageSize)
	db.Find(&response.Data)
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

// 创建管理员
type admin_create_user_req struct {
	Account  string `validate:"required" json:"account"`   // 账号
	Password string `validate:"required" json:"password"`  // 密码
	RoleName string `validate:"required" json:"role_name"` // 角色
	State    int    `validate:"required" json:"state"`     // 状态 1开启,2关闭
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
	token := GetToken(ctx)
	user := new(XAdminUser)
	user.SellerId = token.SellerId
	user.Account = reqdata.Account
	user.Password = xutils.Md5(reqdata.Password)
	user.RoleName = reqdata.RoleName
	user.State = reqdata.State
	user.Memo = reqdata.Memo
	user.CreateTime = xutils.Now()
	db := xapp.Db().Omit(xdb.LoginTime).Create(&user)
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 更新管理员
type admin_update_user_req struct {
	Id       int    `validate:"required" json:"id"` // 管理员Id
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
	updatedata := make(map[string]interface{})
	if reqdata.Password != "" {
		updatedata[xdb.Password] = xutils.Md5(reqdata.Password)
	}
	if reqdata.RoleName != "" {
		updatedata[xdb.RoleName] = reqdata.RoleName
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[xdb.State] = reqdata.State
	}
	updatedata[xdb.Memo] = reqdata.Memo
	db := xapp.Db().Model(new(XAdminUser))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	db = db.Where(xdb.Id+xdb.EQ, reqdata.Id)
	db = db.Updates(updatedata)
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

// 删除管理员
type admin_delete_user_req struct {
	Id int `validate:"required" json:"id"` // 管理员Id
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
	db := xapp.Db().Model(new(XAdminUser))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	db = db.Where(xdb.Id+xdb.EQ, reqdata.Id)
	db = db.Delete(&XAdminUser{})
	if db.Error != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, db.Error.Error()))
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
	Total int64            `json:"total"` // 总数
	Data  []XAdminLoginLog `json:"data"`  // 数据
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
	db := xapp.Db().Model(new(XAdminLoginLog))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.Account != "" {
		db = db.Where(xdb.Account+xdb.EQ, reqdata.Account)
	}
	if reqdata.LoginIp != "" {
		db = db.Where(xdb.LoginIp+xdb.EQ, reqdata.LoginIp)
	}
	if reqdata.StartTime != "" {
		db = db.Where(xdb.CreateTime+xdb.GTE, reqdata.StartTime)
	}
	if reqdata.EndTime != "" {
		db = db.Where(xdb.CreateTime+xdb.LTE, reqdata.EndTime)
	}
	db.Count(&response.Total)
	db = db.Limit(reqdata.PageSize).Offset((reqdata.Page - 1) * reqdata.PageSize).Order(xdb.Id + xdb.DESC)
	db.Find(&response.Data)
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
	Total int64          `json:"total"` // 总数
	Data  []XAdminOptLog `json:"data"`  // 数据
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
	db := xapp.Db().Model(new(XAdminOptLog))
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	if reqdata.Account != "" {
		db = db.Where(xdb.Account+xdb.EQ, reqdata.Account)
	}
	if reqdata.OptName != "" {
		db = db.Where(xdb.OptName+xdb.EQ, reqdata.OptName)
	}
	if reqdata.StartTime != "" {
		db = db.Where(xdb.CreateTime+xdb.GTE, reqdata.StartTime)
	}
	if reqdata.EndTime != "" {
		db = db.Where(xdb.CreateTime+xdb.LTE, reqdata.EndTime)
	}
	db.Count(&response.Total)
	db = db.Limit(reqdata.PageSize).Offset((reqdata.Page - 1) * reqdata.PageSize).Order(xdb.Id + xdb.DESC)
	db.Find(&response.Data)
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
