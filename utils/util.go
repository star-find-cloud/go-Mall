package utils

import "os"

// 获取主机名
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}
