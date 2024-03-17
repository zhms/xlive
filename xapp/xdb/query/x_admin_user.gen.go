// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"xapp/xdb/model"
)

func newXAdminUser(db *gorm.DB, opts ...gen.DOOption) xAdminUser {
	_xAdminUser := xAdminUser{}

	_xAdminUser.xAdminUserDo.UseDB(db, opts...)
	_xAdminUser.xAdminUserDo.UseModel(&model.XAdminUser{})

	tableName := _xAdminUser.xAdminUserDo.TableName()
	_xAdminUser.ALL = field.NewAsterisk(tableName)
	_xAdminUser.ID = field.NewInt64(tableName, "id")
	_xAdminUser.SellerID = field.NewInt32(tableName, "seller_id")
	_xAdminUser.Account = field.NewString(tableName, "account")
	_xAdminUser.Password = field.NewString(tableName, "password")
	_xAdminUser.RoleName = field.NewString(tableName, "role_name")
	_xAdminUser.LoginGoogle = field.NewString(tableName, "login_google")
	_xAdminUser.OptGoogle = field.NewString(tableName, "opt_google")
	_xAdminUser.Agent = field.NewString(tableName, "agent")
	_xAdminUser.State = field.NewInt32(tableName, "state")
	_xAdminUser.Token = field.NewString(tableName, "token")
	_xAdminUser.LoginCount = field.NewInt32(tableName, "login_count")
	_xAdminUser.LoginTime = field.NewTime(tableName, "login_time")
	_xAdminUser.LoginIP = field.NewString(tableName, "login_ip")
	_xAdminUser.Memo = field.NewString(tableName, "memo")
	_xAdminUser.CreateTime = field.NewTime(tableName, "create_time")
	_xAdminUser.RoomID = field.NewInt32(tableName, "room_id")

	_xAdminUser.fillFieldMap()

	return _xAdminUser
}

type xAdminUser struct {
	xAdminUserDo

	ALL         field.Asterisk
	ID          field.Int64  // 自增Id
	SellerID    field.Int32  // 运营商
	Account     field.String // 账号
	Password    field.String // 登录密码
	RoleName    field.String // 角色
	LoginGoogle field.String // 登录谷歌验证码
	OptGoogle   field.String // 渠道商
	Agent       field.String // 上级代理
	State       field.Int32  // 状态 1开启,2关闭
	Token       field.String // 最后登录的token
	LoginCount  field.Int32  // 登录次数
	LoginTime   field.Time   // 最后登录时间
	LoginIP     field.String // 最后登录Ip
	Memo        field.String // 备注
	CreateTime  field.Time   // 创建时间
	RoomID      field.Int32

	fieldMap map[string]field.Expr
}

func (x xAdminUser) Table(newTableName string) *xAdminUser {
	x.xAdminUserDo.UseTable(newTableName)
	return x.updateTableName(newTableName)
}

func (x xAdminUser) As(alias string) *xAdminUser {
	x.xAdminUserDo.DO = *(x.xAdminUserDo.As(alias).(*gen.DO))
	return x.updateTableName(alias)
}

func (x *xAdminUser) updateTableName(table string) *xAdminUser {
	x.ALL = field.NewAsterisk(table)
	x.ID = field.NewInt64(table, "id")
	x.SellerID = field.NewInt32(table, "seller_id")
	x.Account = field.NewString(table, "account")
	x.Password = field.NewString(table, "password")
	x.RoleName = field.NewString(table, "role_name")
	x.LoginGoogle = field.NewString(table, "login_google")
	x.OptGoogle = field.NewString(table, "opt_google")
	x.Agent = field.NewString(table, "agent")
	x.State = field.NewInt32(table, "state")
	x.Token = field.NewString(table, "token")
	x.LoginCount = field.NewInt32(table, "login_count")
	x.LoginTime = field.NewTime(table, "login_time")
	x.LoginIP = field.NewString(table, "login_ip")
	x.Memo = field.NewString(table, "memo")
	x.CreateTime = field.NewTime(table, "create_time")
	x.RoomID = field.NewInt32(table, "room_id")

	x.fillFieldMap()

	return x
}

