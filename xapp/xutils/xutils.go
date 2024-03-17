package xutils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/beego/beego/logs"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

func ToDecimal(v interface{}) decimal.Decimal {
	return decimal.NewFromFloat(cast.ToFloat64(v))
}

func ObjectToMap(obj interface{}) *map[string]interface{} {
	bytes, err := json.Marshal(obj)
	if err != nil {
		logs.Error("ObjectToMap:", err)
		return nil
	}
	data := make(map[string]interface{})
	json.Unmarshal(bytes, &data)
	return &data
}

func UtcOffset() int {
	currentTime := time.Now()
	_, offset := currentTime.Zone()
	utcTime := currentTime.UTC()
	_, utcOffset := utcTime.Zone()
	return int((time.Duration(offset-utcOffset) * time.Second).Hours())
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
