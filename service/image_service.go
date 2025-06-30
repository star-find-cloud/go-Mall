package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/domain"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/utils"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageService interface {
	// Save 上传单张图片
	Save(ctx context.Context, image *domain.Image, file *multipart.FileHeader, c *gin.Context) (int64, string, error)

	// SaveMore 批量上传图片
	//SaveMore(ctx context.Context, image []*domain.Image) ([]int64, int64, error)

	// GetImage 下载图片
	GetImage(ctx context.Context, id int64) (*domain.Image, error)

	//// 删除图片
	//DeleteImage(ctx context.Context, req *proto.ImageRequest) error

	//// 更改图片
	//UpdateImage(ctx context.Context, req *proto.ImageProto) error
}

type Image struct {
	repo repo.ImageRepo
}

func NewImageService(repo repo.ImageRepo) Image {
	return Image{repo: repo}
}

func (s Image) Save(ctx context.Context, image *domain.Image, file *multipart.FileHeader, c *gin.Context) (int64, string, error) {
	uploadDIr := "/var/star-mall/user/images/"
	if err := os.MkdirAll(uploadDIr, 0755); err != nil {
		applog.AppLogger.Error("create upload dir err: ", err)
		return 0, "", err
	}

	uid, err := utils.GenerateUid()
	if err != nil {
		applog.AppLogger.Error("generate uid err: ", err)
		return 0, "", err
	}
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%d_%d_%s", uid, fileExt)
	filePath := filepath.Join(uploadDIr, fileName)

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		applog.AppLogger.Error("save image err: ", err)
		return 0, "", err
	}

	image.FilePath = filePath
	image.ContentType = file.Header.Get("Content-Type")
	id, err := s.repo.UploadImage(ctx, image)
	if err != nil {
		if err = os.Remove(filePath); err != nil {
			applog.AppLogger.Error("remove image err: ", err)
			return 0, "", err
		}
		applog.AppLogger.Error("upload image err: ", err)
		return 0, "", err
	}

	return id, filePath, nil
}

//func (s Image) SaveMore(ctx context.Context, images []*domain.Image) ([]int64, int64, error) {
//	var (
//		wg         sync.WaitGroup
//		successCnt int64 // 原子计数器, 保证并发安全性
//		IDChan     = make(chan int64, len(images))
//		errChan    = make(chan error, len(images))
//		sem        = make(chan struct{}, 20)
//	)
//
//	for i, image := range images {
//		sem <- struct{}{}
//		wg.Add(1)
//		go func(idx int, images *domain.Image) {
//			defer func() {
//				<-sem
//				wg.Done()
//			}()
//
//			select {
//			case <-ctx.Done():
//				return
//			default:
//				id, err := s.Save(ctx, image, )
//				if err != nil {
//					errChan <- fmt.Errorf("upload image err: %v", err)
//					return
//				}
//				atomic.AddInt64(&successCnt, 1)
//				IDChan <- id
//			}
//		}(i, image)
//	}
//
//	go func() {
//		defer close(IDChan)
//		defer close(errChan)
//		wg.Wait()
//	}()
//
//	ids := make([]int64, 0, len(images))
//	errs := make([]error, 0)
//	for id := range IDChan {
//		ids = append(ids, id)
//	}
//
//	for err := range errChan {
//		errs = append(errs, err)
//	}
//
//	if len(errs) > 0 {
//		return ids, successCnt, fmt.Errorf("upload images err: %v", errs)
//	}
//
//	return ids, successCnt, nil
//}

func (s Image) GetImage(ctx context.Context, id int64) (*domain.Image, error) {
	if id <= 0 {
		return nil, fmt.Errorf("image id must be greater than 0")
	}
	return s.repo.GetByID(ctx, id)
}

//type OSSService struct {
//	ossClient *cos.Client
//	repo      repo.ImageRepo
//}
//
//func NewOSSService(ossClient *cos.Client, repo repo.ImageRepo) *OSSService {
//	return &OSSService{
//		ossClient: ossClient,
//		repo:      repo,
//	}
//}

//// UploadImage 返回 图片上传的 ossPath(对象存储路径), imageID(图片编号)
//func (s *OSSService) UploadImage(ctx context.Context, req *proto.ImageProto) (string, int64, error) {
//	// 根据图片所属者类型和所属者ID生成图片名
//	fileName := fmt.Sprintf("%d-%d-%d", req.OwnerType, req.OwnerId, utils.GetTimeNow())
//
//	// 从 grpc 中获取图片数据
//	reader := bytes.NewReader(req.Data)
//	_, err := s.ossClient.Object.Put(ctx, fileName, reader, nil)
//	if err != nil {
//		logger.AppLogger.Errorf("upload image err: %v", err)
//		return "", 0, err
//	}
//
//	// 生成图片唯一编号
//	id, err := utils.GenerateUid()
//	if err != nil {
//		logger.AppLogger.Errorf("generate uid err: %v", err)
//	}
//	image := &domain.Image{
//		ImageID:      id,
//		OwnerType:    req.OwnerType,
//		OwnerID:      req.OwnerId,
//		OssPath:      fileName,
//		SHA256Hash:   req.Sha256Hash,
//		IsCompressed: req.IsCompressed,
//	}
//
//	// 保存图片信息
//	return s.repo.Save(ctx, image)
//}

