package enum

const (
	RedisTokenRw     = 1
	RedisUserRw      = 2
	RedisGameRw      = 3
	RedisOrderRw     = 4
	RedisStatisticRw = 5
	RedisAdminRw     = 6
)

const (
	DbGameRW      = 0
	DbUserRW      = 1
	DbOrderRW     = 2
	DbLogRW       = 3
	DbDaillyRW    = 4
	DbStatisticRW = 5
	DbAdminRW     = 6
)

const (
	Lock_AdminLogin         string = "admin_login:"
	Lock_ChangeGoogleSecret string = "change_google_secret_%d_%d_%s"
	Lock_UserLogin          string = "user_login:"
	Lock_UserRegister       string = "user_register:"
)

const (
	StateYes = 1 // 是
	StateNo  = 2 // 否
)

const (
	EMail_UnRead = 1 // 未读
	EMil_Readed  = 2 // 已读
)

const (
	UserBindAddressType_Guest = 1 // 游客
	UserBindAddressType_User  = 2 // 用户
)

const (
	Game_HashDaXiao    = 1 // 哈希大小
	Game_HashDanShuang = 2 // 哈希单双

	Game_YiFenHashDaXiao    = 301 // 一分哈希大小
	Game_YiFenHashDanShuang = 302 // 一分哈希单双

	Game_HashHeZhiDaXiao    = 401 // 哈希和值大小
	Game_HashHeZhiDanShuang = 402 // 哈希和值单双

	Game_YueDaXiao    = 101 // 余额大小
	Game_YueDanShuang = 102 // 余额单双

	Game_YiFenYueDaXiao    = 201 // 一分余额大小
	Game_YiFenYueDanShuang = 202 // 一分余额单双

	Game_LotteryK31 = 501 // 一分快3
)

const (
	//返奖类型
	OrderRewardType_Normal = 1 // 普通
	//订单审核状态
	OrderAuditState_None   = 1 // 待审核
	OrderAuditState_Pass   = 2 // 审核通过
	OrderAuditState_Refuse = 3 // 审核拒绝
	//订单审核类型
	OrderAuditType_None = 1 // 无需审核
	//订单状态
	OrderState_WaitingOpen = 1   // 待开奖
	OrderState_AmountLess  = 3   // 低于限额
	OrderState_AmountMore  = 4   // 高于限额
	OrderState_InvalidBet  = 5   // 无效投注
	OrderState_Success     = 200 // 处理完成
	//订单输赢状态
	OrderWinState_Win  = 1 // 赢
	OrderWinState_Lose = 2 // 输
)

const (
	AmountChangeType_Unknown    = 0 // 未知
	AmountChangeType_Recharge   = 1 // 充值
	AmountChangeType_Withward   = 2 // 提现
	AmountChangeType_AdminAlter = 3 // 管理员修改

	AmountChangeType_HashBet    = 4 // 哈希下注
	AmountChangeType_HashReturn = 5 // 哈希退回
	AmountChangeType_HashReward = 6 // 哈希返奖

	AmountChangeType_LotteryBet    = 7 // 彩票下注
	AmountChangeType_LotteryReturn = 8 // 彩票退回
	AmountChangeType_LotteryReward = 9 // 彩票返奖

)
