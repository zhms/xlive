package edb

const (
	TableSeller     = "x_seller"       // 运营商表
	TableUser       = "x_user"         // 用户表
	TableRole       = "x_role"         // 角色表
	TableConfig     = "x_config"       // 配置表
	TableUserIdPool = "x_user_id_pool" // 用户池
	TableHostSeller = "x_host_seller"  // 域名对应运营商
)

const (
	StateYes = 1 // 启用
	StateNo  = 2 // 禁用
)

const (
	DESC      = " desc "
	ASC       = " asc "
	EQ        = " = ? "
	NEQ       = " <> ? "
	AND       = " and "
	OR        = " or "
	ISNULL    = " is null "
	ISNOTNULL = " is not null "
	LIKE      = " like ? "
	IN        = " in (?) "
	NOTIN     = " not in (?) "
	GT        = " > ? "
	GTE       = " >= ? "
	LT        = " < ? "
	LTE       = " <= ? "
	BETWEEN   = " between ? and ? "
	PLUS      = " + ? "
	MINUS     = " - ? "
)

const (
	SellerId  = "seller_id"  // 运营商
	Account   = "account"    // 账号
	Password  = "password"   // 密码
	IsVisitor = "is_visitor" // 是否游客
	UserId    = "user_id"    // 用户Id
	State     = "state"      // 状态
	Token     = "token"      // token
	LiveUrl   = "live_url"   // 直播地址
	Title     = "title"      // 标题
	Name      = "name"       // 名称
)
