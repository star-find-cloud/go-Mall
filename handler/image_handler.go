package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/star-find-cloud/star-mall/internal/service"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	appproto "github.com/star-find-cloud/star-mall/protobuf/pb"
	"net/http"
)

type ImageHandler struct {
	OssService *service.OSSService
}

func NewImageHandler(ossService *service.OSSService) *ImageHandler {
	return &ImageHandler{
		OssService: ossService,
	}
}

// @Summary      上传图片到OSS
// @Description  接收客户端上传的Protobuf格式图片数据，存储至OSS并返回存储路径
// @Tags         Image
// @Accept       protobuf
// @Produce      protobuf
// @Param        image  body      appproto.ImageProto  "图片数据"  true
// @Success      200    {object}  appproto.UploadResponse  "OSS存储路径"
// @Failure      400    {object}  appproto.ErrorResponse   "请求解析失败或上传超时"
// @Failure      500    {object}  appproto.ErrorResponse   "内部服务器错误"
// @Router       /image/upload [post]
func (h ImageHandler) UploadImage(c *gin.Context) {
	var req appproto.ImageProto
	if err := c.ShouldBindWith(&req, binding.ProtoBuf); err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_INVALID_CHUNK,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
		return
	}

	ossPath, err := h.OssService.UploadImage(c, &req)
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
// @Accept       protobuf
// @Produce      protobuf
// @Param        images  body      []appproto.ImageProto  "图片数组数据"  true
// @Success      200    {array}     appproto.UploadResponse  "OSS存储路径列表"
// @Failure      400    {object}    appproto.ErrorResponse   "请求解析失败或上传超时"
// @Failure      500    {object}    appproto.ErrorResponse   "内部服务器错误"
// @Router       /image/upload/batch [post]
func (h ImageHandler) UploadImages(c *gin.Context) {
	var reqs []*appproto.ImageProto
	if err := c.ShouldBindWith(&reqs, binding.ProtoBuf); err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_INVALID_CHUNK,
			Message: err.Error(),
		})
		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
		return
	}

	ossPaths, num, err := h.OssService.UploadImages(c, reqs)
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
