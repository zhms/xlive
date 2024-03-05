package xredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"xcom/global"
	"xcom/xutils"

	"github.com/beego/beego/logs"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type XRedisSubCallback func(string)
type XRedis struct {
	client             *redis.Client
	host               string
	port               int
	db                 int
	password           string
	recving            bool
	subscribecallbacks sync.Map
	project            string
}

func (this *XRedis) Init(cfgname string) {
	if this.client != nil {
		return
	}
	this.project = global.Project
	host := viper.GetString(fmt.Sprint(cfgname, ".host"))
	port := viper.GetInt(fmt.Sprint(cfgname, ".port"))
	db := viper.GetInt(fmt.Sprint(cfgname, ".db"))
	password := viper.GetString(fmt.Sprint(cfgname, ".password"))

	maxidle := 5
	if strings.Contains(global.Env, "prd") {
		maxidle = 20
	}
	m := viper.GetInt(fmt.Sprint(cfgname, ".maxidle"))
	if m > 0 {
		maxidle = m
	}

	this.client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%v:%v", host, port),
		Password:     password,
		DB:           db,
		MaxIdleConns: maxidle,
	})
	logs.Debug("连接redis 成功:", host, port, db, cfgname)
}

func (this *XRedis) Close() {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
}

func (this *XRedis) Client() *redis.Client {
	return this.client
}

func (this *XRedis) GetCacheMap(key string, cb func() (*xutils.XMap, error)) (*xutils.XMap, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheMap error:", key, err.Error())
		return nil, err
	}
	if data != nil {
		jdata := map[string]interface{}{}
		json.Unmarshal(data, &jdata)
		xmap := xutils.XMap{}
		xmap.RawData = jdata
		return &xmap, nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheMaps(key string, cb func() (*xutils.XMaps, error)) (*xutils.XMaps, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil {
		logs.Error("GetCacheMaps error:", key, err.Error())
		return nil, err
	}
	if data != nil {
		jdata := []map[string]interface{}{}
		json.Unmarshal(data, &jdata)
		xmaps := xutils.XMaps{RawData: jdata}
		return &xmaps, nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheString(key string, cb func() (string, error)) (string, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheString error:", key, err.Error())
		return "", err
	}
	if data != nil {
		return string(data), nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheStrings(key string, cb func() ([]string, error)) ([]string, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheStrings error:", key, err.Error())
		return nil, err
	}
	if data != nil {
		jdata := []string{}
		json.Unmarshal(data, &jdata)
		return jdata, nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheInt(key string, cb func() (int64, error)) (int64, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheInt error:", key, err.Error())
		return 0, err
	}
	if data != nil {
		return xutils.ToInt64(data), nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheInts(key string, cb func() ([]int64, error)) ([]int64, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheInts error:", key, err.Error())
		return nil, err
	}
	if data != nil {
		jdata := []int64{}
		json.Unmarshal(data, &jdata)
		return jdata, nil
	} else {
		return cb()
	}
}

func (this *XRedis) GetCacheDecimal(key string, cb func() (decimal.Decimal, error)) (decimal.Decimal, error) {
	data, err := this.client.Get(context.Background(), key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetCacheDecimal error:", key, err.Error())
		return decimal.Zero, err
	}
	if data != nil {
		d, e := decimal.NewFromString(string(data))
		if e != nil {
			return decimal.Zero, e
		}
		return d, nil
	} else {
		return cb()
	}
}

func (this *XRedis) Lock(key string, expire_second int) bool {
	key = fmt.Sprintf("%v:lock:%v", this.project, key)
	if expire_second <= 0 {
		r, err := this.client.Do(context.Background(), "setnx", key, "1").Result()
		if err != nil {
			logs.Error(err.Error())
			return false
		}
		if r == nil {
			return false
		}
		ir := xutils.ToInt(r)
		return ir == 1
	} else {
		r, err := this.client.Do(context.Background(), "set", key, "1", "EX", expire_second, "NX").Result()
		if err != nil {
			return false
		}
		if r == nil {
			return false
		}
		ir := xutils.ToString(r)
		return ir == "OK"
	}
}

// redis解锁
func (this *XRedis) UnLock(key string) {
	key = fmt.Sprintf("%v:lock:%v", this.project, key)
	this.client.Del(context.Background(), key)
}
