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

func newXTest(db *gorm.DB, opts ...gen.DOOption) xTest {
	_xTest := xTest{}

	_xTest.xTestDo.UseDB(db, opts...)
	_xTest.xTestDo.UseModel(&model.XTest{})

	tableName := _xTest.xTestDo.TableName()
	_xTest.ALL = field.NewAsterisk(tableName)
	_xTest.A = field.NewString(tableName, "a")
	_xTest.B = field.NewString(tableName, "b")
	_xTest.C = field.NewString(tableName, "c")

	_xTest.fillFieldMap()

	return _xTest
}

type xTest struct {
	xTestDo

	ALL field.Asterisk
	A   field.String
	B   field.String
	C   field.String

	fieldMap map[string]field.Expr
}

func (x xTest) Table(newTableName string) *xTest {
	x.xTestDo.UseTable(newTableName)
	return x.updateTableName(newTableName)
}

func (x xTest) As(alias string) *xTest {
	x.xTestDo.DO = *(x.xTestDo.As(alias).(*gen.DO))
	return x.updateTableName(alias)
}

func (x *xTest) updateTableName(table string) *xTest {
	x.ALL = field.NewAsterisk(table)
	x.A = field.NewString(table, "a")
	x.B = field.NewString(table, "b")
	x.C = field.NewString(table, "c")

	x.fillFieldMap()

	return x
}

func (x *xTest) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := x.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (x *xTest) fillFieldMap() {
	x.fieldMap = make(map[string]field.Expr, 3)
	x.fieldMap["a"] = x.A
	x.fieldMap["b"] = x.B
	x.fieldMap["c"] = x.C
}

func (x xTest) clone(db *gorm.DB) xTest {
	x.xTestDo.ReplaceConnPool(db.Statement.ConnPool)
	return x
}

func (x xTest) replaceDB(db *gorm.DB) xTest {
	x.xTestDo.ReplaceDB(db)
	return x
}

type xTestDo struct{ gen.DO }

type IXTestDo interface {
	gen.SubQuery
	Debug() IXTestDo
	WithContext(ctx context.Context) IXTestDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IXTestDo
	WriteDB() IXTestDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IXTestDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IXTestDo
	Not(conds ...gen.Condition) IXTestDo
	Or(conds ...gen.Condition) IXTestDo
	Select(conds ...field.Expr) IXTestDo
	Where(conds ...gen.Condition) IXTestDo
	Order(conds ...field.Expr) IXTestDo
	Distinct(cols ...field.Expr) IXTestDo
	Omit(cols ...field.Expr) IXTestDo
	Join(table schema.Tabler, on ...field.Expr) IXTestDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IXTestDo
	RightJoin(table schema.Tabler, on ...field.Expr) IXTestDo
	Group(cols ...field.Expr) IXTestDo
	Having(conds ...gen.Condition) IXTestDo
	Limit(limit int) IXTestDo
	Offset(offset int) IXTestDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IXTestDo
	Unscoped() IXTestDo
	Create(values ...*model.XTest) error
	CreateInBatches(values []*model.XTest, batchSize int) error
	Save(values ...*model.XTest) error
	First() (*model.XTest, error)
	Take() (*model.XTest, error)
	Last() (*model.XTest, error)
	Find() ([]*model.XTest, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XTest, err error)
	FindInBatches(result *[]*model.XTest, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.XTest) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IXTestDo
	Assign(attrs ...field.AssignExpr) IXTestDo
	Joins(fields ...field.RelationField) IXTestDo
	Preload(fields ...field.RelationField) IXTestDo
	FirstOrInit() (*model.XTest, error)
	FirstOrCreate() (*model.XTest, error)
	FindByPage(offset int, limit int) (result []*model.XTest, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IXTestDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (x xTestDo) Debug() IXTestDo {
	return x.withDO(x.DO.Debug())
}

func (x xTestDo) WithContext(ctx context.Context) IXTestDo {
	return x.withDO(x.DO.WithContext(ctx))
}

func (x xTestDo) ReadDB() IXTestDo {
	return x.Clauses(dbresolver.Read)
}

func (x xTestDo) WriteDB() IXTestDo {
	return x.Clauses(dbresolver.Write)
}

func (x xTestDo) Session(config *gorm.Session) IXTestDo {
	return x.withDO(x.DO.Session(config))
}

func (x xTestDo) Clauses(conds ...clause.Expression) IXTestDo {
	return x.withDO(x.DO.Clauses(conds...))
}

func (x xTestDo) Returning(value interface{}, columns ...string) IXTestDo {
	return x.withDO(x.DO.Returning(value, columns...))
}

func (x xTestDo) Not(conds ...gen.Condition) IXTestDo {
	return x.withDO(x.DO.Not(conds...))
}

func (x xTestDo) Or(conds ...gen.Condition) IXTestDo {
	return x.withDO(x.DO.Or(conds...))
}

func (x xTestDo) Select(conds ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Select(conds...))
}

func (x xTestDo) Where(conds ...gen.Condition) IXTestDo {
	return x.withDO(x.DO.Where(conds...))
}

func (x xTestDo) Order(conds ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Order(conds...))
}

