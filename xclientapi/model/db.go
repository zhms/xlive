package model

type XUser struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`                              // 用户Id
	SellerId   int    `gorm:"column:seller_id;comment:运营商"`                                    // 运营商Id
	UserId     int    `gorm:"column:user_id;comment:玩家Id"`                                     // 玩家Id
	State      int    `gorm:"column:state;default:1;comment:状态 1启用,2禁用"`                       // 用户状态，1启用，2禁用
	Account    string `gorm:"column:account;type:varchar(32);comment:账号"`                      // 用户账号
	Password   string `gorm:"column:password;type:varchar(64);comment:密码"`                     // 用户密码
	IsVisitor  int    `gorm:"column:is_visitor;default:0;comment:是否游客 1是,2否"`                  // 是否游客
	Token      string `gorm:"-"`                                                               // 最后登录的token
	CreateTime string `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"` // 创建时间
}

func (this *XUser) TableName() string {
	return "x_user"
}