func (x *xAdminUser) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := x.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (x *xAdminUser) fillFieldMap() {
	x.fieldMap = make(map[string]field.Expr, 16)
	x.fieldMap["id"] = x.ID
	x.fieldMap["seller_id"] = x.SellerID
	x.fieldMap["account"] = x.Account
	x.fieldMap["password"] = x.Password
	x.fieldMap["role_name"] = x.RoleName
	x.fieldMap["login_google"] = x.LoginGoogle
	x.fieldMap["opt_google"] = x.OptGoogle
	x.fieldMap["agent"] = x.Agent
	x.fieldMap["state"] = x.State
	x.fieldMap["token"] = x.Token
	x.fieldMap["login_count"] = x.LoginCount
	x.fieldMap["login_time"] = x.LoginTime
	x.fieldMap["login_ip"] = x.LoginIP
	x.fieldMap["memo"] = x.Memo
	x.fieldMap["create_time"] = x.CreateTime
	x.fieldMap["room_id"] = x.RoomID
}

func (x xAdminUser) clone(db *gorm.DB) xAdminUser {
	x.xAdminUserDo.ReplaceConnPool(db.Statement.ConnPool)
	return x
}

func (x xAdminUser) replaceDB(db *gorm.DB) xAdminUser {
	x.xAdminUserDo.ReplaceDB(db)
	return x
}

type xAdminUserDo struct{ gen.DO }

type IXAdminUserDo interface {
	gen.SubQuery
	Debug() IXAdminUserDo
	WithContext(ctx context.Context) IXAdminUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IXAdminUserDo
	WriteDB() IXAdminUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IXAdminUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IXAdminUserDo
	Not(conds ...gen.Condition) IXAdminUserDo
	Or(conds ...gen.Condition) IXAdminUserDo
	Select(conds ...field.Expr) IXAdminUserDo
	Where(conds ...gen.Condition) IXAdminUserDo
	Order(conds ...field.Expr) IXAdminUserDo
	Distinct(cols ...field.Expr) IXAdminUserDo
	Omit(cols ...field.Expr) IXAdminUserDo
	Join(table schema.Tabler, on ...field.Expr) IXAdminUserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IXAdminUserDo
	RightJoin(table schema.Tabler, on ...field.Expr) IXAdminUserDo
	Group(cols ...field.Expr) IXAdminUserDo
	Having(conds ...gen.Condition) IXAdminUserDo
	Limit(limit int) IXAdminUserDo
	Offset(offset int) IXAdminUserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IXAdminUserDo
	Unscoped() IXAdminUserDo
	Create(values ...*model.XAdminUser) error
	CreateInBatches(values []*model.XAdminUser, batchSize int) error
	Save(values ...*model.XAdminUser) error
	First() (*model.XAdminUser, error)
	Take() (*model.XAdminUser, error)
	Last() (*model.XAdminUser, error)
	Find() ([]*model.XAdminUser, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XAdminUser, err error)
	FindInBatches(result *[]*model.XAdminUser, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.XAdminUser) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IXAdminUserDo
	Assign(attrs ...field.AssignExpr) IXAdminUserDo
	Joins(fields ...field.RelationField) IXAdminUserDo
	Preload(fields ...field.RelationField) IXAdminUserDo
	FirstOrInit() (*model.XAdminUser, error)
	FirstOrCreate() (*model.XAdminUser, error)
	FindByPage(offset int, limit int) (result []*model.XAdminUser, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IXAdminUserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (x xAdminUserDo) Debug() IXAdminUserDo {
	return x.withDO(x.DO.Debug())
}

func (x xAdminUserDo) WithContext(ctx context.Context) IXAdminUserDo {
	return x.withDO(x.DO.WithContext(ctx))
}

func (x xAdminUserDo) ReadDB() IXAdminUserDo {
	return x.Clauses(dbresolver.Read)
}

func (x xAdminUserDo) WriteDB() IXAdminUserDo {
	return x.Clauses(dbresolver.Write)
}

func (x xAdminUserDo) Session(config *gorm.Session) IXAdminUserDo {
	return x.withDO(x.DO.Session(config))
}

func (x xAdminUserDo) Clauses(conds ...clause.Expression) IXAdminUserDo {
	return x.withDO(x.DO.Clauses(conds...))
}

func (x xAdminUserDo) Returning(value interface{}, columns ...string) IXAdminUserDo {
	return x.withDO(x.DO.Returning(value, columns...))
}

func (x xAdminUserDo) Not(conds ...gen.Condition) IXAdminUserDo {
	return x.withDO(x.DO.Not(conds...))
}

func (x xAdminUserDo) Or(conds ...gen.Condition) IXAdminUserDo {
	return x.withDO(x.DO.Or(conds...))
}

func (x xAdminUserDo) Select(conds ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Select(conds...))
}

func (x xAdminUserDo) Where(conds ...gen.Condition) IXAdminUserDo {
	return x.withDO(x.DO.Where(conds...))
}

func (x xAdminUserDo) Order(conds ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Order(conds...))
}

