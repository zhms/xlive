package xutils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/beego/logs"
	"github.com/shopspring/decimal"
)

// interface转string
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return fmt.Sprint(v.(int))
	case int32:
		return fmt.Sprint(v.(int32))
	case int64:
		return fmt.Sprint(v.(int64))
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case decimal.Decimal:
		return v.(decimal.Decimal).String()
	case uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case uint32:
		return strconv.FormatUint(uint64(v.(uint32)), 10)
	case uint:
		return strconv.FormatUint(uint64(v.(uint)), 10)
	default:
		if bytes, ok := v.([]byte); ok {
			return string(bytes)
		}
	}
	return ""
}

// interface转int64
func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch v.(type) {
	case string:
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0
		}
		return i
	case int:
		return int64(v.(int))
	case int32:
		return int64(v.(int32))
	case int64:
		return int64(v.(int64))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case decimal.Decimal:
		return v.(decimal.Decimal).IntPart()
	default:
		if bytes, ok := v.([]byte); ok {
			i, err := strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				return 0
			}
			return i
		}
	}
	return 0
}

// interface转int
func ToInt(v interface{}) int {
	return int(ToInt64(v))
}

// interface转int32
func ToInt32(v interface{}) int32 {
	return int32(ToInt64(v))
}

func ToDecimal(v interface{}) decimal.Decimal {
	if v == nil {
		return decimal.Zero
	}
	switch v.(type) {
	case string:
		i, err := decimal.NewFromString(v.(string))
		if err != nil {
			return decimal.Zero
		}
		return i.Round(6)
	case int:
		return decimal.NewFromInt(int64(v.(int)))
	case int32:
		return decimal.NewFromInt(int64(v.(int32)))
	case int64:
		return decimal.NewFromInt(int64(v.(int64)))
	case float32:
		return decimal.NewFromFloat32(v.(float32))
	case float64:
		return decimal.NewFromFloat(v.(float64))
	case decimal.Decimal:
		return v.(decimal.Decimal)
	default:
		if bytes, ok := v.([]byte); ok {
			i, err := decimal.NewFromString(string(bytes))
			if err != nil {
				return decimal.Zero
			}
			return i
		}
	}
	return decimal.Zero
}

// go对象转map
func ObjectToMap(obj any) *map[string]interface{} {
	bytes, err := json.Marshal(obj)
	if err != nil {
		logs.Error("ObjectToMap:", err)
		return nil
	}
	data := map[string]interface{}{}
	json.Unmarshal(bytes, &data)
	return &data
}
