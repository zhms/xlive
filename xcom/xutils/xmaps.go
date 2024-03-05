package xutils

import "encoding/json"

type XMaps struct {
	RawData []map[string]interface{}
}

// 从json串,解析获得XMaps
func (this *XMaps) FromBytes(bytes []byte) error {
	this.RawData = []map[string]interface{}{}
	err := json.Unmarshal(bytes, &this.RawData)
	return err
}

// 从go对象,解析获得XMaps
func (this *XMaps) FromObjects(obj []interface{}) error {
	this.RawData = []map[string]interface{}{}
	for i := 0; i < len(obj); i++ {
		this.RawData = append(this.RawData, obj[i].(map[string]interface{}))
	}
	return nil
}

// 获取原始map切片 []map[string]interface{}
func (this *XMaps) Maps() *[]map[string]interface{} {
	if this.RawData == nil {
		return nil
	}
	return &this.RawData
}

// 获取切片长度
func (this *XMaps) Length() int {
	if this.RawData == nil {
		return 0
	}
	return len(this.RawData)
}

// 根据下标获取元素
func (this *XMaps) Index(index int) *XMap {
	if this.RawData == nil {
		return nil
	}
	if index < 0 {
		return nil
	}
	if index >= len(this.RawData) {
		return nil
	}
	return &XMap{RawData: this.RawData[index]}
}

// 根据下标删除元素
func (this *XMaps) Remove(index int) {
	if this.RawData == nil {
		return
	}
	if index < 0 {
		return
	}
	if index >= len(this.RawData) {
		return
	}
	this.RawData = append(this.RawData[:index], this.RawData[index+1:]...)
}

// 遍历元素
func (this *XMaps) ForEach(cb func(*XMap) bool) {
	if this.RawData == nil {
		return
	}
	for i := 0; i < len(this.RawData); i++ {
		if !cb(&XMap{RawData: this.RawData[i]}) {
			break
		}
	}
}

func (this *XMaps) ToString() string {
	if this.RawData == nil {
		return ""
	}
	bytes, _ := json.Marshal(this.RawData)
	return string(bytes)
}
