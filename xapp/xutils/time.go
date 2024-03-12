package xutils

import (
	"strings"
	"time"
)

const (
	TimeLayout   = "2006-01-02 15:04:05"
	DateLayout   = "2006-01-02"
	MounthLayout = "200601"
)

// 时间戳转本地日期 eg: 1609459200 -> 2021-01-01
func TimeStampToLocalDate(tvalue int64) string {
	if tvalue == 0 {
		return ""
	}
	tm := time.Unix(tvalue, 0)
	tstr := tm.Format("2006-01-02")
	return strings.Split(tstr, " ")[0]
}

// 本地日期转时间戳(秒) eg: 2021-01-01 -> 1609459200
func LocalDateToTimeStamp(timestr string) int64 {
	t, _ := time.ParseInLocation("2006-01-02", timestr, time.Local)
	return t.Local().Unix()
}

// 时间戳(秒)转本地时间 eg: 1609459200 -> 2021-01-01 08:00:00
func TimeStampToLocalTime(tvalue int64) string {
	if tvalue == 0 {
		return ""
	}
	tm := time.Unix(tvalue, 0)
	tstr := tm.Format("2006-01-02 15:04:05")
	return tstr
}

// 本地时间转时间戳(秒) eg: 2021-01-01 08:00:00 -> 1609459200
func LocalTimeToTimeStamp(timestr string) int64 {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timestr, time.Local)
	return t.Local().Unix()
}

// 本地时间转utc时间 eg: 2021-01-01 08:00:00 -> 2021-01-01T00:00:00Z
func LocalTimeToUtcFormat(timestr string) string {
	if len(timestr) == 0 {
		return timestr
	}
	if len(timestr) == 10 {
		timestr = timestr + " 00:00:00"
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timestr, time.Local)
	r := t.UTC().Format("2006-01-02T15:04:05Z")
	return r
}

// utc时间转本地时间 eg: 2021-01-01T00:00:00Z -> 2021-01-01 08:00:00
func UtcToLocalTime(timestr string) string {
	if len(timestr) == 0 {
		return ""
	}
	t, err := time.Parse(time.RFC3339, timestr)
	if err != nil {
		return ""
	}
	localTime := t.Local()
	return localTime.In(time.Local).Format("2006-01-02 15:04:05")
}

// 获取本地时间 eg: 2021-01-01 00:00:00
func GetLocalTime() string {
	tm := time.Now()
	return tm.In(time.Local).Format("2006-01-02 15:04:05")
}

// 获取本地日期 eg: 2021-01-01
func GetLocalDate() string {
	tm := time.Now()
	return tm.In(time.Local).Format("2006-01-02")
}

func Now() string {
	return GetLocalTime()
}

func UtcOffset() int {
	currentTime := time.Now()
	_, offset := currentTime.Zone()
	utcTime := currentTime.UTC()
	_, utcOffset := utcTime.Zone()
	return int((time.Duration(offset-utcOffset) * time.Second).Hours())
}

func UtcNow() string {
	tm := time.Now()
	return tm.UTC().Format("2006-01-02 15:04:05")
}

func GetUtcDate() string {
	tm := time.Now()
	return tm.UTC().Format("2006-01-02")
}

func LocalTimeToUtcTime(timestr string) string {
	if len(timestr) == 0 {
		return timestr
	}
	if len(timestr) == 10 {
		timestr = timestr + " 00:00:00"
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timestr, time.Local)
	r := t.UTC().Format("2006-01-02 15:04:05")
	return r
}
