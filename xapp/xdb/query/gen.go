// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q              = new(Query)
	XAdminLoginLog *xAdminLoginLog
	XAdminOptLog   *xAdminOptLog
	XAdminRole     *xAdminRole
	XAdminUser     *xAdminUser
	XChat          *xChat
	XChatBanIP     *xChatBanIP
	XHongbao       *xHongbao
	XHostSeller    *xHostSeller
	XKv            *xKv
	XLiveProvider  *xLiveProvider
	XLiveRoom      *xLiveRoom
	XSeller        *xSeller
	XUser          *xUser
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	XAdminLoginLog = &Q.XAdminLoginLog
	XAdminOptLog = &Q.XAdminOptLog
	XAdminRole = &Q.XAdminRole
	XAdminUser = &Q.XAdminUser
	XChat = &Q.XChat
	XChatBanIP = &Q.XChatBanIP
	XHongbao = &Q.XHongbao
	XHostSeller = &Q.XHostSeller
	XKv = &Q.XKv
	XLiveProvider = &Q.XLiveProvider
	XLiveRoom = &Q.XLiveRoom
	XSeller = &Q.XSeller
	XUser = &Q.XUser
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:             db,
		XAdminLoginLog: newXAdminLoginLog(db, opts...),
		XAdminOptLog:   newXAdminOptLog(db, opts...),
		XAdminRole:     newXAdminRole(db, opts...),
		XAdminUser:     newXAdminUser(db, opts...),
		XChat:          newXChat(db, opts...),
		XChatBanIP:     newXChatBanIP(db, opts...),
		XHongbao:       newXHongbao(db, opts...),
		XHostSeller:    newXHostSeller(db, opts...),
		XKv:            newXKv(db, opts...),
		XLiveProvider:  newXLiveProvider(db, opts...),
		XLiveRoom:      newXLiveRoom(db, opts...),
		XSeller:        newXSeller(db, opts...),
		XUser:          newXUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	XAdminLoginLog xAdminLoginLog
	XAdminOptLog   xAdminOptLog
	XAdminRole     xAdminRole
	XAdminUser     xAdminUser
	XChat          xChat
	XChatBanIP     xChatBanIP
	XHongbao       xHongbao
	XHostSeller    xHostSeller
	XKv            xKv
	XLiveProvider  xLiveProvider
	XLiveRoom      xLiveRoom
	XSeller        xSeller
	XUser          xUser
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		XAdminLoginLog: q.XAdminLoginLog.clone(db),
		XAdminOptLog:   q.XAdminOptLog.clone(db),
		XAdminRole:     q.XAdminRole.clone(db),
		XAdminUser:     q.XAdminUser.clone(db),
		XChat:          q.XChat.clone(db),
		XChatBanIP:     q.XChatBanIP.clone(db),
		XHongbao:       q.XHongbao.clone(db),
		XHostSeller:    q.XHostSeller.clone(db),
		XKv:            q.XKv.clone(db),
		XLiveProvider:  q.XLiveProvider.clone(db),
		XLiveRoom:      q.XLiveRoom.clone(db),
		XSeller:        q.XSeller.clone(db),
		XUser:          q.XUser.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		XAdminLoginLog: q.XAdminLoginLog.replaceDB(db),
		XAdminOptLog:   q.XAdminOptLog.replaceDB(db),
		XAdminRole:     q.XAdminRole.replaceDB(db),
		XAdminUser:     q.XAdminUser.replaceDB(db),
		XChat:          q.XChat.replaceDB(db),
		XChatBanIP:     q.XChatBanIP.replaceDB(db),
		XHongbao:       q.XHongbao.replaceDB(db),
		XHostSeller:    q.XHostSeller.replaceDB(db),
		XKv:            q.XKv.replaceDB(db),
		XLiveProvider:  q.XLiveProvider.replaceDB(db),
		XLiveRoom:      q.XLiveRoom.replaceDB(db),
		XSeller:        q.XSeller.replaceDB(db),
		XUser:          q.XUser.replaceDB(db),
	}
}

type queryCtx struct {
	XAdminLoginLog IXAdminLoginLogDo
	XAdminOptLog   IXAdminOptLogDo
	XAdminRole     IXAdminRoleDo
	XAdminUser     IXAdminUserDo
	XChat          IXChatDo
	XChatBanIP     IXChatBanIPDo
	XHongbao       IXHongbaoDo
	XHostSeller    IXHostSellerDo
	XKv            IXKvDo
	XLiveProvider  IXLiveProviderDo
	XLiveRoom      IXLiveRoomDo
	XSeller        IXSellerDo
	XUser          IXUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		XAdminLoginLog: q.XAdminLoginLog.WithContext(ctx),
		XAdminOptLog:   q.XAdminOptLog.WithContext(ctx),
		XAdminRole:     q.XAdminRole.WithContext(ctx),
		XAdminUser:     q.XAdminUser.WithContext(ctx),
		XChat:          q.XChat.WithContext(ctx),
		XChatBanIP:     q.XChatBanIP.WithContext(ctx),
		XHongbao:       q.XHongbao.WithContext(ctx),
		XHostSeller:    q.XHostSeller.WithContext(ctx),
		XKv:            q.XKv.WithContext(ctx),
		XLiveProvider:  q.XLiveProvider.WithContext(ctx),
		XLiveRoom:      q.XLiveRoom.WithContext(ctx),
		XSeller:        q.XSeller.WithContext(ctx),
		XUser:          q.XUser.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
