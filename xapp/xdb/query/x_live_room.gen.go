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

func newXLiveRoom(db *gorm.DB, opts ...gen.DOOption) xLiveRoom {
	_xLiveRoom := xLiveRoom{}

	_xLiveRoom.xLiveRoomDo.UseDB(db, opts...)
	_xLiveRoom.xLiveRoomDo.UseModel(&model.XLiveRoom{})

	tableName := _xLiveRoom.xLiveRoomDo.TableName()
	_xLiveRoom.ALL = field.NewAsterisk(tableName)
	_xLiveRoom.ID = field.NewInt32(tableName, "id")
	_xLiveRoom.SellerID = field.NewInt32(tableName, "seller_id")
	_xLiveRoom.Name = field.NewString(tableName, "name")
	_xLiveRoom.Account = field.NewString(tableName, "account")
	_xLiveRoom.Title = field.NewString(tableName, "title")
	_xLiveRoom.PushURL = field.NewString(tableName, "push_url")
	_xLiveRoom.PullURL = field.NewString(tableName, "pull_url")
	_xLiveRoom.LiveURL = field.NewString(tableName, "live_url")
	_xLiveRoom.State = field.NewInt32(tableName, "state")
	_xLiveRoom.CreateTime = field.NewTime(tableName, "create_time")

	_xLiveRoom.fillFieldMap()

	return _xLiveRoom
}

type xLiveRoom struct {
	xLiveRoomDo

	ALL        field.Asterisk
	ID         field.Int32  // id
	SellerID   field.Int32  // 运营商
	Name       field.String // 直播间名称
	Account    field.String // 主播账号
	Title      field.String // 直播间标题
	PushURL    field.String // 推流地址
	PullURL    field.String // 拉流地址
	LiveURL    field.String // 前端地址
	State      field.Int32  // 状态1正在直播,2未在直播
	CreateTime field.Time   // 创建时间

	fieldMap map[string]field.Expr
}

func (x xLiveRoom) Table(newTableName string) *xLiveRoom {
	x.xLiveRoomDo.UseTable(newTableName)
	return x.updateTableName(newTableName)
}

func (x xLiveRoom) As(alias string) *xLiveRoom {
	x.xLiveRoomDo.DO = *(x.xLiveRoomDo.As(alias).(*gen.DO))
	return x.updateTableName(alias)
}

func (x *xLiveRoom) updateTableName(table string) *xLiveRoom {
	x.ALL = field.NewAsterisk(table)
	x.ID = field.NewInt32(table, "id")
	x.SellerID = field.NewInt32(table, "seller_id")
	x.Name = field.NewString(table, "name")
	x.Account = field.NewString(table, "account")
	x.Title = field.NewString(table, "title")
	x.PushURL = field.NewString(table, "push_url")
	x.PullURL = field.NewString(table, "pull_url")
	x.LiveURL = field.NewString(table, "live_url")
	x.State = field.NewInt32(table, "state")
	x.CreateTime = field.NewTime(table, "create_time")

	x.fillFieldMap()

	return x
}

func (x *xLiveRoom) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := x.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (x *xLiveRoom) fillFieldMap() {
	x.fieldMap = make(map[string]field.Expr, 10)
	x.fieldMap["id"] = x.ID
	x.fieldMap["seller_id"] = x.SellerID
	x.fieldMap["name"] = x.Name
	x.fieldMap["account"] = x.Account
	x.fieldMap["title"] = x.Title
	x.fieldMap["push_url"] = x.PushURL
	x.fieldMap["pull_url"] = x.PullURL
	x.fieldMap["live_url"] = x.LiveURL
	x.fieldMap["state"] = x.State
	x.fieldMap["create_time"] = x.CreateTime
}

func (x xLiveRoom) clone(db *gorm.DB) xLiveRoom {
	x.xLiveRoomDo.ReplaceConnPool(db.Statement.ConnPool)
	return x
}

func (x xLiveRoom) replaceDB(db *gorm.DB) xLiveRoom {
	x.xLiveRoomDo.ReplaceDB(db)
	return x
}

