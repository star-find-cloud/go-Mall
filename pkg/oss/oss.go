package oss

import "context"

type OSS interface {
	// GeneratePresignedUploadURL 生成上传预签名URL
	GeneratePresignedUploadURL(ctx context.Context, filePath, ContentType string) (string, error)

	// GeneratePresignedDownloadURL 生成下载签名URL
	GeneratePresignedDownloadURL(ctx context.Context, filePath, ContentType string) (string, error)

	// GeneratePresignedDeleteURL 删除签名URL
	GeneratePresignedDeleteURL(ctx context.Context, filePath, ContentType string) (string, error)
}
