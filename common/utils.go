package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

const charset = "1234567890"

// 时间戳转日期
func TimestampToDate(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 获取当前时间
func GetUnixTime() int64 {
	fmt.Println(time.Now().Unix())
	return time.Now().Unix()
}

// 获取时间戳的 Nano 时间
func GetUnixNanoTime() int64 {
	return time.Now().UnixNano()
}

func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// Md5 加密
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return string(hex.EncodeToString(m.Sum(nil)))
}

// 验证邮箱
func VerifyEmail(email string) bool {
	pattern := `^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 获取日期
func FormatDate() string {
	return time.Now().Format("20060102")
}

// 生成订单号
var sequence int64 = 0 // 原子计数器
func GenerateOrderID() string {
	now := time.Now()

	// 时间基础部分（精确到秒）
	baseTime := now.Format("20060102150405")

	// 毫秒部分（固定3位）
	millis := fmt.Sprintf("%03d", now.Nanosecond()/1e6)

	// 进程ID（固定3位）
	pid := fmt.Sprintf("%03d", os.Getpid()%1000)

	// 原子递增序列（固定3位）
	seq := atomic.AddInt64(&sequence, 1) % 1000
	sequencePart := fmt.Sprintf("%03d", seq)

	return fmt.Sprintf("%s%s%s%s", baseTime, millis, pid, sequencePart)
}

func Mul(price float64, num int) float64 {
	return price * float64(num)
}
