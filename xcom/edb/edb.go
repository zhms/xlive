package edb

const (
	TableSeller      = "x_seller"        // 运营商表
	TableUser        = "x_user"          // 用户表
	TableRole        = "x_role"          // 角色表
	TableConfig      = "x_config"        // 配置表
	TableUserIdPool  = "x_user_id_pool"  // 用户池
	TableHostSeller  = "x_host_seller"   // 域名对应运营商
	TableKv          = "x_kv"            // 键值对
	TableLiveRoom    = "x_live_room"     // 直播间
	TableChatList    = "x_chat_list"     // 聊天列表
	TableChatBanIp   = "x_chat_ban_ip"   // IP禁言
	TableChatBanUser = "x_chat_ban_user" // 用户禁言
	TableHongbao     = "x_hongbao"       // 红包
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
	SellerId      = "seller_id"      // 运营商
	Account       = "account"        // 账号
	Password      = "password"       // 密码
	IsVisitor     = "is_visitor"     // 是否游客
	UserId        = "user_id"        // 用户Id
	State         = "state"          // 状态
	Token         = "token"          // token
	LiveUrl       = "live_url"       // 直播地址
	Title         = "title"          // 标题
	Name          = "name"           // 名称
	K             = "k"              // 键
	V             = "v"              // 值
	RoleName      = "role_name"      // 角色名称
	LoginIp       = "login_ip"       // 登录IP
	CreateTime    = "create_time"    // 创建时间
	Memo          = "memo"           // 备注
	ConfigName    = "config_name"    // 配置名称
	ConfigValue   = "config_value"   // 配置值
	ChannelId     = "channel_id"     // 渠道Id
	SellerName    = "seller_name"    // 运营商名称
	LoginGoogle   = "login_google"   // 谷歌登录
	Id            = "id"             // Id
	OptName       = "opt_name"       // 操作名称
	Parent        = "parent"         // 父级
	RoleData      = "role_data"      // 角色数据
	LoginTime     = "login_time"     // 登录时间
	LoginCount    = "login_count"    // 登录次数
	Agent         = "agent"          // 代理商
	ChatState     = "chat_state"     // 聊天状态
	LoginLocation = "login_location" // IP位置
	RoomId        = "room_id"        // 房间Id
	Content       = "content"        // 内容
	Ip            = "ip"             // IP
	IpLocation    = "ip_location"    // IP位置
	AdminAccount  = "admin_account"  // 管理员账号
	TotalAmount   = "total_amount"   // 总金额
	TotalCount    = "total_count"    // 总数量
	Sender        = "sender"         // 发送者
)
