package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/conf"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GetHostName 获取主机名
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

// ParsePathParamInt64 解析 url 中的整数参数为 int64
func ParsePathParamInt64(c *gin.Context, paramName string) (int64, error) {
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
	return int64(id), nil
}

// ExtractInt 匹配第一个连续的整数
func ExtractInt(s string) (int, error) {
	// 匹配第一个连续的整数（支持负号）
	re := regexp.MustCompile(`-?\d+`)
	match := re.FindString(s)
	if match == "" {
		return 0, errors.New("no integer found")
	}
	return strconv.Atoi(match)
}

// ToBase62 将 num 转换为 base62, 降低id传输长度
func ToBase62(num int64) (string, error) {
	if num <= 10000000 {
		return "", fmt.Errorf("num 不存在, 且")
	}

	var result []byte
	for num > 0 {
		var remainder = num % 62
		result = append(result, charset[remainder])
		num /= 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result), nil
}

// FromBase62 将 base62 转换为 十进制
func FromBase62(base62 string) (int64, error) {
	var result int64
	for i := 0; i < len(base62); i++ {
		var index = strings.IndexByte(charset, base62[i])
		if index < 0 {
			return 0, fmt.Errorf("invalid base62 character: %c", base62[i])
		}
		result = result*62 + int64(index)
	}
	return result, nil
}

func IsEnableOSS() bool {
	var c = conf.GetConfig()
	return c.OSS.Enable
}
