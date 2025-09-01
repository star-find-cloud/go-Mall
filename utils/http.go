package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net"
	"net/http/httputil"
	"os"
	"strings"
)

// 判断是否发生broken pipe 错误
func IsBrokenPipeErr(err interface{}) bool {
	if network, ok := err.(*net.OpError); ok {
		if se, ok := network.Err.(*os.SyscallError); ok {
			lowerErr := strings.ToLower(se.Error())
			return strings.Contains(lowerErr, "broken pipe") || strings.Contains(lowerErr, "connection reset by peer")
		}
	}
	return false
}

// 捕获并格式化请求
func DumpRequest(c *gin.Context) string {
	reqBytes, err := httputil.DumpRequest(c.Request, false)
	if err != nil {
		return "Failed to dump request: " + err.Error()
	}

	// 过滤请求头中的敏感信息
	headers := bytes.Split(reqBytes, []byte("\r\n"))
	for i, header := range headers {
		if bytes.HasPrefix(header, []byte("Authorization:")) {
			headers[i] = []byte("Authorization: ***")
		}
	}
	return string(bytes.Join(headers, []byte("\r\n")))
}
