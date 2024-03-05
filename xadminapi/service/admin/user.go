package service_admin

import (
	"fmt"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/global"
	"xcom/xutils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminUserLoginReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}

type AdminUserLoginRes struct {
	SellerId   int    `json:"seller_id"`   // 运营商
	Account    string `json:"account"`     // 账号
	Token      string `json:"token"`       // token
	LoginCount int    `json:"login_count"` // 登录次数
	AuthData   string `json:"auth_data"`   // 权限数据
	UtcOffset  int    `json:"utc_offset"`  // 当地时区与utc的偏移量
	LoginIp    string `json:"login_ip"`    // 登录Ip
	LoginTime  string `json:"login_time"`  // 登录时间
}

// 管理员登录
func (this *ServiceAdmin) AdminUserLogin(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(AdminUserLoginReq)
	ip := ctx.ClientIP()
	verifycode := ctx.Request.Header.Get("VerifyCode")

	if global.IsEnvPrd() && verifycode == "" {
		return nil, enum.VerifyNotFoundCode, nil
	}
	if global.IsEnvPrd() && len(verifycode) != 6 {
		return nil, enum.VerifyInCorrectCode, nil
	}

	locker := enum.Lock_AdminLogin + reqdata.Account
	if !server.Redis().Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}

	var adminuser model.XAdminUser
	sellerid := 1
	channelid := 0
	db := server.Db()
	db = db.Where(edb.SellerId+edb.EQ, sellerid)
	db = db.Where(edb.ChannelId+edb.EQ, channelid)
	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	err = db.First(&adminuser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, enum.UserNotFound, nil
		}
		return nil, nil, err
	}

	if global.IsEnvPrd() && adminuser.LoginGoogle != "" && !xutils.VerifyGoogleCode(adminuser.LoginGoogle, verifycode) {
		return nil, enum.VerifyInCorrectCode, nil
	}

	password := xutils.Md5(reqdata.Password)
	if password != adminuser.Password {
		return nil, enum.UserPasswordError, nil
	}

	if adminuser.State != enum.StateYes {
		return nil, enum.UserStateError, nil
	}

	roledata := model.XAdminRole{}
	db = server.Db().Model(&model.XAdminRole{})
	db = db.Where(edb.SellerId+edb.EQ, sellerid)
	db = db.Where(edb.RoleName+edb.EQ, adminuser.RoleName)
	err = db.First(&roledata).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, enum.RoleNotFound, nil
		}
		return nil, nil, err
	}
	if roledata.State != enum.StateYes {
		return nil, enum.RoleStateError, nil
	}
	server.DelToken(adminuser.Token)

	tokendata := server.TokenData{}
	tokendata.SellerId = sellerid
	tokendata.UserId = int(adminuser.Id)
	tokendata.Account = reqdata.Account
	tokendata.AuthData = roledata.RoleData
	tokendata.GoogleSecret = adminuser.OptGoogle
	token := uuid.New().String()
	server.SetToken(token, &tokendata)
	response := &AdminUserLoginRes{}
	response.SellerId = sellerid
	response.Account = reqdata.Account
	response.Token = token
	response.AuthData = roledata.RoleData
	response.UtcOffset = xutils.UtcOffset()
	response.LoginCount = adminuser.LoginCount + 1
	response.LoginIp = ip
	response.LoginTime = adminuser.LoginTime

	db = server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.Id+edb.EQ, adminuser.Id)
	err = db.Updates(map[string]interface{}{
		edb.Token:      token,
		edb.LoginIp:    ip,
		edb.LoginTime:  xutils.Now(),
		edb.LoginCount: gorm.Expr(edb.LoginCount+edb.PLUS, 1),
	}).Error
	if err != nil {
		return nil, nil, err
	}
	loginlog := model.XAdminLoginLog{}
	loginlog.SellerId = tokendata.SellerId
	loginlog.Account = tokendata.Account
	loginlog.CreateTime = xutils.Now()
	loginlog.Token = token
	loginlog.LoginIp = ip
	db = server.Db().Model(&loginlog).Create(&loginlog)
	if db.Error != nil {
		logs.Error("loginlog create error", db.Error)
		return nil, nil, db.Error
	}
	return response, nil, nil
}

type GetAdminUserReq struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Account  string `json:"account"`   // 账号
}

type GetAdminUserRes struct {
	Total int64              `json:"total"` // 总数
	Data  []model.XAdminUser `json:"data"`  // 数据
}

// 获取管理员列表
func (this *ServiceAdmin) GetAdminUserList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetAdminUserReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	db := server.Db().Model(&model.XAdminUser{})
	db = xutils.DbWhere(db, edb.SellerId+edb.EQ, token.SellerId, int(0))
	db = xutils.DbWhere(db, edb.Account+edb.EQ, reqdata.Account, "")

	data := GetAdminUserRes{}

	err = db.Count(&data.Total).Error
	if err != nil {
		return nil, nil, err
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return nil, nil, err
	}
	return data, nil, err
}

