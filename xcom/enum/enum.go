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
	Lock_ChangeGoogleSecret string = "change_google_secret_%d_%s"
	Lock_UserLogin          string = "user_login:"
	Lock_UserRegister       string = "user_register:"
)

const (
	StateYes = 1 // 是
	StateNo  = 2 // 否
)
