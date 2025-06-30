package domain

// 用户短信
type UserSms struct {
	ID        int
	IP        string
	Phone     string
	SendCount int
	AddDay    int
	AddTime   int
	Sign      string // 签名, 验证消息是否被篡改
}
