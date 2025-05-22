package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/star-find-cloud/star-mall/model"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	proto "github.com/star-find-cloud/star-mall/protobuf/pb"
	"github.com/star-find-cloud/star-mall/repo"
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

	// 更改图片
	UpdateImage(ctx context.Context, req *proto.ImageProto) error
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

func (s *OSSService) UploadImage(ctx context.Context, req *proto.ImageProto) (string, int64, error) {
	fileName := fmt.Sprintf("%d-%d-%d", req.OwnerType, req.OwnerId, utils.GetTimeNow())

	reader := bytes.NewReader(req.Data)
	_, err := s.oosClient.Object.Put(ctx, fileName, reader, nil)
	if err != nil {
		logger.AppLogger.Errorf("upload image err: %v", err)
		return "", 0, err
	}

	id, err := utils.GenerateUid()
	if err != nil {
		logger.AppLogger.Errorf("generate uid err: %v", err)
	}
	image := &model.Image{
		ImageID:      id,
		OwnerType:    req.OwnerType,
		OwnerID:      req.OwnerId,
		OssPath:      fileName,
		SHA256Hash:   req.Sha256Hash,
		IsCompressed: req.IsCompressed,
	}

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
				resp, _, err := s.UploadImage(ctx, req)
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
	// image_id 获取图片信息
	image, err := s.repo.GetByID(ctx, int(req.ImageId))
	if err != nil {
		logger.AppLogger.Errorf("failed to get image by ID %d: %v", req.ImageId, err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// 从 OSS 下载图片数据
	resp, err := s.oosClient.Object.Get(ctx, image.OssPath, nil)
	if err != nil {
		logger.AppLogger.Errorf("failed to download image from OSS (path: %s): %v", image.OssPath, err)
		return nil, fmt.Errorf("failed to download image from OSS: %w", err)
	}
	defer resp.Body.Close()

	// 读取图片内容到字节数组
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		logger.AppLogger.Errorf("failed to read image data (path: %s): %v", image.OssPath, err)
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 4. 构造 ImageProto 返回值
	imageProto := &proto.ImageProto{
		ImageId:   int64(image.ImageID),
		OwnerType: image.OwnerType,
		Data:      buffer.Bytes(),
	}

	return imageProto, nil
}

func (s *OSSService) DeleteImage(ctx context.Context, id int64) error {
	image, err := s.repo.GetByID(ctx, int(id))
	if err != nil {
		logger.AppLogger.Errorf("failed to get image by ID %d: %v", id, err)
		return fmt.Errorf("failed to get image: %w", err)
	}

	resp, err := s.oosClient.Object.Delete(ctx, image.OssPath)
	if err != nil {
		logger.AppLogger.Errorf("failed to delete image from OSS (path: %s): %v", image.OssPath, err)
		return fmt.Errorf("failed to delete image from OSS: %w", err)
	}
	logger.AppLogger.Infof("delete image from OSS (path: %s): %v", image.OssPath, resp)
	defer resp.Body.Close()

	err = s.repo.Delete(ctx, int(id))
	if err != nil {
		logger.AppLogger.Errorf("failed to delete image from OSS (path: %s): %v", image.OssPath, err)
		return fmt.Errorf("failed to delete image from OSS: %w", err)
	}
	return nil
}

func (s *OSSService) UpdateImage(ctx context.Context, req *proto.ImageProto) error {
	imageID := req.ImageId

	err := s.DeleteImage(ctx, imageID)
	if err != nil {
		return err
	}

	_, _, err = s.UploadImage(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
