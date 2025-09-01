package conf

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	if c == nil {
		t.Fatal("配置对象为 nil，初始化失败")
	}

	// 打印完整结构体（调试用）
	fmt.Println(c.OSS.OSSAccessKeyID)
	fmt.Println(c.OSS.OSSAccessKeySecret)
}
