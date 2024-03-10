package model

type XSeller struct {
	Id         uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT;comment:'自增Id'" json:"id"`          // 自增Id
	SellerId   int    `gorm:"column:seller_id;type:int;comment:'运营商'" json:"seller_id"`               // 运营商
	State      int    `gorm:"column:state;type:int;default:1;comment:'状态 1开启,2关闭'" json:"state"`      // 状态 1开启 2关闭
	SellerName string `gorm:"column:seller_name;type:varchar(32);comment:'运营商名称'" json:"seller_name"` // 运营商名称
	Memo       string `gorm:"column:memo;type:varchar(256);comment:'备注'" json:"memo"`                 // 备注
}

func (XSeller) TableName() string {
	return "x_seller"
}

type XAdminLoginLog struct {
	Id         uint64 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增Id'" json:"id"`         // 自增Id
	SellerId   int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                     // 运营商
	Account    string `gorm:"column:account;type:varchar(32);comment:'账号'" json:"account"`         // 账号
	Token      string `gorm:"column:token;type:varchar(64);comment:'登录的token'" json:"-"`           // 登录的token
	LoginIp    string `gorm:"column:login_ip;type:varchar(32);comment:'最近一次登录Ip'" json:"login_ip"` // 最近一次登录Ip
	Memo       string `gorm:"column:memo;type:varchar(256);comment:'备注'" json:"memo"`              // 备注
	IpLocation string `gorm:"-" json:"ip_location"`                                                //ip地理位置
	CreateTime string `gorm:"column:create_time" json:"create_time"`                               // 创建时间
}

func (XAdminLoginLog) TableName() string {
	return "x_admin_login_log"
}

