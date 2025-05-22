package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"os"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

// 获取主机名
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

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

// 解析 url 中的整数参数
func ParsePathParamInt(c *gin.Context, paramName string) (int, error) {
	// 从 user 中获取对应的参数
	strID := c.Param(paramName)
	if strID == "" {
		return 0, errors.New("missing path parameter: " + paramName)
	}
	// 将参数从 str 转换为 int
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ParsePathParamInt64(c *gin.Context, paramName string) (int64, error) {
	// 从 user 中获取对应的参数
	strID := c.Param(paramName)
	if strID == "" {
		return 0, errors.New("missing path parameter: " + paramName)
	}
	// 将参数从 str 转换为 int
	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func ExtractInt(s string) (int, error) {
	// 匹配第一个连续的整数（支持负号）
	re := regexp.MustCompile(`-?\d+`)
	match := re.FindString(s)
	if match == "" {
		return 0, errors.New("no integer found")
	}
	return strconv.Atoi(match)
}

// 生成 uid
func GenerateUid() (int64, error) {
	nodeID, _ := ExtractInt(GetHostName())
	node, err := snowflake.NewNode(int64(nodeID))
	if err != nil {
		return 0, err
	}

	return node.Generate().Int64(), nil
}

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
