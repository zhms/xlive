package service_admin

import (
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/utils"

	"github.com/gin-gonic/gin"
)

type GetAdminRoleReq struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	RoleName string `json:"role_name"` // 角色名
}

type GetAdminRoleRes struct {
	Total int64              `json:"total"` // 总数
	Data  []model.XAdminRole `json:"data"`  // 数据
}

// 获取角色列表
func (this *ServiceAdmin) GetRoleList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetAdminRoleReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetAdminRoleRes{}
	db := server.Db().Model(&model.XAdminRole{})
	db = utils.DbWhere(db, edb.SellerId, token.SellerId, int(0))
	db = utils.DbWhere(db, edb.RoleName, reqdata.RoleName, "")
	err = db.Count(&data.Total).Error
	if err != nil {
		return nil, nil, err
	}
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return nil, nil, err
	}
	return data, nil, err
}

type CreateAdminRoleReq struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
	Parent   string `validate:"required" json:"parent"`    // 上级角色
	RoleData string `validate:"required" json:"role_data"` // 权限数据
	State    int    `validate:"required" json:"state"`     // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// 创建角色
func (this *ServiceAdmin) CreateRole(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(CreateAdminRoleReq)
	token := server.GetToken(ctx)
	exists, err := this.role_exists(token.SellerId, reqdata.Parent)
	if err != nil {
		return nil, err
	}
	if !exists {
		return enum.ParentRoleNotFound, nil
	}
	db := server.Db().Model(&model.XAdminRole{})
	db = db.Create(map[string]interface{}{
		edb.SellerId: token.SellerId,
		edb.RoleName: reqdata.RoleName,
		edb.Parent:   reqdata.Parent,
		edb.RoleData: reqdata.RoleData,
		edb.State:    reqdata.State,
		edb.Memo:     reqdata.Memo,
	})
	return nil, db.Error
}

type UpdateAdminRoleReq struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
	Parent   string `json:"parent"`                        // 上级角色
	RoleData string `json:"role_data"`                     // 权限数据
	State    int    `json:"state"`                         // 状态 1开启,2关闭
	Memo     string `json:"memo"`                          // 备注
}

// 更新角色
func (this *ServiceAdmin) UpdateRole(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateAdminRoleReq)
	if reqdata.RoleName == "超级管理员" || reqdata.RoleName == "运营商超管" {
		return enum.RoleNotEditable, nil
	}
	token := server.GetToken(ctx)
	updatedata := map[string]interface{}{}
	utils.MapSet(&updatedata, edb.Memo, reqdata.Memo, "")
	if reqdata.Parent != "" {
		exists, err := this.role_exists(token.SellerId, reqdata.Parent)
		if err != nil {
			return nil, err
		}
		if !exists {
			return enum.ParentRoleNotFound, nil
		}
		utils.MapSet(&updatedata, edb.Parent, reqdata.Parent, "")
	}
	utils.MapSet(&updatedata, edb.RoleData, reqdata.RoleData, "")
	utils.MapSetIn(&updatedata, edb.State, reqdata.State, []interface{}{int(1), int(2)})
	db := server.Db().Model(&model.XAdminRole{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.RoleName+edb.EQ, reqdata.RoleName)
	db = db.Updates(updatedata)
	return nil, db.Error
}

type DeleteAdminRoleReq struct {
	RoleName string `validate:"required" json:"role_name"` // 角色
}

// 删除角色
func (this *ServiceAdmin) DeleteRole(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	token := server.GetToken(ctx)
	reqdata := idata.(DeleteAdminRoleReq)
	if reqdata.RoleName == "超级管理员" || reqdata.RoleName == "运营商超管" {
		return enum.RoleCantDelete, nil
	}
	db := server.Db().Model(&model.XAdminRole{})
	db = db.Where(edb.SellerId+edb.EQ, token.SellerId)
	db = db.Where(edb.RoleName+edb.EQ, reqdata.RoleName)
	db = db.Delete(&model.XAdminRole{})
	return nil, db.Error
}
