package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/idgen"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	"github.com/star-find-cloud/star-mall/repo"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// ImageMetaDataService 图片元数据服务接口
type ImageMetaDataService interface {
	// Save 保存图片元数据
	Save(ctx context.Context, image *domain.Image) (int64, error)

	// SaveMore 批量保存图片元数据
	SaveMore(ctx context.Context, image []*domain.Image) (map[int]int64, error)

	// Get 获取图片元数据
	Get(ctx context.Context, id int64) (*domain.Image, error)

	// Delete 删除图片元数据
	Delete(ctx context.Context, id int64) error

	// Update 更改图片元数据
	Update(ctx context.Context, oldImageID int64, newImage *domain.Image) (int64, error)
}

// ImageLocalService 图片本地服务接口
type ImageLocalService interface {
	// Upload 上传图片到服务器
	Upload(ctx context.Context, image *domain.Image, file *multipart.FileHeader) (int64, error)

	// UploadMore 批量上传图片到服务器
	UploadMore(ctx context.Context, images []*domain.Image, files []*multipart.FileHeader) (map[int]int64, error)

	// Download 获取图片
	Download(ctx context.Context, id int64) (string, error)

	// Remove 删除图片
	Remove(ctx context.Context, path string) error
}

// PresignedService 预签名服务接口
type PresignedService interface {
	// GenerateUploadURL 生成上传预签名 返回 path 和 预签名
	GenerateUploadURL(ctx context.Context, image *domain.Image) (int64, string, error)

	// GenerateUploadURLs 批量生成上传预签名
	GenerateUploadURLs(ctx context.Context, images []*domain.Image) (map[int64]string, error)

	// GenerateDownloadURL 生成下载预签名
	GenerateDownloadURL(ctx context.Context, id int64) (string, error)

	// GenerateDeleteURL 生成删除预签名
	GenerateDeleteURL(ctx context.Context, id int64) (string, error)
}

type ImageForDB struct {
	repo  repo.ImageRepo
	cache *database.Redis
}

func NewImageForDBService(repo repo.ImageRepo) *ImageForDB {
	return &ImageForDB{repo: repo}
}

// Save 保存元数据
func (s *ImageForDB) Save(ctx context.Context, image *domain.Image) (int64, error) {
	id, err := s.repo.UploadImage(ctx, image)
	if err != nil {
		if err = os.Remove(image.Path); err != nil {
			log.AppLogger.Errorln("remove image err: ", err)
			return 0, err
		}
		log.AppLogger.Errorln("upload image err: ", err)
		return 0, err
	}

	return id, nil
}

func (s *ImageForDB) SaveMore(ctx context.Context, image []*domain.Image) (map[int]int64, error) {
	var (
		total      = len(image)
		resultChan = make(chan map[int]int64, total)
		errChan    = make(chan error, total)
		wg         = &sync.WaitGroup{}
	)
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(idx int, img *domain.Image) {
			defer wg.Done()

			id, err := s.Save(ctx, img)
			if err != nil {
				select {
				case errChan <- err:
				default:
					log.AppLogger.Warnf("save failed for image %d, error dropped due to full channel: %v", idx, err)
				}
				return
			}
			select {
			case resultChan <- map[int]int64{idx: id}:
			case <-ctx.Done():
				log.AppLogger.Warnf("Context done while sending result for image %d", idx)
			}
		}(i, image[i])
	}

	var (
		imgs     = make(map[int]int64)
		errSlice []error
	)

	for resultChan != nil || errChan != nil {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				for k, v := range result {
					imgs[k] = v
				}
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errSlice = append(errSlice, err)
			}
		case <-ctx.Done():
			log.AppLogger.Warnln("save more operation cancelled via context")
			return imgs, ctx.Err()
		}
	}

	if len(errSlice) > 0 {
		for _, err := range errSlice {
			log.AppLogger.Warnf("save image matedata error: %v", err)
		}
		return imgs, fmt.Errorf("%d out of %d uploads failed", len(errSlice), total)
	}

	return imgs, nil
}

