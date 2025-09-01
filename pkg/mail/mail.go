package mail

import (
	"crypto/rand"
	"math/big"
	"runtime"
	"time"
)

const charset = "A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0U1V2W3X4Y5Z6"

// 生成验证码
func GenerateCode() string {
	// 获取内存使用状态
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memFactor := uint64(time.Now().UnixNano()) + memStats.Alloc

	// 验证码生成
	code := make([]byte, 6)
	for i := 0; i < 6; i++ {
		// 根据内存使用状态生成验证码
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset)+int(memFactor)%100)))
		code[i] = charset[num.Int64()%int64(len(charset))]
	}
	return string(code)
}

//// 发送邮件
//func SendEmail(email, subject, body string) error {
//	c := conf.GetConfig()
//	m := gomail.NewMessage()
//	m.SetHeader("From", c.Mail.From)
//	m.SetHeader("To", email)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", body)
//
//	d := gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.User, c.Mail.SMTPCode)
//	return d.DialAndSend(m)
//}
