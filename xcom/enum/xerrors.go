package enum

import "github.com/gin-gonic/gin"

func MakeError(data map[string]interface{}, errmsg string) *map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range data {
		m[k] = v
	}
	m["data"] = errmsg
	return &m
}

func MakeSucess(value interface{}) *map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range Success {
		m[k] = v
	}
	m["data"] = value
	return &m
}

func MakePageSucess(total int64, value interface{}) *map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range Success {
		m[k] = v
	}
	m["data"] = gin.H{"total": total, "data": value}
	return &m
}

// 通用
var Success = map[string]interface{}{"code": 0, "msg": "成功"}
var BadParams = map[string]interface{}{"code": 1, "msg": "参数错误"}
var InternalError = map[string]interface{}{"code": 2, "msg": "内部错误"}
var TooManyRequest = map[string]interface{}{"code": 3, "msg": "请求频繁"}
var GetUserInfoError = map[string]interface{}{"code": 4, "msg": "获取用户信息失败"}
var SellerNotFound = map[string]interface{}{"code": 5, "msg": "商户不存在"}
var LiveNotAvailable = map[string]interface{}{"code": 6, "msg": "直播不可用"}
var AlreadyAudited = map[string]interface{}{"code": 7, "msg": "已审核"}

// 验证码
var VerifyNotFoundCode = map[string]interface{}{"code": 100001, "msg": "未填写验证码"}
var VerifyNotFoundSecret = map[string]interface{}{"code": 100002, "msg": "未绑定验证码"}
var VerifyInCorrectCode = map[string]interface{}{"code": 100003, "msg": "验证码不正确"}

// 鉴权
var AuthTokenNotFound = map[string]interface{}{"code": 100101, "msg": "未填写token"}
var AuthGetTokenError = map[string]interface{}{"code": 100102, "msg": "获取token失败"}
var AuthTokenExpired = map[string]interface{}{"code": 100103, "msg": "未登录或登录已过期"}
var AuthNotFoundMainMenu = map[string]interface{}{"code": 100104, "msg": "权限不足"}
var AuthNotFoundSubMenu = map[string]interface{}{"code": 100105, "msg": "权限不足"}
var AuthNotFoundOpt = map[string]interface{}{"code": 100106, "msg": "权限不足"}
var AuthNotAllow = map[string]interface{}{"code": 100107, "msg": "权限不足"}
var ParentRoleNotFound = map[string]interface{}{"code": 100108, "msg": "父级角色不存在"}
var RoleNotEditable = map[string]interface{}{"code": 100109, "msg": "角色不允许编辑"}
var RoleCantDelete = map[string]interface{}{"code": 100110, "msg": "角色不允许删除"}
var RoleNotFound = map[string]interface{}{"code": 100111, "msg": "角色不存在"}
var RoleStateError = map[string]interface{}{"code": 100112, "msg": "角色状态异常"}

// 用户
var UserNotFound = map[string]interface{}{"code": 100201, "msg": "用户不存在"}
var UserPasswordError = map[string]interface{}{"code": 100202, "msg": "密码错误"}
var UserStateError = map[string]interface{}{"code": 100203, "msg": "用户状态异常"}
var UserExist = map[string]interface{}{"code": 100204, "msg": "账号已存在"}
var UserCantDelete = map[string]interface{}{"code": 100205, "msg": "账号不允许删除"}
var NewIdError = map[string]interface{}{"code": 100206, "msg": "新建用户Id错误"}
