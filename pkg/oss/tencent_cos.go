package oss

import (
	"context"
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

type TencentCos struct {
	client *cos.Client
	bucket string
}

var tencentCos = &TencentCos{}

func init() {
	c := conf.GetConfig()
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.OSS.TencentCOS.Bucket, c.OSS.TencentCOS.Region))
	baseURL := &cos.BaseURL{BucketURL: u}

	// 创建 腾讯云 cos 客户端
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 60 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.OSS.OSSAccessKeyID,
			SecretKey: c.OSS.OSSAccessKeySecret,
		},
	})
	tencentCos.client = client
	tencentCos.bucket = c.OSS.TencentCOS.Bucket
}

func NewTencentCos() *TencentCos {
	return tencentCos
}

func (oss *TencentCos) GetTencentCosClient() *cos.Client {
	return oss.client
}

// GeneratePresignedUploadURL 生成上传预签名URL, filePath 为文件永久路径，如：/test/123.jpg, opt 为请求头数组
func (oss *TencentCos) GeneratePresignedUploadURL(ctx context.Context, filePath, ContentType string) (string, error) {
	var c = conf.GetConfig()
	expiry := 30 * time.Minute

	var opt = &cos.PresignedURLOptions{
		Header: &http.Header{},
	}
	opt.Header.Set("Content-Type", ContentType)

	presignedURL, err := oss.client.Object.GetPresignedURL(
		ctx,
		http.MethodPut,
		filePath,
		c.OSS.OSSAccessKeyID,
		c.OSS.OSSAccessKeySecret,
		expiry,
		opt,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// GeneratePresignedDownloadURL 生成下载预签名URL, filePath 为文件永久路径，如：/test/123.jpg, opt 为请求头数组
func (oss *TencentCos) GeneratePresignedDownloadURL(ctx context.Context, filePath, ContentType string) (string, error) {
	var c = conf.GetConfig()
	expiry := 1 * time.Hour

	opt := &cos.PresignedURLOptions{
		Header: &http.Header{},
	}
	opt.Header.Set("Content-Type", ContentType)

	presignedURL, err := oss.client.Object.GetPresignedURL(ctx, http.MethodGet, filePath,
		c.OSS.OSSAccessKeyID,
		c.OSS.OSSAccessKeySecret,
		expiry,
		opt)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (oss *TencentCos) GeneratePresignedDeleteURL(ctx context.Context, filePath, ContentType string) (string, error) {
	var c = conf.GetConfig()
	expiry := 5 * time.Minute

	opt := &cos.PresignedURLOptions{
		Header: &http.Header{},
	}
	opt.Header.Set("Content-Type", ContentType)

	presignedURL, err := oss.client.Object.GetPresignedURL(ctx, http.MethodDelete, filePath,
		c.OSS.OSSAccessKeyID,
		c.OSS.OSSAccessKeySecret,
		expiry,
		opt)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
