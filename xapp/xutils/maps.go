package xutils

import "encoding/json"

type XMaps struct {
	RawData []map[string]interface{}
}

func (this *XMaps) FromBytes(bytes []byte) error {
	this.RawData = []map[string]interface{}{}
	err := json.Unmarshal(bytes, &this.RawData)
	return err
}

func (this *XMaps) ToBytes() []byte {
	b, _ := json.Marshal(this.RawData)
	return b
}

func (this *XMaps) FromObjects(obj []interface{}) error {
	this.RawData = []map[string]interface{}{}
	for i := 0; i < len(obj); i++ {
		this.RawData = append(this.RawData, obj[i].(map[string]interface{}))
	}
	return nil
}

func (this *XMaps) Length() int {
	if this.RawData == nil {
		return 0
	}
	return len(this.RawData)
}

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
