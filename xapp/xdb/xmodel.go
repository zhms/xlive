package xdb

type XLiveRoom struct {
	Id         int    `gorm:"column:id;primaryKey;autoIncrement;comment:'id'" json:"id"`                         // id
	SellerId   int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                   // 运营商
	Name       string `gorm:"column:name;type:varchar(32);comment:'直播间名称'" json:"name"`                          // 直播间名称
	Account    string `gorm:"column:account;type:varchar(32);charset:utf8mb4;comment:'主播账号'" json:"account"`     // 主播账号
	PushURL    string `gorm:"column:push_url;type:varchar(1024);charset:utf8mb4;comment:'推流地址'" json:"push_url"` // 推流地址
	PullURL    string `gorm:"column:pull_url;type:varchar(1024);charset:utf8mb4;comment:'拉流地址'" json:"pull_url"` // 拉流地址
	LiveURL    string `gorm:"column:live_url;type:varchar(1024);charset:utf8mb4;comment:'前端地址'" json:"live_url"` // 前端地址
	State      int    `gorm:"column:state;comment:'状态1正在直播,2未直播'" json:"state"`                                  //状态 1正在直播,2未直播
	Title      string `gorm:"column:title;type:varchar(32);charset:utf8mb4;comment:'直播间标题'" json:"title"`        // 直播间标题
	CreateTime string `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`    // 创建时间
}

func (XLiveRoom) TableName() string {
	return "x_live_room"
}

type XUser struct {
	Id            int    `gorm:"column:id;primaryKey;autoIncrement;comment:'id'" json:"id"`
	SellerId      int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`
	Account       string `gorm:"column:account;type:varchar(32);charset:utf8mb4;comment:'用户账号'" json:"account"`
	Password      string `gorm:"column:password;type:varchar(32);charset:utf8mb4;comment:'用户密码'" json:"password"`
	IsVisitor     int    `gorm:"column:is_visitor;comment:'是否是游客'" json:"is_visitor"`
	State         int    `gorm:"column:state;comment:'状态 1开启,2关闭'" json:"state"`
	Token         string `gorm:"column:token;type:varchar(32);charset:utf8mb4;comment:'token'" json:"token"`
	Agent         string `gorm:"column:agent;type:varchar(32);charset:utf8mb4;comment:'所属管理员'" json:"agent"`
	LoginIP       string `gorm:"column:login_ip;type:varchar(64);charset:utf8mb4;comment:'登录ip'" json:"login_ip"`
	LoginLocation string `gorm:"column:login_location;type:varchar(64);charset:utf8mb4;comment:'登录ip地区'" json:"login_location"`
	LoginCount    int    `gorm:"column:login_count;comment:'登录次数'" json:"login_count"`
	LoginTime     string `gorm:"column:login_time;comment:'登录时间'" json:"login_time"`
	ChatState     int    `gorm:"column:chat_state;comment:'禁言 1是,2否'" json:"chat_state"`
	CreateTime    string `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
}

func (XUser) TableName() string {
	return "x_user"
}
