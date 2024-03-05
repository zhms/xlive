package xutils

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type XMap struct {
	RawData map[string]interface{}
}

// 从json串解析获得XMap
func (this *XMap) FromBytes(bytes []byte) error {
	this.RawData = map[string]interface{}{}
	return json.Unmarshal(bytes, &this.RawData)
}

// 从go对象解析获得XMap
func (this *XMap) FromObject(obj interface{}) error {
	bytes, err := json.Marshal(&obj)
	if err != nil {
		return err
	}
	this.RawData = map[string]interface{}{}
	json.Unmarshal(bytes, &this.RawData)
	return json.Unmarshal(bytes, &this.RawData)
}

func (this *XMap) map_field(field string) interface{} {
	if this.RawData == nil {
		return nil
	}
	return (this.RawData)[field]
}

// 获取原始map数据
func (this *XMap) Map() *map[string]interface{} {
	return &this.RawData
}

// 获取字段 转换为int型
func (this *XMap) Int(field string) int {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int(ToInt(data))
}

// 获取字段 转换为int32型
func (this *XMap) Int32(field string) int32 {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int32(ToInt(data))
}

// 获取字段 转换为int64型
func (this *XMap) Int64(field string) int64 {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int64(ToInt(data))
}

// 获取字段 转换为string型
func (this *XMap) String(field string) string {
	data := this.map_field(field)
	if data == nil {
		return ""
	}
	return ToString(data)
}

// 获取字段 转换为decimal型
func (this *XMap) Decimal(field string) decimal.Decimal {
	data := this.map_field(field)
	if data == nil {
		return decimal.Zero
	}
	r, e := decimal.NewFromString(ToString(data))
	if e != nil {
		return decimal.Zero
	}
	return r
}

// 获取字段 转换为[]byte
func (this *XMap) Bytes(field string) []byte {
	data := this.map_field(field)
	if data == nil {
		return []byte{}
	}
	return []byte(ToString(data))
}

// 删除字段
func (this *XMap) Delete(field string) {
	if this.RawData == nil {
		return
	}
	delete(this.RawData, field)
}

// 设置字段值
func (this *XMap) Set(field string, value interface{}) {
	if this.RawData == nil {
		this.RawData = map[string]interface{}{}
	}
	this.RawData[field] = value
}

// 遍历map
func (this *XMap) ForEach(cb func(string, interface{}) bool) {
	if this.RawData == nil {
		return
	}
	for k, v := range this.RawData {
		if !cb(k, v) {
			break
		}
	}
}

// 字段是否存在
func (this *XMap) Exists(field string) bool {
	if this.RawData == nil {
		return false
	}
	_, ok := this.RawData[field]
	return ok
}

// 字段是否存在
func (this *XMap) ToObject(data any) {
	if this.RawData == nil {
		return
	}
	jdata, _ := json.Marshal(this.RawData)
	json.Unmarshal(jdata, data)
}

func MapSet(m *map[string]interface{}, key string, value interface{}, invalidvalue interface{}) {
	if m == nil {
		return
	}
	if value == invalidvalue {
		return
	}
	(*m)[key] = value
}

func MapSetIn(m *map[string]interface{}, key string, value interface{}, validvalue []interface{}) {
	if m == nil {
		return
	}
	finded := false
	for _, v := range validvalue {
		if v == value {
			finded = true
			break
		}
	}
	if !finded {
		return
	}
	(*m)[key] = value
}

func (this *XMap) ToString() string {
	bytes, _ := json.Marshal(this.RawData)
	return string(bytes)
}
