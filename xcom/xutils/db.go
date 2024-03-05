package xutils

import (
	"database/sql"

	"github.com/beego/beego/logs"
	"gorm.io/gorm"
)

// 获取sql.Rows返回的一条数据
func DbFirst(rows *sql.Rows) *XMap {
	if rows == nil {
		return nil
	}
	data := XMap{}
	if rows.Next() {
		data.RawData = *dbgetone(rows)
	} else {
		return nil
	}
	rows.Close()
	return &data
}

func DbFirstEx(db *gorm.DB) (*XMap, error) {
	rows, err := db.Rows()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	data := XMap{}
	if rows.Next() {
		data.RawData = *dbgetone(rows)
	} else {
		return nil, nil
	}
	rows.Close()
	return &data, nil
}

// 获取sql.Rows返回的数据
func DbResult(rows *sql.Rows) *XMaps {
	if rows == nil {
		return nil
	}
	data := XMaps{}
	for rows.Next() {
		data.RawData = append(data.RawData, *dbgetone(rows))
	}
	rows.Close()
	return &data
}

func DbResultEx(db gorm.DB) (*XMaps, error) {
	rows, err := db.Rows()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	data := XMaps{}
	for rows.Next() {
		data.RawData = append(data.RawData, *dbgetone(rows))
	}
	rows.Close()
	return &data, nil
}

func dbgetone(rows *sql.Rows) *map[string]interface{} {
	data := make(map[string]interface{})
	fields, _ := rows.Columns()
	scans := make([]interface{}, len(fields))
	for i := range scans {
		scans[i] = &scans[i]
	}
	err := rows.Scan(scans...)
	if err != nil {
		logs.Error(err)
		return nil
	}
	for i := range fields {
		if scans[i] != nil {
			data[fields[i]] = ToString(scans[i])
		} else {
			data[fields[i]] = nil
		}
	}
	return &data
}

func DbWhere(db *gorm.DB, field string, value interface{}, invalidvalue interface{}) *gorm.DB {
	if value != invalidvalue {
		db = db.Where(field, value)
	}
	return db
}