func (s *ImageForDB) Get(ctx context.Context, id int64) (*domain.Image, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id is not exist \n")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ImageForDB) Update(ctx context.Context, oldImageID int64, newImage *domain.Image) (int64, error) {
	if oldImageID <= 0 || newImage == nil {
		return 0, fmt.Errorf("id is not exist \n")
	}
	return s.repo.Update(ctx, oldImageID, newImage)
}

func (s *ImageForDB) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is not exist \n")
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("delete image err: %v \n", err)
		return errors.New("delete failed")
	}

	return nil
}

// Upload 上传图片到服务器
func (s *ImageForDB) Upload(ctx context.Context, image *domain.Image, file *multipart.FileHeader) (int64, error) {
	var (
		c       = conf.GetConfig()
		saveDir = c.App.ImageDir
	)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	uid, err := idgen.GenerateUid()
	if err != nil {
		log.AppLogger.Errorln("get uid error: ", err)
		return 0, err
	}

	select {
	case <-ctx.Done():
		return 0, errors.New("busy servers")
	default:
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("image/%d/%d/%d.%s \n", image.OwnerType, image.OwnerID, uid, fileExt)
	filePath := filepath.Join(saveDir, fileName)

	// 访问图片
	src, err := file.Open()
	if err != nil {
		log.AppLogger.Warnln("open image error: ", err)
		return 0, fmt.Errorf("open image err: %v", err)
	}
	defer src.Close()

	// 创建图片地址
	dst, err := os.Create(filePath)
	if err != nil {
		log.AppLogger.Warnln("save image error:", err)
		return 0, fmt.Errorf("save image error: %v", err)
	}
	defer dst.Close()

	done := make(chan error, 1)
	go func() {
		// 将访问到的图片进行持久化
		_, err = io.Copy(dst, src)
		done <- err
	}()

	select {
	case <-ctx.Done():
		// 超时取消文件上传, 清理文件
		os.Remove(filePath)
		return 0, ctx.Err()
	case err := <-done:
		if err != nil {
			log.AppLogger.Warnln("copy image error:", err)
			return 0, fmt.Errorf("copy image error: %v", err)
		}
	}

	return uid, nil
}

// UploadMore 批量上传图片
func (s *ImageForDB) UploadMore(ctx context.Context, images []*domain.Image, files []*multipart.FileHeader) (map[int]int64, error) {
	// 判断图片数量是否一致
	if len(images) != len(files) {
		log.AppLogger.Warnln("The number of pictures and files is not equal")
		return nil, errors.New("the number of pictures and files is not equal")
	}

	var (
		total      = len(images)
		resultChan = make(chan map[int]int64, len(images))
		errChan    = make(chan error, len(images))
		wg         = &sync.WaitGroup{}
	)
	wg.Add(total)
	// 并发处理图片
	for i := 0; i < total; i++ {
		go func(idx int, img *domain.Image, f *multipart.FileHeader) {
			defer wg.Done()

			id, err := s.Upload(ctx, img, f)
			if err != nil {
				select {
				case errChan <- err:
				default:
					log.AppLogger.Warnf("Upload failed for image %d, error dropped due to full channel: %v", idx, err)
				}
				return
			}

			select {
			case resultChan <- map[int]int64{idx: id}:
			case <-ctx.Done():
				log.AppLogger.Warnf("Context done while sending result for image %d", idx)
			}
		}(i, images[i], files[i])
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	var (
		errSlice []error
		imgs     = make(map[int]int64)
	)

	for resultChan != nil || errChan != nil {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				for k, v := range result {
					imgs[k] = v
				}
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errSlice = append(errSlice, err)
			}
		case <-ctx.Done():
			log.AppLogger.Warnln("Upload more operation cancelled via context")
			return imgs, ctx.Err()
		}
	}

	if len(errSlice) > 0 {
		for _, err := range errSlice {
			log.AppLogger.Errorln("Upload more operation failed:", err)
		}
		return imgs, fmt.Errorf("%d out of %d uploads failed", len(errSlice), total)
	}

	return imgs, nil
}

