package domain

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
	StatsCode       string // 统计代码
}
