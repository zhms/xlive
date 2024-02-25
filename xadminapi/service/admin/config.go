package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"

	"gorm.io/gorm"
)

type GetXConfigReq struct {
	SellerId   int      `validate:"required" from:"seller_id"` // 运营商
	ChannelId  int      `from:"channel_id"`                    // 渠道商
	ConfigName []string `from:"config_name"`                   // 配置名称
}

type GetXConfigRes struct {
	Data []model.XConfig `json:"data"`
}

// 获取配置
func (this *ServiceAdmin) GetXConfig(reqdata *GetXConfigReq) (data *GetXConfigRes, merr map[string]interface{}, err error) {
	data = &GetXConfigRes{}
	db := server.Db().Model(&model.XConfig{})
	db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
	db = db.Where(edb.ChannelId+edb.EQ, reqdata.ChannelId)
	db = db.Where(edb.ConfigName+edb.IN, reqdata.ConfigName)
	db = db.Find(&data.Data)
	if db.Error != nil {
		return nil, nil, err
	}
	return data, nil, err
}

type UpdateXConfigData struct {
	ChannelId   int    `json:"channel_id"`                      // 渠道商
	ConfigName  string `validate:"required" json:"config_name"` // 配置名称
	ConfigValue string `json:"config_value"`                    // 配置数据
	Memo        string `json:"memo"`                            // 备注
}

type UpdateXConfigReq struct {
	SellerId int                 `validate:"required" json:"seller_id"` // 运营商
	Configs  []UpdateXConfigData `json:"configs"`                       // 配置数据
}

// 更新配置
func (this *ServiceAdmin) UpdateXConfig(reqdata *UpdateXConfigReq) (merr map[string]interface{}, err error) {
	err = server.Db().Transaction(func(db *gorm.DB) error {
		db = db.Model(&model.XConfig{})
		for _, v := range reqdata.Configs {
			db = db.Where(edb.SellerId+edb.EQ, reqdata.SellerId)
			db = db.Where(edb.ChannelId+edb.EQ, v.ChannelId)
			db = db.Where(edb.ConfigName+edb.EQ, v.ConfigName)
			db = db.Updates(map[string]interface{}{
				edb.ConfigValue: v.ConfigValue,
				edb.Memo:        v.Memo,
			})
			if db.Error != nil {
				return err
			}
		}
		return nil
	})
	return nil, err
}
