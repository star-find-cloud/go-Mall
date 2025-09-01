package service

import (
	"context"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	"github.com/star-find-cloud/star-mall/repo"
	"testing"
)

func TestImageForOSS_GenerateUploadURLs(t *testing.T) {
	imagea := &domain.Image{
		OwnerID:      123,
		OwnerType:    101,
		SHA256Hash:   "ab1234@",
		IsCompressed: false,
		ContentType:  _const.WEBP,
	}

	imageb := &domain.Image{
		OwnerID:      352,
		OwnerType:    101,
		SHA256Hash:   "A15Wc@",
		IsCompressed: false,
		ContentType:  _const.PNG,
	}

	images := []*domain.Image{imagea, imageb}
	cos := oss.NewTencentCos()
	db, err := database.NewMySQL()
	cache, err := database.NewRedis()
	imageRepo := repo.NewImageRepo(db)
	imageService := NewImageForOssService(imageRepo, cos, cache)
	url, err := imageService.GenerateUploadURLs(context.Background(), images)
	fmt.Println(url, err)
}