type xLiveRoomDo struct{ gen.DO }

type IXLiveRoomDo interface {
	gen.SubQuery
	Debug() IXLiveRoomDo
	WithContext(ctx context.Context) IXLiveRoomDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IXLiveRoomDo
	WriteDB() IXLiveRoomDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IXLiveRoomDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IXLiveRoomDo
	Not(conds ...gen.Condition) IXLiveRoomDo
	Or(conds ...gen.Condition) IXLiveRoomDo
	Select(conds ...field.Expr) IXLiveRoomDo
	Where(conds ...gen.Condition) IXLiveRoomDo
	Order(conds ...field.Expr) IXLiveRoomDo
	Distinct(cols ...field.Expr) IXLiveRoomDo
	Omit(cols ...field.Expr) IXLiveRoomDo
	Join(table schema.Tabler, on ...field.Expr) IXLiveRoomDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IXLiveRoomDo
	RightJoin(table schema.Tabler, on ...field.Expr) IXLiveRoomDo
	Group(cols ...field.Expr) IXLiveRoomDo
	Having(conds ...gen.Condition) IXLiveRoomDo
	Limit(limit int) IXLiveRoomDo
	Offset(offset int) IXLiveRoomDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IXLiveRoomDo
	Unscoped() IXLiveRoomDo
	Create(values ...*model.XLiveRoom) error
	CreateInBatches(values []*model.XLiveRoom, batchSize int) error
	Save(values ...*model.XLiveRoom) error
	First() (*model.XLiveRoom, error)
	Take() (*model.XLiveRoom, error)
	Last() (*model.XLiveRoom, error)
	Find() ([]*model.XLiveRoom, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XLiveRoom, err error)
	FindInBatches(result *[]*model.XLiveRoom, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.XLiveRoom) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IXLiveRoomDo
	Assign(attrs ...field.AssignExpr) IXLiveRoomDo
	Joins(fields ...field.RelationField) IXLiveRoomDo
	Preload(fields ...field.RelationField) IXLiveRoomDo
	FirstOrInit() (*model.XLiveRoom, error)
	FirstOrCreate() (*model.XLiveRoom, error)
	FindByPage(offset int, limit int) (result []*model.XLiveRoom, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IXLiveRoomDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (x xLiveRoomDo) Debug() IXLiveRoomDo {
	return x.withDO(x.DO.Debug())
}

func (x xLiveRoomDo) WithContext(ctx context.Context) IXLiveRoomDo {
	return x.withDO(x.DO.WithContext(ctx))
}

func (x xLiveRoomDo) ReadDB() IXLiveRoomDo {
	return x.Clauses(dbresolver.Read)
}

func (x xLiveRoomDo) WriteDB() IXLiveRoomDo {
	return x.Clauses(dbresolver.Write)
}

func (x xLiveRoomDo) Session(config *gorm.Session) IXLiveRoomDo {
	return x.withDO(x.DO.Session(config))
}

func (x xLiveRoomDo) Clauses(conds ...clause.Expression) IXLiveRoomDo {
	return x.withDO(x.DO.Clauses(conds...))
}

func (x xLiveRoomDo) Returning(value interface{}, columns ...string) IXLiveRoomDo {
	return x.withDO(x.DO.Returning(value, columns...))
}

func (x xLiveRoomDo) Not(conds ...gen.Condition) IXLiveRoomDo {
	return x.withDO(x.DO.Not(conds...))
}

func (x xLiveRoomDo) Or(conds ...gen.Condition) IXLiveRoomDo {
	return x.withDO(x.DO.Or(conds...))
}

func (x xLiveRoomDo) Select(conds ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Select(conds...))
}

func (x xLiveRoomDo) Where(conds ...gen.Condition) IXLiveRoomDo {
	return x.withDO(x.DO.Where(conds...))
}

func (x xLiveRoomDo) Order(conds ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Order(conds...))
}

func (x xLiveRoomDo) Distinct(cols ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Distinct(cols...))
}

func (x xLiveRoomDo) Omit(cols ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Omit(cols...))
}