//func (s *OSSService) SaveMore(ctx context.Context, reqs []*proto.ImageProto) ([]int64, int64, error) {
//	var (
//		wg         sync.WaitGroup
//		successCnt int64 // 原子计数器, 保证并发安全性
//		//urlChan    = make(chan string, len(reqs))
//		IDChan  = make(chan int64, len(reqs))
//		errChan = make(chan error, len(reqs))
//		sem     = make(chan struct{}, 20)
//	)
//
//	for i, req := range reqs {
//		sem <- struct{}{}
//		wg.Add(1)
//		go func(idx int, req *proto.ImageProto) {
//			defer func() {
//				<-sem
//				wg.Done()
//			}()
//
//			select {
//			case <-ctx.Done():
//				return
//			default:
//				_, id, err := s.UploadImage(ctx, req)
//				if err != nil {
//					errChan <- fmt.Errorf("upload image err: %v", err)
//					return
//				}
//				atomic.AddInt64(&successCnt, 1)
//				IDChan <- id
//			}
//		}(i, req)
//	}
//
//	go func() {
//		defer close(IDChan)
//		defer close(errChan)
//		wg.Wait()
//	}()
//
//	ids := make([]int64, 0, len(reqs))
//	errs := make([]error, 0)
//	for id := range IDChan {
//		ids = append(ids, id)
//	}
//
//	for err := range errChan {
//		errs = append(errs, err)
//	}
//
//	if len(errs) > 0 {
//		return ids, successCnt, fmt.Errorf("upload images err: %v", errs)
//	}
//
//	return ids, successCnt, nil
//}
//
//func (s *OSSService) DownloadImage(ctx context.Context, req *proto.ImageRequest) (*proto.ImageProto, error) {
//	// image_id 获取图片信息
//	image, err := s.repo.GetByID(ctx, int(req.ImageId))
//	if err != nil {
//		logger.AppLogger.Errorf("failed to get image by ID %d: %v", req.ImageId, err)
//		return nil, fmt.Errorf("failed to get image: %w", err)
//	}
//
//	// 从 OSS 下载图片数据
//	resp, err := s.ossClient.Object.Get(ctx, image.OssPath, nil)
//	if err != nil {
//		logger.AppLogger.Errorf("failed to download image from OSS (path: %s): %v", image.OssPath, err)
//		return nil, fmt.Errorf("failed to download image from OSS: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// 读取图片内容到字节数组
//	var buffer bytes.Buffer
//	_, err = buffer.ReadFrom(resp.Body)
//	if err != nil {
//		logger.AppLogger.Errorf("failed to read image data (path: %s): %v", image.OssPath, err)
//		return nil, fmt.Errorf("failed to read image data: %w", err)
//	}
//
//	// 4. 构造 ImageProto 返回值
//	imageProto := &proto.ImageProto{
//		ImageId:   image.ImageID,
//		OwnerType: image.OwnerType,
//		Data:      buffer.Bytes(),
//	}
//
//	return imageProto, nil
//}
//
//func (s *OSSService) DeleteImage(ctx context.Context, id int64) error {
//	image, err := s.repo.GetByID(ctx, int(id))
//	if err != nil {
//		logger.AppLogger.Errorf("failed to get image by ID %d: %v", id, err)
//		return fmt.Errorf("failed to get image: %w", err)
//	}
//
//	resp, err := s.ossClient.Object.Delete(ctx, image.OssPath)
//	if err != nil {
//		logger.AppLogger.Errorf("failed to delete image from OSS (path: %s): %v", image.OssPath, err)
//		return fmt.Errorf("failed to delete image from OSS: %w", err)
//	}
//	logger.AppLogger.Infof("delete image from OSS (path: %s): %v", image.OssPath, resp)
//	defer resp.Body.Close()
//
//	err = s.repo.Delete(ctx, int(id))
//	if err != nil {
//		logger.AppLogger.Errorf("failed to delete image from OSS (path: %s): %v", image.OssPath, err)
//		return fmt.Errorf("failed to delete image from OSS: %w", err)
//	}
//	return nil
//}
//
//func (s *OSSService) UpdateImage(ctx context.Context, req *proto.ImageProto) error {
//	imageID := req.ImageId
//
//	err := s.DeleteImage(ctx, imageID)
//	if err != nil {
//		return err
//	}
//
//	_, _, err = s.UploadImage(ctx, req)
//	if err != nil {
//		return err
//	}
//	return nil
//}