func (x xAdminUserDo) Distinct(cols ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Distinct(cols...))
}

func (x xAdminUserDo) Omit(cols ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Omit(cols...))
}

func (x xAdminUserDo) Join(table schema.Tabler, on ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Join(table, on...))
}

func (x xAdminUserDo) LeftJoin(table schema.Tabler, on ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.LeftJoin(table, on...))
}

func (x xAdminUserDo) RightJoin(table schema.Tabler, on ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.RightJoin(table, on...))
}

func (x xAdminUserDo) Group(cols ...field.Expr) IXAdminUserDo {
	return x.withDO(x.DO.Group(cols...))
}

func (x xAdminUserDo) Having(conds ...gen.Condition) IXAdminUserDo {
	return x.withDO(x.DO.Having(conds...))
}

func (x xAdminUserDo) Limit(limit int) IXAdminUserDo {
	return x.withDO(x.DO.Limit(limit))
}

func (x xAdminUserDo) Offset(offset int) IXAdminUserDo {
	return x.withDO(x.DO.Offset(offset))
}

func (x xAdminUserDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IXAdminUserDo {
	return x.withDO(x.DO.Scopes(funcs...))
}

func (x xAdminUserDo) Unscoped() IXAdminUserDo {
	return x.withDO(x.DO.Unscoped())
}

func (x xAdminUserDo) Create(values ...*model.XAdminUser) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Create(values)
}

func (x xAdminUserDo) CreateInBatches(values []*model.XAdminUser, batchSize int) error {
	return x.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (x xAdminUserDo) Save(values ...*model.XAdminUser) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Save(values)
}

func (x xAdminUserDo) First() (*model.XAdminUser, error) {
	if result, err := x.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.XAdminUser), nil
	}
}

func (x xAdminUserDo) Take() (*model.XAdminUser, error) {
	if result, err := x.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.XAdminUser), nil
	}
}

func (x xAdminUserDo) Last() (*model.XAdminUser, error) {
	if result, err := x.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.XAdminUser), nil
	}
}

func (x xAdminUserDo) Find() ([]*model.XAdminUser, error) {
	result, err := x.DO.Find()
	return result.([]*model.XAdminUser), err
}

func (x xAdminUserDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XAdminUser, err error) {
	buf := make([]*model.XAdminUser, 0, batchSize)
	err = x.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (x xAdminUserDo) FindInBatches(result *[]*model.XAdminUser, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return x.DO.FindInBatches(result, batchSize, fc)
}

func (x xAdminUserDo) Attrs(attrs ...field.AssignExpr) IXAdminUserDo {
	return x.withDO(x.DO.Attrs(attrs...))
}

func (x xAdminUserDo) Assign(attrs ...field.AssignExpr) IXAdminUserDo {
	return x.withDO(x.DO.Assign(attrs...))
}

func (x xAdminUserDo) Joins(fields ...field.RelationField) IXAdminUserDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Joins(_f))
	}
	return &x
}

func (x xAdminUserDo) Preload(fields ...field.RelationField) IXAdminUserDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Preload(_f))
	}
	return &x
}

func (x xAdminUserDo) FirstOrInit() (*model.XAdminUser, error) {
	if result, err := x.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.XAdminUser), nil
	}
}

func (x xAdminUserDo) FirstOrCreate() (*model.XAdminUser, error) {
	if result, err := x.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.XAdminUser), nil
	}
}

func (x xAdminUserDo) FindByPage(offset int, limit int) (result []*model.XAdminUser, count int64, err error) {
	result, err = x.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = x.Offset(-1).Limit(-1).Count()
	return
}

func (x xAdminUserDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = x.Count()
	if err != nil {
		return
	}

	err = x.Offset(offset).Limit(limit).Scan(result)
	return
}

func (x xAdminUserDo) Scan(result interface{}) (err error) {
	return x.DO.Scan(result)
}

func (x xAdminUserDo) Delete(models ...*model.XAdminUser) (result gen.ResultInfo, err error) {
	return x.DO.Delete(models)
}

func (x *xAdminUserDo) withDO(do gen.Dao) *xAdminUserDo {
	x.DO = *do.(*gen.DO)
	return x
}