func (s *ImageForDB) Download(ctx context.Context, id int64) (string, error) {
	path, err := s.repo.GetPathByImageID(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("get path by image id error: %v \n", err)
		return "", err
	}
	return path, nil
}

func (s *ImageForDB) Remove(ctx context.Context, path string) error {
	err := os.Remove(path)
	if err != nil {
		log.AppLogger.Errorln("remove image err: ", err)
		return err
	}
	return nil
}

type ImageForOSS struct {
	repo  repo.ImageRepo
	oss   oss.OSS
	cache *database.Redis
}

func NewImageForOssService(repo repo.ImageRepo, oss oss.OSS, cache *database.Redis) *ImageForOSS {
	return &ImageForOSS{repo: repo, oss: oss, cache: cache}
}

func (s *ImageForOSS) GenerateUploadURL(ctx context.Context, image *domain.Image) (int64, string, error) {
	// 验证文件后缀
	v, exist := _const.ContentTypeStringMap[image.ContentType]
	if !exist {
		log.AppLogger.Warnln("image content type is not exist")
		return 0, "", fmt.Errorf("image content type is not exist")
	}

	// 生成唯一文件名
	uid, err := idgen.GetUid(ctx, s.cache)
	if err != nil {
		log.AppLogger.Warnf("get uid error:%v", err)
		return 0, "", err
	}

	path := fmt.Sprintf("image/%d/%d/%d.%s", image.OwnerType, image.OwnerID, uid, v)

	url, err := s.oss.GeneratePresignedUploadURL(ctx, path, v)
	if err != nil {
		log.AppLogger.Warnf("生成上传预签名失败: %v", err)
		return 0, "", err
	}

	cacheMap := map[string]interface{}{
		"imageID":     uid,
		"ownerID":     image.OwnerID,
		"ownerType":   image.OwnerType,
		"path":        path,
		"sha256":      image.SHA256Hash,
		"contentType": image.ContentType,
	}
	err = idgen.MakeUidCache(ctx, uid, s.cache, cacheMap)
	if err != nil {
		log.AppLogger.Warnln("make uid cache error")
		log.RedisLogger.Warnf("make uid cache error: %v", err)
		return 0, "", err
	}

	return uid, url, nil
}

// GenerateUploadURLs 批量生成上传预签名
func (s *ImageForOSS) GenerateUploadURLs(ctx context.Context, images []*domain.Image) (map[int64]string, error) {
	var (
		urlChan = make(chan map[int64]string, len(images))
		errChan = make(chan error, len(images))
		wg      = &sync.WaitGroup{}
	)

	for _, image := range images {
		wg.Add(1)
		go func(img *domain.Image) {
			defer wg.Done()
			id, url, err := s.GenerateUploadURL(ctx, img)
			fmt.Printf("id: %d url: %s, err: %v \n", id, url, err)
			if err != nil {
				errChan <- err
			}

			urlChan <- map[int64]string{id: url}
		}(image)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(urlChan)
	}()

	var errorSlice []error
	for err := range errChan {
		errorSlice = append(errorSlice, err)
	}

	if len(errorSlice) > 0 {
		for _, err := range errorSlice {
			log.AppLogger.Errorf("Generate upload urls error: %v \n", err)
		}
		return nil, fmt.Errorf("Generate upload urls error: %v \n", errorSlice)
	}

	urls := make(map[int64]string)
	for urlMap := range urlChan {
		for k, v := range urlMap {
			urls[k] = v
		}
	}

	return urls, nil
}

func (s *ImageForOSS) GenerateDownloadURL(ctx context.Context, id int64) (string, error) {
	path, contentType, err := s.repo.GetPathAndContentTypeByImageID(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("get path by image id error: %v \n", err)
		return "", err
	}

	var contentTypeStr = _const.ContentTypeStringMap[contentType]

	url, err := s.oss.GeneratePresignedDownloadURL(ctx, path, contentTypeStr)
	if err != nil {
		log.AppLogger.Errorf("generate download url error: %v \n", err)
		return "", err
	}
	return url, nil
}

