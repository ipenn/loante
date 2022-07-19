package tools

import "time"

func GetUnixTime() int64 {
	return time.Now().Unix()
}

func GetFormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func UnixTimeToStr(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func StrToUnixTime(str string) int64 {
	tm2, _ := time.Parse("2006-01-02 15:04:05", str)
	return tm2.Unix() - 8*3600
}

func ToDay(t int64) string { //t 小时
	return UnixTimeToStr(GetUnixTime() + t*3600)
}

func ToAddMonth(number int) string {
	nowTime := time.Now()
	var getTime time.Time
	getTime = nowTime.AddDate(0, number, 0)
	return getTime.Format("2006-01-02 15:04:05")
}

func ToAddDay(number int) string {
	nowTime := time.Now()
	var getTime time.Time
	getTime = nowTime.AddDate(0, 0, number)
	return getTime.Format("2006-01-02 15:04:05")
}