type CreateAdminUserReq struct {
	Account  string `validate:"required" json:"account"`   // 账号
	Password string `validate:"required" json:"password"`  // 登录密码
	RoleName string `validate:"required" json:"role_name"` // 角色
	State    int    `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// 创建管理员
func (this *ServiceAdmin) CreateAdminUser(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(CreateAdminUserReq)
	token := server.GetToken(ctx)
	exists, err := this.role_exists(token.SellerId, reqdata.RoleName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return enum.RoleNotFound, nil
	}
	err = server.Db().Model(&model.XAdminUser{}).Create(map[string]interface{}{
		edb.SellerId: token.SellerId,
		edb.Account:  reqdata.Account,
		edb.Password: xutils.Md5(reqdata.Password),
		edb.RoleName: reqdata.RoleName,
		edb.State:    reqdata.State,
		edb.Memo:     reqdata.Memo,
	}).Error
	return nil, err
}

type UpdateAdminUserReq struct {
	Account  string `validate:"required" json:"account"` // 账号
	Password string `json:"password"`                    // 登录密码
	RoleName string `json:"role_name"`                   // 角色
	State    int    `json:"state"`                       // 状态 1开启,2关闭
	Memo     string `json:"memo"`                        // 备注
}

// 更新管理员
func (this *ServiceAdmin) UpdateAdminUser(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateAdminUserReq)
	updatedata := map[string]interface{}{}
	token := server.GetToken(ctx)
	updatedata[edb.Memo] = reqdata.Memo
	xutils.MapSet(&updatedata, edb.RoleName, reqdata.RoleName, "")
	if reqdata.RoleName != "" {
		exists, err := this.role_exists(token.SellerId, reqdata.RoleName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return enum.RoleNotFound, nil
		}
		updatedata[edb.RoleName] = reqdata.RoleName
	}
	if reqdata.State == 1 || reqdata.State == 2 {
		updatedata[edb.State] = reqdata.State
	}
	db := server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	db = db.Updates(updatedata)
	return nil, db.Error
}

type DeleteAdminUserReq struct {
	Account string `validate:"required" json:"account"` // 账号
}

// 删除管理员
func (this *ServiceAdmin) DeleteAdminUser(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(DeleteAdminUserReq)
	if reqdata.Account == "admin" {
		return enum.UserCantDelete, nil
	}
	token := server.GetToken(ctx)
	db := server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	db = db.Delete(&model.XAdminUser{})
	return nil, db.Error
}

type SetLoginGoogleReq struct {
	Account string `validate:"required" json:"account"` // 账号
}

type SetLoginGoogleRes struct {
	Url string `json:"url"` // 二维码
}

// 设置登录验证码
func (this *ServiceAdmin) SetLoginGoogle(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(SetLoginGoogleReq)
	token := server.GetToken(ctx)
	locker := fmt.Sprintf(enum.Lock_ChangeGoogleSecret, token.SellerId, reqdata.Account)
	if !server.Redis().Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}
	verifycode := ctx.Request.Header.Get("VerifyCode")
	tokendata := server.GetToken(ctx)

	userdata := model.XAdminUser{}
	db := server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, tokendata.SellerId)
	db = db.Where(edb.Account+edb.EQ, tokendata.Account)
	db = db.First(&userdata)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, enum.UserNotFound, nil
		}
		return nil, nil, db.Error
	}
	if userdata.LoginGoogle == "" && userdata.LoginGoogle == reqdata.Account {
		return nil, enum.VerifyNotFoundSecret, nil
	}

	if userdata.OptGoogle != "" {
		if global.IsEnvPrd() && !xutils.VerifyGoogleCode(userdata.OptGoogle, verifycode) {
			return nil, enum.VerifyInCorrectCode, nil
		}
	} else if userdata.LoginGoogle != "" {
		if global.IsEnvPrd() && !xutils.VerifyGoogleCode(userdata.LoginGoogle, verifycode) {
			return nil, enum.VerifyInCorrectCode, nil
		}
	}
	secret, url := xutils.NewGoogleSecret("直播登录", reqdata.Account)
	db = server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	db = db.Updates(map[string]interface{}{
		edb.LoginGoogle: secret,
	})
	if db.Error != nil {
		logs.Error("SetLoginGoogle error", db.Error)
		return nil, nil, db.Error
	}
	googlesecret := &SetLoginGoogleRes{Url: url}
	return googlesecret, nil, db.Error
}

type SetOptGoogleReq struct {
	Account string `validate:"required" json:"account"` // 账号
}
type SetOptGoogleRes struct {
	Url string `json:"url"` // 二维码
}

// 设置操作验证码
func (this *ServiceAdmin) SetOptGoogle(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(SetOptGoogleReq)
	token := server.GetToken(ctx)
	locker := fmt.Sprintf(enum.Lock_ChangeGoogleSecret, token.SellerId, reqdata.Account)
	if !server.Redis().Lock(locker, 5) {
		return nil, enum.TooManyRequest, nil
	}

	verifycode := ctx.Request.Header.Get("VerifyCode")
	tokendata := server.GetToken(ctx)

	userdata := model.XAdminUser{}
	db := server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, tokendata.SellerId)
	db = db.Where(edb.Account+edb.EQ, tokendata.Account)
	db = db.First(&userdata)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, enum.UserNotFound, nil
		}
		return nil, nil, db.Error
	}
	if userdata.LoginGoogle == "" && userdata.LoginGoogle == reqdata.Account {
		return nil, enum.VerifyNotFoundSecret, nil
	}

	if userdata.OptGoogle != "" {
		if global.IsEnvPrd() && !xutils.VerifyGoogleCode(userdata.OptGoogle, verifycode) {
			return nil, enum.VerifyInCorrectCode, nil
		}
	} else if userdata.LoginGoogle != "" {
		if global.IsEnvPrd() && !xutils.VerifyGoogleCode(userdata.LoginGoogle, verifycode) {
			return nil, enum.VerifyInCorrectCode, nil
		}
	}
	secret, url := xutils.NewGoogleSecret("直播操作", reqdata.Account)
	db = server.Db().Model(&model.XAdminUser{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.Account+edb.EQ, reqdata.Account)
	db = db.Updates(map[string]interface{}{
		"opt_google": secret,
	})
	if db.Error != nil {
		logs.Error("SetOptGoogle error", db.Error)
		return nil, nil, db.Error
	}
	googlesecret := &SetOptGoogleRes{Url: url}
	return googlesecret, nil, db.Error
}
