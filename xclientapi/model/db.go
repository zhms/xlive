package model

import "github.com/shopspring/decimal"

type XUser struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`                              // 用户Id
	SellerId   int    `gorm:"column:seller_id;comment:运营商"`                                    // 运营商Id
	ChannelId  int    `gorm:"column:channel_id;comment:渠道"`                                    // 渠道Id
	UserId     int    `gorm:"column:user_id;comment:玩家Id"`                                     // 玩家Id
	State      int    `gorm:"column:state;default:1;comment:状态 1启用,2禁用"`                       // 用户状态，1表示启用，2表示禁用
	Account    string `gorm:"column:account;type:varchar(32);comment:账号"`                      // 用户账号
	Password   string `gorm:"column:password;type:varchar(64);comment:密码"`                     // 用户密码
	NickName   string `gorm:"column:nick_name;type:varchar(32);comment:昵称"`                    // 用户昵称
	PhoneNum   string `gorm:"column:phone_num;type:varchar(32);comment:电话号码"`                  // 用户电话号码
	Email      string `gorm:"column:email;type:varchar(64);comment:Email地址"`                   // 用户Email地址
	TopAgent   int    `gorm:"column:top_agent;comment:顶级代理"`                                   // 顶级代理Id
	Agent      int    `gorm:"column:agent;comment:上级代理"`                                       // 上级代理Id
	Agents     string `gorm:"column:agents;type:text;comment:代理"`                              // 代理列表
	Token      string `gorm:"column:token;type:varchar(64);comment:最后登录的token"`                // 最后登录的token
	RegIp      string `gorm:"column:reg_ip;type:varchar(32);comment:注册Ip"`                     // 注册Ip
	LoginIp    string `gorm:"column:login_ip;type:varchar(32);comment:最后登录Ip"`                 // 最后登录Ip
	LoginTime  string `gorm:"column:login_time;comment:最后登录时间"`                                // 最后登录时间
	CreateTime string `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"` // 创建时间
}

func (this *XUser) TableName() string {
	return "x_user"
}

type XUserScore struct {
	Id           uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                  // 自增Id
	SellerId     int             `gorm:"column:seller_id" json:"seller_id"`                                             // 运营商
	ChannelId    int             `gorm:"column:channel_id" json:"channel_id"`                                           // 渠道
	UserId       int             `gorm:"column:user_id" json:"user_id"`                                                 // 玩家Id
	Symbol       string          `gorm:"column:symbol" json:"symbol"`                                                   // 币种
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(50,6);default:0.000000" json:"amount"`               // 可用金额
	FrozenAmount decimal.Decimal `gorm:"column:frozen_amount;type:decimal(50,6);default:0.000000" json:"frozen_amount"` // 冻结金额
	BankAmount   decimal.Decimal `gorm:"column:bank_amount;type:decimal(50,6);default:0.000000" json:"bank_amount"`     // 保险箱金额
	CreateTime   string          `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`               // 创建时间
}

func (this *XUserScore) TableName() string {
	return "x_user_score"
}

type XUserPool struct {
	UserId     int    `gorm:"primary_key;comment:'玩家Id'" json:"user_id"`                       // 玩家Id
	CreateTime string `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"` // 创建时间
}

func (this *XUserPool) TableName() string {
	return "x_user_pool"
}
