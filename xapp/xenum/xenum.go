package xenum

import "github.com/gin-gonic/gin"

func MakeError(data map[string]interface{}, errmsg string) *map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range data {
		m[k] = v
	}
	m["data"] = errmsg
	return &m
}

func MakeSucess(value interface{}) *map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range Success {
		m[k] = v
	}
	m["data"] = value
	return &m
}

func MakePageSucess(total int64, value interface{}) *map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range Success {
		m[k] = v
	}
	m["data"] = gin.H{"total": total, "data": value}
	return &m
}

var Success = map[string]interface{}{"code": 0, "msg": "成功"}
var BadParams = map[string]interface{}{"code": 1, "msg": "参数错误"}
var TooManyRequest = map[string]interface{}{"code": 2, "msg": "请求频繁"}
var InternalError = map[string]interface{}{"code": 3, "msg": "内部错误"}
var LiveNotAvailable = map[string]interface{}{"code": 4, "msg": "直播不可用"}

var VerifyCodeNotFound = map[string]interface{}{"code": 100101, "msg": "请填写验证码"}
var GetTokenError = map[string]interface{}{"code": 100102, "msg": "获取token失败"}
var TokenExpired = map[string]interface{}{"code": 100103, "msg": "未登录或登录已过期"}
var VerifySecretNotFound = map[string]interface{}{"code": 100104, "msg": "账号未绑定验证秘钥"}
var VerifyCodeError = map[string]interface{}{"code": 100105, "msg": "验证码不正确"}
var Unauthorized1 = map[string]interface{}{"code": 100106, "msg": "权限不足1"}
var Unauthorized2 = map[string]interface{}{"code": 100106, "msg": "权限不足2"}
var Unauthorized3 = map[string]interface{}{"code": 100106, "msg": "权限不足3"}
var Unauthorized4 = map[string]interface{}{"code": 100106, "msg": "权限不足4"}
var Unauthorized5 = map[string]interface{}{"code": 100106, "msg": "权限不足5"}
var Unauthorized6 = map[string]interface{}{"code": 100106, "msg": "权限不足6"}
var Unauthorized7 = map[string]interface{}{"code": 100106, "msg": "权限不足7"}
var TokenNotFound = map[string]interface{}{"code": 100107, "msg": "请填写token"}

var UserNotFound = map[string]interface{}{"code": 100201, "msg": "用户不存在"}
var UserPasswordError = map[string]interface{}{"code": 100202, "msg": "密码错误"}
var UserBanded = map[string]interface{}{"code": 100203, "msg": "用户已被禁用"}
var UserCantLogin = map[string]interface{}{"code": 100204, "msg": "用户无法登录此系统"}
var UserStateError = map[string]interface{}{"code": 100205, "msg": "用户状态异常"}
var UserExist = map[string]interface{}{"code": 100206, "msg": "账号已存在"}
var UserCantDelete = map[string]interface{}{"code": 100207, "msg": "账号不允许删除"}

var RoleNotFound = map[string]interface{}{"code": 100301, "msg": "角色不存在"}
var RoleBaned = map[string]interface{}{"code": 100302, "msg": "角色已被禁用"}

var AuthTokenNotFound = map[string]interface{}{"code": 100401, "msg": "未填写token"}
var AuthGetTokenError = map[string]interface{}{"code": 100402, "msg": "获取token失败"}
var AuthTokenExpired = map[string]interface{}{"code": 100403, "msg": "未登录或登录已过期"}
