package model

// 设置
type Setting struct {
	ID              int
	SiteTitle       string // 网站标题
	SiteLogo        string // 网站logo
	SiteKeywords    string
	SiteDescription string // 网站描述
	NoPicture       string
	SiteIcp         string // 网站备案号
	SiteTelephone   string // 网站联系电话
	SearchKeywords  string
	AppID           string // 第三方服务 api 凭证ID
	AppSecret       string
	EndPoint        string // 对象存储服务节点
	BucketName      string // 存储桶名称
	OssStatus       int    // 对象存储服务状态
	StatsCode       string // 统计代码
}
