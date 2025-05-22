package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	appproto "github.com/star-find-cloud/star-mall/protobuf/pb"
	"github.com/star-find-cloud/star-mall/service"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	OssService *service.OSSService
}

func NewImageHandler(ossService *service.OSSService) *ImageHandler {
	return &ImageHandler{OssService: ossService}
}

// @Summary      上传图片到OSS
// @Description  接收客户端上传的Protobuf格式图片数据，存储至OSS并返回存储路径
// @Tags         Image
// @Accept application/x-protobuf
// @Produce application/x-protobuf
// @Param        image  body      appproto.ImageProto  true  "图片数据"
// @Success      200    {object}  appproto.UploadResponse  "OSS存储路径"
// @Failure      400    {object}  appproto.ErrorResponse   "请求解析失败或上传超时"
// @Failure      500    {object}  appproto.ErrorResponse   "内部服务器错误"
// @Router /api/v1/{owner_type}/{id}/image/upload [post]
func (h ImageHandler) UploadImage(c *gin.Context) {
	var req *appproto.ImageProto
	if err := c.ShouldBindWith(&req, binding.ProtoBuf); err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_INVALID_CHUNK,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
		return
	}

	ownerType := c.Param("owner_type")
	ownerID := c.Param("id")
	req.OwnerType, _ = strconv.ParseInt(ownerType, 10, 64)
	req.OwnerId, _ = strconv.ParseInt(ownerID, 10, 64)

	ossPath, _, err := h.OssService.UploadImage(c, req)
	if err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_UPLOAD_TIMEOUT,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to upload image: %v", err)
		return
	}

	// 返回响应
	c.ProtoBuf(http.StatusOK, &appproto.UploadResponse{OssPath: ossPath})
}

// @Summary      批量上传图片到OSS
// @Description  接收客户端上传的Protobuf格式图片数组，批量存储至OSS并返回存储路径列表
// @Tags         Image
// @Accept application/x-protobuf
// @Produce application/x-protobuf
// @Param        images  body      appproto.ImagesProto  true  "图片数组数据"
// @Success      200    {array}     appproto.UploadResponse  "OSS存储路径列表"
// @Failure      400    {object}    appproto.ErrorResponse   "请求解析失败或上传超时"
// @Failure      500    {object}    appproto.ErrorResponse   "内部服务器错误"
// @Router /api/v1/{owner_type}/{id}/images/upload [post]
func (h ImageHandler) UploadImages(c *gin.Context) {
	var reqs appproto.ImagesProto
	if err := c.ShouldBindWith(&reqs, binding.ProtoBuf); err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_INVALID_CHUNK,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
		return
	}

	ownerType, err := strconv.ParseInt(c.Param("owner_type"), 10, 64)
	ownerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errors.New("invalid owner_type or id")
		return
	}

	for _, image := range reqs.Images {
		image.OwnerType = ownerType
		image.OwnerId = ownerID
	}

	ossPaths, num, err := h.OssService.UploadImages(c, reqs.Images)
	if err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_UPLOAD_TIMEOUT,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to upload images: %v", err)
		return
	}

	for _, ossPath := range ossPaths {
		c.ProtoBuf(http.StatusOK, &appproto.UploadResponse{OssPath: ossPath, SuccessCount: uint32(num)})
	}
}