func (x xTestDo) Distinct(cols ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Distinct(cols...))
}

func (x xTestDo) Omit(cols ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Omit(cols...))
}

func (x xTestDo) Join(table schema.Tabler, on ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Join(table, on...))
}

func (x xTestDo) LeftJoin(table schema.Tabler, on ...field.Expr) IXTestDo {
	return x.withDO(x.DO.LeftJoin(table, on...))
}

func (x xTestDo) RightJoin(table schema.Tabler, on ...field.Expr) IXTestDo {
	return x.withDO(x.DO.RightJoin(table, on...))
}

func (x xTestDo) Group(cols ...field.Expr) IXTestDo {
	return x.withDO(x.DO.Group(cols...))
}

func (x xTestDo) Having(conds ...gen.Condition) IXTestDo {
	return x.withDO(x.DO.Having(conds...))
}

func (x xTestDo) Limit(limit int) IXTestDo {
	return x.withDO(x.DO.Limit(limit))
}

func (x xTestDo) Offset(offset int) IXTestDo {
	return x.withDO(x.DO.Offset(offset))
}

func (x xTestDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IXTestDo {
	return x.withDO(x.DO.Scopes(funcs...))
}

func (x xTestDo) Unscoped() IXTestDo {
	return x.withDO(x.DO.Unscoped())
}

func (x xTestDo) Create(values ...*model.XTest) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Create(values)
}

func (x xTestDo) CreateInBatches(values []*model.XTest, batchSize int) error {
	return x.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (x xTestDo) Save(values ...*model.XTest) error {
	if len(values) == 0 {
		return nil
	}
	return x.DO.Save(values)
}

func (x xTestDo) First() (*model.XTest, error) {
	if result, err := x.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.XTest), nil
	}
}

func (x xTestDo) Take() (*model.XTest, error) {
	if result, err := x.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.XTest), nil
	}
}

func (x xTestDo) Last() (*model.XTest, error) {
	if result, err := x.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.XTest), nil
	}
}

func (x xTestDo) Find() ([]*model.XTest, error) {
	result, err := x.DO.Find()
	return result.([]*model.XTest), err
}

func (x xTestDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.XTest, err error) {
	buf := make([]*model.XTest, 0, batchSize)
	err = x.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (x xTestDo) FindInBatches(result *[]*model.XTest, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return x.DO.FindInBatches(result, batchSize, fc)
}

func (x xTestDo) Attrs(attrs ...field.AssignExpr) IXTestDo {
	return x.withDO(x.DO.Attrs(attrs...))
}

func (x xTestDo) Assign(attrs ...field.AssignExpr) IXTestDo {
	return x.withDO(x.DO.Assign(attrs...))
}

func (x xTestDo) Joins(fields ...field.RelationField) IXTestDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Joins(_f))
	}
	return &x
}

func (x xTestDo) Preload(fields ...field.RelationField) IXTestDo {
	for _, _f := range fields {
		x = *x.withDO(x.DO.Preload(_f))
	}
	return &x
}

func (x xTestDo) FirstOrInit() (*model.XTest, error) {
	if result, err := x.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.XTest), nil
	}
}

func (x xTestDo) FirstOrCreate() (*model.XTest, error) {
	if result, err := x.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.XTest), nil
	}
}

func (x xTestDo) FindByPage(offset int, limit int) (result []*model.XTest, count int64, err error) {
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

func (x xTestDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = x.Count()
	if err != nil {
		return
	}

	err = x.Offset(offset).Limit(limit).Scan(result)
	return
}

func (x xTestDo) Scan(result interface{}) (err error) {
	return x.DO.Scan(result)
}

func (x xTestDo) Delete(models ...*model.XTest) (result gen.ResultInfo, err error) {
	return x.DO.Delete(models)
}

func (x *xTestDo) withDO(do gen.Dao) *xTestDo {
	x.DO = *do.(*gen.DO)
	return x
}