type XConfig struct {
	Id          uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`          // 配置Id
	SellerId    int    `gorm:"column:seller_id" json:"seller_id" comment:"运营商"`       // 运营商Id
	ConfigName  string `gorm:"column:config_name" json:"config_name" comment:"配置名"`   // 配置名称
	ConfigValue string `gorm:"column:config_value" json:"config_value" comment:"配置值"` // 配置值
	Memo        string `gorm:"column:memo" json:"memo" comment:"备注"`                  // 备注
	CreateTime  string `gorm:"column:create_time" json:"create_time"`                 // 创建时间
}

func (XConfig) TableName() string {
	return "x_config"
}

type XAdminOptLog struct {
	Id         uint64 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增Id'" json:"id"`                      // 自增Id
	SellerId   int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                  // 运营商
	Account    string `gorm:"column:account;type:varchar(32);charset:utf8mb4;comment:'账号'" json:"account"`      // 账号
	ReqPath    string `gorm:"column:req_path;type:varchar(256);charset:utf8mb4;comment:'请求路径'" json:"req_path"` // 请求路径
	OptName    string `gorm:"column:opt_name;type:varchar(64);charset:utf8mb4;comment:'操作名称'" json:"opt_name"`  // 请求路径
	ReqData    string `gorm:"column:req_data;type:varchar(256);charset:utf8mb4;comment:'请求参数'" json:"req_data"` // 请求参数
	ReqIp      string `gorm:"column:req_ip;type:varchar(32);charset:utf8mb4;comment:'请求的Ip'" json:"req_ip"`     // 请求的Ip
	IpLocation string `gorm:"-" json:"ip_location"`                                                             // Ip地理位置
	CreateTime string `gorm:"column:create_time" json:"create_time"`                                            // 创建时间
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

type XUser struct {
	Id            int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                   // 自增Id
	SellerId      int    `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                // 运营商
	UserId        int    `gorm:"column:user_id;comment:'用户id'" json:"user_id"`                                   // 用户id
	Account       string `gorm:"column:account;type:varchar(32);comment:'用户账号'" json:"account"`                  // 用户账号
	Password      string `gorm:"column:password;type:varchar(32);comment:'用户密码'" json:"password"`                // 用户密码
	IsVisitor     int    `gorm:"column:is_visitor;comment:'是否是游客'" json:"is_visitor"`                            // 是否是游客
	State         int    `gorm:"column:state;comment:'状态 1开启,2关闭'" json:"state"`                                 // 状态 1开启,2关闭
	Agent         string `gorm:"column:agent;type:varchar(32);comment:'所属管理员'" json:"agent"`                     // 所属管理员
	LoginIP       string `gorm:"column:login_ip;type:varchar(64);comment:'登录ip'" json:"login_ip"`                // 登录ip
	LoginLocation string `gorm:"column:login_location;type:varchar(64);comment:'登录ip地区'" json:"login_location"`  // 登录ip地区
	LoginCount    int    `gorm:"column:login_count;comment:'登录次数'" json:"login_count"`                           // 登录次数
	ChatState     int    `gorm:"column:chat_state;comment:'聊天状态 1开启,2关闭'" json:"chat_state"`                     // 聊天状态 1开启,2关闭
	LoginTime     string `gorm:"column:login_time;comment:'登录时间'" json:"login_time"`                             // 登录时间
	CreateTime    string `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"` // 创建时间
}

func (XUser) TableName() string {
	return "x_user"
}

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

type XChatList struct {
	Id         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`            // 自增Id
	UserId     int    `gorm:"column:user_id" json:"user_id" comment:"会员id"`            // 会员id
	Account    string `gorm:"column:account" json:"account" comment:"会员账号"`            // 会员账号
	RoomId     int    `gorm:"column:room_id" json:"room_id" comment:"房价id"`            // 房价id
	Content    string `gorm:"column:content" json:"content" comment:"内容"`              // 内容
	Ip         string `gorm:"column:ip" json:"ip" comment:"ip"`                        // ip
	IpLocation string `gorm:"column:ip_location" json:"ip_location"`                   // ip_location
	State      int    `gorm:"column:state" json:"state" comment:"状态 1待审核,2审核通过,3审核拒绝"` // 状态 1待审核,2审核通过,3审核拒绝
	CreateTime string `gorm:"column:create_time" json:"create_time"`                   // 创建时间
}

func (XChatList) TableName() string {
	return "x_chat_list"
}

type XChatBanIP struct {
	Id           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Ip           string `gorm:"column:ip;type:varchar(32);charset:utf8mb4;not null" json:"ip"`
	AdminAccount string `gorm:"column:admin_account;type:varchar(32);charset:utf8mb4" json:"admin_account"`
	CreateTime   string `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
}

func (XChatBanIP) TableName() string {
	return "x_chat_ban_ip"
}

type XHongbao struct {
	Id          int     `gorm:"column:id;primary_key;AUTO_INCREMENT;comment:'id'" json:"id"`                      // id
	SellerId    int     `gorm:"column:seller_id;comment:'运营商'" json:"seller_id"`                                  // 运营商
	RoomId      int     `gorm:"column:room_id;comment:'房间Id'" json:"room_id"`                                     // 房间Id
	TotalAmount float64 `gorm:"column:total_amount;type:decimal(50,6);comment:'红包总金额'" json:"total_amount"`       // 红包总金额
	TotalCount  int     `gorm:"column:total_count;comment:'红包个数'" json:"total_count"`                             // 红包个数
	UsedAmount  float64 `gorm:"column:used_amount;type:decimal(50,6);comment:'已领取金额'" json:"used_amount"`         // 已领取金额
	UsedCount   int     `gorm:"column:used_count;comment:'已领取个数'" json:"used_count"`                              // 已领取个数
	Sender      string  `gorm:"column:sender;type:varchar(32);charset:utf8mb4;comment:'发送者'" json:"sender"`       // 发送者
	Memo        string  `gorm:"column:memo;type:varchar(255);charset:utf8mb4;comment:'备注'" json:"memo"`           // 备注
	CreateTime  string  `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:'红包发送时间'" json:"create_time"` // 红包发送时间
}

func (XHongbao) TableName() string {
	return "x_hongbao"
}