func (x xLiveRoomDo) Join(table schema.Tabler, on ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Join(table, on...))
}

func (x xLiveRoomDo) LeftJoin(table schema.Tabler, on ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.LeftJoin(table, on...))
}

func (x xLiveRoomDo) RightJoin(table schema.Tabler, on ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.RightJoin(table, on...))
}

func (x xLiveRoomDo) Group(cols ...field.Expr) IXLiveRoomDo {
	return x.withDO(x.DO.Group(cols...))
}

func (x xLiveRoomDo) Having(conds ...gen.Condition) IXLiveRoomDo {
	return x.withDO(x.DO.Having(conds...))
}

func (x xLiveRoomDo) Limit(limit int) IXLiveRoomDo {
	return x.withDO(x.DO.Limit(limit))
}

func (x xLiveRoomDo) Offset(offset int) IXLiveRoomDo {
	return x.withDO(x.DO.Offset(offset))
}

func (x xLiveRoomDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IXLiveRoomDo {
	return x.withDO(x.DO.Scopes(funcs...))
}

func (x xLiveRoomDo) Unscoped() IXLiveRoomDo {
	return x.withDO(x.DO.Unscoped())
}

func (x xLiveRoomDo) Create(values ...*model.XLiveRoom) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Create(values)
}

func (x xLiveRoomDo) CreateInBatches(values []*model.XLiveRoom, batchSize int) error {
	return x.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (x xLiveRoomDo) Save(values ...*model.XLiveRoom) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Save(values)
}

func (x xLiveRoomDo) First() (*model.XLiveRoom, error) {
	if result, err := x.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.XLiveRoom), nil
	}
}

func (x xLiveRoomDo) Take() (*model.XLiveRoom, error) {
	if result, err := x.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.XLiveRoom), nil
	}
}

func (x xLiveRoomDo) Last() (*model.XLiveRoom, error) {
	if result, err := x.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.XLiveRoom), nil
	}
}

func (x xLiveRoomDo) Find() ([]*model.XLiveRoom, error) {
	result, err := x.DO.Find()
	return result.([]*model.XLiveRoom), err
}

func (x xLiveRoomDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XLiveRoom, err error) {
	buf := make([]*model.XLiveRoom, 0, batchSize)
	err = x.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (x xLiveRoomDo) FindInBatches(result *[]*model.XLiveRoom, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return x.DO.FindInBatches(result, batchSize, fc)
}

func (x xLiveRoomDo) Attrs(attrs ...field.AssignExpr) IXLiveRoomDo {
	return x.withDO(x.DO.Attrs(attrs...))
}

func (x xLiveRoomDo) Assign(attrs ...field.AssignExpr) IXLiveRoomDo {
	return x.withDO(x.DO.Assign(attrs...))
}

func (x xLiveRoomDo) Joins(fields ...field.RelationField) IXLiveRoomDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Joins(_f))
	}
	return &x
}

func (x xLiveRoomDo) Preload(fields ...field.RelationField) IXLiveRoomDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Preload(_f))
	}
	return &x
}

func (x xLiveRoomDo) FirstOrInit() (*model.XLiveRoom, error) {
	if result, err := x.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.XLiveRoom), nil
	}
}

func (x xLiveRoomDo) FirstOrCreate() (*model.XLiveRoom, error) {
	if result, err := x.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.XLiveRoom), nil
	}
}

func (x xLiveRoomDo) FindByPage(offset int, limit int) (result []*model.XLiveRoom, count int64, err error) {
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

func (x xLiveRoomDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = x.Count()
	if err != nil {
		return
	}

	err = x.Offset(offset).Limit(limit).Scan(result)
	return
}

func (x xLiveRoomDo) Scan(result interface{}) (err error) {
	return x.DO.Scan(result)
}

func (x xLiveRoomDo) Delete(models ...*model.XLiveRoom) (result gen.ResultInfo, err error) {
	return x.DO.Delete(models)
}

func (x *xLiveRoomDo) withDO(do gen.Dao) *xLiveRoomDo {
	x.DO = *do.(*gen.DO)
	return x
}
