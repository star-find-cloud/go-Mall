package utils

import "time"

// FormatDate 获取日期
func FormatDate() string {
	return time.Now().Format("20060102")
}

// GetTimeNow 获取当前时间戳
func GetTimeNow() int64 {
	// 构造起始时间
	startTime := time.Date(2008, 10, 1, 12, 31, 0, 0, time.UTC)

	// 获取当前时间
	currentTime := time.Now().UTC() // 保持时区一致

	// 计算时间差并转换为纳秒
	duration := currentTime.Sub(startTime)
	nanoseconds := duration.Nanoseconds()

	return nanoseconds
}

// TimestampToDate 时间戳转日期
func TimestampToDate(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}
