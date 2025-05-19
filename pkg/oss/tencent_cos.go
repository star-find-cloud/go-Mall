package oss

import (
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

var _cosClient *cos.Client

func init() {
	c := conf.GetConfig()
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.OSS.TencentCOS.Bucket, c.OSS.TencentCOS.Region))
	baseURL := &cos.BaseURL{BucketURL: u}

	// 创建 腾讯云 cos 客户端
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 60 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.OSS.TencentCOS.SecretId,
			SecretKey: c.OSS.TencentCOS.SecretKey,
		},
	})
	_cosClient = client
}

func GetCos() *cos.Client {
	return _cosClient
}