func (s *ImageForOSS) GenerateDeleteURL(ctx context.Context, id int64) (string, error) {
	path, contentType, err := s.repo.GetPathAndContentTypeByImageID(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("get path by image id error: %v \n", err)
		return "", err
	}

	var contentTypeStr = _const.ContentTypeStringMap[contentType]

	url, err := s.oss.GeneratePresignedDeleteURL(ctx, path, contentTypeStr)
	if err != nil {
		log.AppLogger.Errorf("generate download url error: %v \n", err)
		return "", err
	}
	return url, nil
}

func (s *ImageForOSS) Save(ctx context.Context, image *domain.Image) (int64, error) {
	imageMap, err := idgen.GetUidCache(ctx, image.ImageID, s.cache)
	if err != nil {
		log.AppLogger.Errorf("get uid cache error: %v \n", err)
		return 0, err
	}

	image.ImageID, err = strconv.ParseInt(imageMap["imageID"], 10, 64)
	image.OwnerID, err = strconv.ParseInt(imageMap["ownerID"], 10, 64)
	image.OwnerType, err = strconv.ParseInt(imageMap["ownerType"], 10, 64)
	image.Path = imageMap["path"]
	image.SHA256Hash = imageMap["sha256"]
	image.ContentType, err = strconv.ParseInt(imageMap["contentType"], 10, 64)
	if err != nil {
		log.AppLogger.Errorf("parse uid cache error: %v \n", err)
		return 0, err
	}

	id, err := s.repo.UploadImage(ctx, image)
	if err != nil {
		log.AppLogger.Errorf("save image error: %v \n", err)
		return 0, err
	}
	return id, nil
}

func (s *ImageForOSS) SaveMore(ctx context.Context, image []*domain.Image) (map[int]int64, error) {
	var (
		total      = len(image)
		resultChan = make(chan map[int]int64, total)
		errChan    = make(chan error, total)
		wg         = &sync.WaitGroup{}
	)
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(idx int, img *domain.Image) {
			defer wg.Done()

			id, err := s.Save(ctx, img)
			if err != nil {
				select {
				case errChan <- err:
				default:
					log.AppLogger.Warnf("save failed for image %d, error dropped due to full channel: %v", idx, err)
				}
				return
			}
			select {
			case resultChan <- map[int]int64{idx: id}:
			case <-ctx.Done():
				log.AppLogger.Warnf("Context done while sending result for image %d", idx)
			}
		}(i, image[i])
	}

	var (
		imgs     = make(map[int]int64)
		errSlice []error
	)

	for resultChan != nil || errChan != nil {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				for k, v := range result {
					imgs[k] = v
				}
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errSlice = append(errSlice, err)
			}
		case <-ctx.Done():
			log.AppLogger.Warnln("save more operation cancelled via context")
			return imgs, ctx.Err()
		}
	}

	if len(errSlice) > 0 {
		for _, err := range errSlice {
			log.AppLogger.Warnf("save image matedata error: %v", err)
		}
		return imgs, fmt.Errorf("%d out of %d uploads failed", len(errSlice), total)
	}

	return imgs, nil
}

func (s *ImageForOSS) Get(ctx context.Context, id int64) (*domain.Image, error) {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("get image error: %v \n", err)
		return nil, err
	}
	return image, nil
}

func (s *ImageForOSS) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("get image error: %v \n", err)
		return err
	}
	return nil
}

func (s *ImageForOSS) Update(ctx context.Context, oldImageID int64, newImage *domain.Image) (int64, error) {
	id, err := s.repo.Update(ctx, oldImageID, newImage)
	if err != nil {
		log.AppLogger.Errorf("update image error: %v \n", err)
		return 0, err
	}
	return id, nil
}
