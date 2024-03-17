package xutils

import (
	"encoding/json"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type XMap struct {
	RawData map[string]interface{}
}

func (this *XMap) FromBytes(bytes []byte) error {
	this.RawData = make(map[string]interface{})
	return json.Unmarshal(bytes, &this.RawData)
}

func (this *XMap) ToBytes() ([]byte, error) {
	return json.Marshal(this.RawData)
}

func (this *XMap) FromObject(obj interface{}) error {
	bytes, err := json.Marshal(&obj)
	if err != nil {
		return err
	}
	this.RawData = make(map[string]interface{})
	json.Unmarshal(bytes, &this.RawData)
	return json.Unmarshal(bytes, &this.RawData)
}

func (this *XMap) ToObject(data any) {
	if this.RawData == nil {
		return
	}
	jdata, _ := json.Marshal(this.RawData)
	json.Unmarshal(jdata, data)
}

func (this *XMap) map_field(field string) interface{} {
	if this.RawData == nil {
		return nil
	}
	return (this.RawData)[field]
}

func (this *XMap) Map() *map[string]interface{} {
	return &this.RawData
}

func (this *XMap) Int(field string) int {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int(cast.ToInt(data))
}

func (this *XMap) Int32(field string) int32 {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int32(cast.ToInt(data))
}

func (this *XMap) Int64(field string) int64 {
	data := this.map_field(field)
	if data == nil {
		return 0
	}
	return int64(cast.ToInt(data))
}

func (this *XMap) String(field string) string {
	data := this.map_field(field)
	if data == nil {
		return ""
	}
	return cast.ToString(data)
}

func (this *XMap) Decimal(field string) decimal.Decimal {
	data := this.map_field(field)
	if data == nil {
		return decimal.Zero
	}
	r, e := decimal.NewFromString(cast.ToString(data))
	if e != nil {
		return decimal.Zero
	}
	return r
}

func (this *XMap) Bytes(field string) []byte {
	data := this.map_field(field)
	if data == nil {
		return []byte{}
	}
	return []byte(cast.ToString(data))
}

func (this *XMap) Delete(field string) {
	if this.RawData == nil {
		return
	}
	delete(this.RawData, field)
}

func (this *XMap) Set(field string, value interface{}) {
	if this.RawData == nil {
		this.RawData = make(map[string]interface{})
	}
	this.RawData[field] = value
}

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

func (this *XMap) Exists(field string) bool {
	if this.RawData == nil {
		return false
	}
	_, ok := this.RawData[field]
	return ok
}

func (this *XMap) SetEx(key string, value interface{}, invalidvalue interface{}) {
	if value == invalidvalue {
		return
	}
	this.RawData[key] = value
}

func (this *XMap) SetInEx(key string, value interface{}, validvalue []interface{}) {
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
	this.RawData[key] = value
}
