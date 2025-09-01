package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// GenerateOrderID 生成订单号（返回 int64 类型）
func GenerateOrderID() (int64, error) {
	var sequence int64 = 0 // 原子计数器
	now := time.Now()

	// 时间基础部分（精确到秒），格式化为数字字符串
	baseTime := now.Format("20060102150405")

	// 毫秒部分（固定3位）
	millis := fmt.Sprintf("%03d", now.Nanosecond()/1e6)

	// 进程ID（固定3位）
	pid := fmt.Sprintf("%03d", os.Getpid()%1000)

	// 原子递增序列（固定3位）
	seq := atomic.AddInt64(&sequence, 1) % 1000
	sequencePart := fmt.Sprintf("%03d", seq)

	// 拼接订单号
	orderIDStr := fmt.Sprintf("%s%s%s%s", baseTime, millis, pid, sequencePart)

	// 转换为 int64 类型
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("failed to generate order ID as int64")
	}

	return orderID, nil
}

// GetSumPrice 获取总价
func GetSumPrice(price float64, num int) float64 {
	return price * float64(num)
}

// ParsePathParamInt 解析 url 中的整数参数
func ParsePathParamInt(c *gin.Context, paramName string) (int, error) {
	// 从 url 中获取对应的参数
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
