package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/star-find-cloud/star-mall/internal/repo"
	"github.com/star-find-cloud/star-mall/model"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	proto "github.com/star-find-cloud/star-mall/protobuf/pb"
	"github.com/star-find-cloud/star-mall/utils"
	"github.com/tencentyun/cos-go-sdk-v5"
	"sync"
	"sync/atomic"
)

type ImageService interface {
	// 上传单张图片
	UploadImage(ctx context.Context, req *proto.ImageProto) (string, error)

	// 批量上传图片
	UploadImages(ctx context.Context, reqs []*proto.ImageProto) ([]string, error)

	// 下载图片
	DownloadImage(ctx context.Context, req *proto.ImageRequest) (*proto.ImageProto, error)

	// 删除图片
	DeleteImage(ctx context.Context, req *proto.ImageRequest) error
}

type OSSService struct {
	oosClient *cos.Client
	repo      repo.ImageRepository
}

func NewOSSService(oosClient *cos.Client, repo repo.ImageRepository) *OSSService {
	return &OSSService{
		oosClient: oosClient,
		repo:      repo,
	}
}

func (s *OSSService) UploadImage(ctx context.Context, req *proto.ImageProto) (string, error) {
	fileName := fmt.Sprintf("%d-%d-%d", req.OwnerType, req.OwnerId, utils.GetTimeNow())

	reader := bytes.NewReader(req.Data)
	_, err := s.oosClient.Object.Put(ctx, fileName, reader, nil)
	if err != nil {
		logger.AppLogger.Errorf("upload image err: %v", err)
		return "", err
	}

	image := &model.Image{
		ImageID:      req.ImageId,
		OwnerType:    req.OwnerType,
		OwnerID:      req.OwnerId,
		OssPath:      fileName,
		SHA256Hash:   req.Sha256Hash,
		IsCompressed: req.IsCompressed,
	}

	// 获取返回值, OssPath, err, 通过 err 是否为nil判断mysql OssPath 传递成功与否
	return s.repo.Create(ctx, image)
}

func (s *OSSService) UploadImages(ctx context.Context, reqs []*proto.ImageProto) ([]string, int, error) {
	var (
		wg         sync.WaitGroup
		successCnt int64 // 原子计数器, 保证并发安全性
		urlChan    = make(chan string, len(reqs))
		errChan    = make(chan error, len(reqs))
		sem        = make(chan struct{}, 20)
	)

	for i, req := range reqs {
		sem <- struct{}{}
		wg.Add(1)
		go func(idx int, req *proto.ImageProto) {
			defer func() {
				<-sem
				wg.Done()
			}()

			select {
			case <-ctx.Done():
				return
			default:
				resp, err := s.UploadImage(ctx, req)
				if err != nil {
					errChan <- fmt.Errorf("upload image err: %v", err)
					return
				}
				atomic.AddInt64(&successCnt, 1)
				urlChan <- resp
			}
		}(i, req)
	}

	go func() {
		defer close(urlChan)
		defer close(errChan)
		wg.Wait()
	}()

	urls := make([]string, 0, len(reqs))
	errs := make([]error, 0)
	for url := range urlChan {
		urls = append(urls, url)
	}

	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return urls, int(successCnt), fmt.Errorf("upload images err: %v", errs)
	}

	return urls, int(successCnt), nil
}

func (s *OSSService) DownloadImage(ctx context.Context, req *proto.ImageRequest) (*proto.ImageProto, error) {

	return nil, nil
}

func (s *OSSService) DeleteImage(ctx context.Context, req *proto.ImageRequest) error {
	return nil
}
