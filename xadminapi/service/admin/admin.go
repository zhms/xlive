package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"

	"xcom/edb"
	"xcom/utils"

	"github.com/beego/beego/logs"
	"gorm.io/gorm"
)

type ServiceAdmin struct {
}

// 初始化权限
func (this *ServiceAdmin) Init() {
	result, err := server.Db().Raw("call x_init_auth()").Rows()
	if err != nil {
		logs.Error("x_init_auth error:", err.Error())
		return
	}
	if result != nil {
		data := utils.DbResult(result)
		if data.Length() > 0 {
			logs.Warn(data.Index(0).String("Warning"))
		}
	}
}

func (this *ServiceAdmin) role_exists(sellerid int, parent string) (exists bool, err error) {
	parentmodel := &model.XAdminRole{}
	db := server.Db().Model(parentmodel)
	db = db.Where(edb.SellerId+edb.EQ, sellerid)
	db = db.Where(edb.RoleName+edb.EQ, parent)
	db = db.First(&parentmodel)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, db.Error
	}
	return true, err
}
