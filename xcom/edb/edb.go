package edb

const (
	TableSeller          = "x_seller"            // 运营商表
	TableUser            = "x_user"              // 用户表
	TableRole            = "x_role"              // 角色表
	TableChannel         = "x_channel"           // 渠道商表
	TableConfig          = "x_config"            // 配置表
	TableUserIdPool      = "x_user_id_pool"      // 用户池
	TableUserScore       = "x_user_score"        // 用户金额表
	TableUserScoreLog    = "x_user_score_log"    // 用户金额日志表
	TableAgent           = "x_agent"             // 代理表
	TableAgentChild      = "x_agent_child"       // 代理关系表
	TableOrderHash       = "x_order_hash_"       // 哈希订单表
	TableUserBindAddress = "x_user_bind_address" // 绑定地址表
	TableLotteryHistory  = "x_lottery_history"   // 开奖历史表
	TableUserGameInfo    = "x_user_game_info"    // 用户游戏信息表
	TableSymbol          = "x_symbol"            // 货币符号表
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
	SellerId           = "seller_id"            // 运营商
	ChannelId          = "channel_id"           // 渠道商
	ChannelName        = "channel_name"         // 渠道商名称
	Account            = "account"              // 账号
	CreateTime         = "create_time"          // 创建时间
	Password           = "password"             // 登录密码
	RoleName           = "role_name"            // 角色名称
	LoginGoogle        = "login_google"         // 登录验证码
	OptGoogle          = "opt_google"           // 操作验证码
	State              = "state"                // 状态
	Token              = "token"                // 最后登录的token
	LoginCount         = "login_count"          // 登录次数
	LoginTime          = "login_time"           // 最后登录时间
	LoginIp            = "login_ip"             // 最后登录Ip
	Memo               = "memo"                 // 备注
	Id                 = "id"                   // 自增Id
	ConfigName         = "config_name"          // 配置名称
	ConfigValue        = "config_value"         // 配置数据
	OptName            = "opt_name"             // 操作名
	SellerName         = "seller_name"          // 运营商名称
	RoleData           = "role_data"            // 角色数据
	Parent             = "parent"               // 上级角色
	UserId             = "user_id"              // 用户Id
	NickName           = "nick_name"            // 昵称
	RegIp              = "reg_ip"               // 注册Ip
	Symbol             = "symbol"               // 货币符号
	Usdt               = "usdt"                 // 可用金额
	Trx                = "trx"                  // 可用金额
	UsdtFrozen         = "usdt_frozen"          // 冻结金额
	UsdtBank           = "usdt_bank"            // 保险箱金额
	Agent              = "agent"                // 上级代理
	Amount             = "amount"               // 可用金额
	FrozenAmount       = "frozen_amount"        // 冻结金额
	BankAmount         = "bank_amount"          // 保险箱金额
	BeforeAmount       = "before_amount"        // 变动前金额
	AfterAmount        = "after_amount"         // 变动后金额
	Reason             = "reason"               // 变动原因
	BeforeFrozenAmount = "before_frozen_amount" // 变动前冻结金额
	AfterFrozenAmount  = "after_frozen_amount"  // 变动后冻结金额
	BeforeBankAmount   = "before_bank_amount"   // 变动前保险箱金额
	AfterBankAmount    = "after_bank_amount"    // 变动后保险箱金额
	Agents             = "agents"               // 代理列表
	TopAgent           = "top_agent"            // 顶级代理
	ChildId            = "child_id"             // 下级代理
	ChildLevel         = "child_level"          // 下级代理等级
	ChildCount         = "child_count"          // 下级代理数量
	AuditTime          = "audit_time"           // 审核时间
	RewardTime         = "reward_time"          // 发放时间
	Title              = "title"                // 标题
	Sort               = "sort"                 // 排序
	Picture            = "picture"              // 图片
	Content            = "content"              // 内容
	ExLink             = "exlink"               // 链接
	MailId             = "mail_id"              // 邮件id
	StartTime          = "start_time"           // 开始时间
	EndTime            = "end_time"             // 结束时间
	Attachment         = "attachment"           // 附件
	ActiveId           = "active_id"            // 活动id
	IsTopAgent         = "is_top_agent"         // 是否顶级代理
	Address            = "address"              // 地址
	GameId             = "game_id"              // 游戏Id
	GameType           = "game_type"            // 游戏类型
	GameName           = "game_name"            // 游戏名称
	OpenTime           = "open_time"            // 开奖时间
	Period             = "period"               // 期号
	BlockHash          = "block_hash"           // 区块哈希
	OpenArea           = "open_area"            // 开奖区域
	OpenWord           = "open_word"            // 开奖字
	BetArea            = "bet_area"             // 投注区域
	IsWin              = "is_win"               // 是否中奖
	BlockNum           = "block_num"            // 区块高度
	RewardFee          = "reward_fee"           // 中奖金额
	RewardType         = "reward_type"          // 发放类型
	RewardRate         = "reward_rate"          // 发放比例
	FromAddress        = "from_address"         // 转出地址
	ToAddress          = "to_address"           // 转入地址
	TxId               = "tx_id"                // 交易id
	BetAmount          = "bet_amount"           // 投注金额
	NextBlockHash      = "next_block_hash"      // 下一区块哈希
	BlockMaker         = "block_maker"          // 区块生产者
	MakeTotal          = "make_total"           // 区块生产者总数
	BlockTime          = "block_time"           // 区块时间
	RoomLevel          = "room_level"           // 房间等级
	WinAmount          = "win_amount"           // 中奖金额
	ExtraData          = "extra_data"           // 额外数据
	LastGame           = "last_game"            // 最后一次游戏
	AuditState         = "audit_state"          // 审核状态
	BindType           = "bind_type"            // 绑定类型
	ActiveTime         = "active_time"          // 激活时间
	Data               = "data"                 // 数据
	BetTrx             = "bet_trx"              // 投注金额
	WinTrx             = "win_trx"              // 中奖金额
	BetUsdt            = "bet_usdt"             // 投注金额
	WinUsdt            = "win_usdt"             // 中奖金额
)

const (
	GameIgnore = "game:ignore" // 忽略的游戏
	GameInfo   = "game:info"   // 游戏信息
)
