package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
	"strconv"
)

type ImageHandler struct {
	ImageService service.ImageService
}

func NewImageHandler(ImageService service.ImageService) *ImageHandler {
	return &ImageHandler{ImageService: ImageService}
}

//// UploadImageRequest 图片上传请求体
//// @Description 图片上传请求体
//type UploadImageRequest struct {
//	// @Description 图片所属者类型
//	OwnerType int64
//	// @Description 图片所属者ID
//	OwnerId int64
//	// @Description 是否压缩
//	IsCompressed bool
//}

// UploadImageResponse 图片上传响应体
// @Description 图片上传响应体
type UploadImageResponse struct {
	// @Description 返回消息
	Message string
	// @Description 图片ID
	ImageID int64
	// @Description 图片路径
	FilePath string
}

// UploadImage 图片上传
// @Summary 图片上传
// @Description 图片上传
// @Accept json
// @Produce json
// @Tags	Image
// @Success 200 {object} UploadImageResponse
// @Failure 401 {object} string "没有权限"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/image/upload [post]
func (h ImageHandler) UploadImage(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", exists)
		return
	}

	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", ok)
		return
	}
	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", customClaims.Roles)
		return
	}
	if customClaims.UserID == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not merchant"))
		return
	}

	isCompressed, err := strconv.ParseBool(c.PostForm("is_compressed"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "upload failed", err)
		return
	}

	sha256hash := c.PostForm("sha256hash")

	var image = &domain.Image{
		IsCompressed: isCompressed,
		SHA256Hash:   sha256hash,
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "upload failed", err)
		return
	}

	id, filePath, err := h.ImageService.Save(c.Request.Context(), image, file, c)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "upload failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, UploadImageResponse{
		Message:  "upload successfully",
		ImageID:  id,
		FilePath: filePath,
	})
}

// GetImageRequest 获取图片请求体
type GetImageRequest struct {
	// @Description 图片ID
	ID int64 `json:"id"`
}

// GetImage 获取图片
// @Summary GetImage 获取图片
// @Description 获取图片
// @Accept json
// @Produce json
// @Tags	Image
// @Param image body GetImageRequest true "image"
// @Success 200 {object} domain.Image "get successfully"
// @Failure 401 {object} string "没有权限"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/image/getImage [get]
func (h ImageHandler) GetImage(c *gin.Context) {
	var req = &GetImageRequest{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	image, err := h.ImageService.GetImage(c.Request.Context(), req.ID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "get failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, image)
}

//// UploadImage
//// @Summary      上传图片到OSS
//// @Description  接收客户端上传的Protobuf格式图片数据，存储至OSS并返回图片ID
//// @Tags         Image
//// @Accept application/x-protobuf
//// @Produce application/x-protobuf
//// @Param        image  body      appproto.ImageProto  true  "图片数据"
//// @Success      200    {object}  appproto.UploadResponse  "OSS存储路径"
//// @Failure      400    {object}  appproto.ErrorResponse   "请求解析失败或上传超时"
//// @Failure      500    {object}  appproto.ErrorResponse   "内部服务器错误"
//// @Router /api/v1/image/{owner_type}/{owner_id}/upload [post]
//func (h ImageHandler) UploadImage(c *gin.Context) {
//	var req = &appproto.ImageProto{}
//	if err := c.ShouldBindWith(&req, binding.ProtoBuf); err != nil {
//		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
//			Code:    appproto.ErrorCode_INVALID_CHUNK,
//			Message: err.Error(),
//		})
//		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
//		return
//	}
//
//	ownerType := c.Param("owner_type")
//	ownerID := c.Param("owner_id")
//	req.OwnerType, _ = strconv.ParseInt(ownerType, 10, 64)
//	req.OwnerId, _ = strconv.ParseInt(ownerID, 10, 64)
//
//	_, id, err := h.OssService.UploadImage(c, req)
//	if err != nil {
//		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
//			Code:    appproto.ErrorCode_UPLOAD_TIMEOUT,
//			Message: err.Error(),
//		})
//		applog.AppLogger.Errorf("failed to upload image: %v", err)
//		return
//	}
//
//	// 返回响应
//	c.ProtoBuf(http.StatusOK, &appproto.UploadResponse{
//		ImageId: id,
//	})
//}
//
//// @Summary      批量上传图片到OSS
//// @Description  接收客户端上传的Protobuf格式图片数组，批量存储至OSS并返回图片ID列表
//// @Tags         Image
//// @Accept application/x-protobuf
//// @Produce application/x-protobuf
//// @Param        images  body      appproto.ImagesProto  true  "图片数组数据"
//// @Success      200    {array}     appproto.UploadResponse  "OSS存储路径列表"
//// @Failure      400    {object}    appproto.ErrorResponse   "请求解析失败或上传超时"
//// @Failure      500    {object}    appproto.ErrorResponse   "内部服务器错误"
//// @Router /api/v1/image/{owner_type}/{owner_id}/images/upload [post]
//func (h ImageHandler) SaveMore(c *gin.Context) {
//	var reqs appproto.ImagesProto
//	if err := c.ShouldBindWith(&reqs, binding.ProtoBuf); err != nil {
//		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
//			Code:    appproto.ErrorCode_INVALID_CHUNK,
//			Message: err.Error(),
//		})
//		applog.AppLogger.Errorf("failed to unmarshal request: %v", err)
//		return
//	}
//
//	ownerType, err := strconv.ParseInt(c.Param("owner_type"), 10, 64)
//	ownerID, err := strconv.ParseInt(c.Param("owner_id"), 10, 64)
//	if err != nil {
//		errors.New("invalid owner_type or ownerID")
//		return
//	}
//
//	for _, image := range reqs.Images {
//		image.OwnerType = ownerType
//		image.OwnerId = ownerID
//	}
//
//	ids, num, err := h.OssService.SaveMore(c, reqs.Images)
//	if err != nil {
//		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
//			Code:    appproto.ErrorCode_UPLOAD_TIMEOUT,
//			Message: err.Error(),
//		})
//		applog.AppLogger.Errorf("failed to upload images: %v", err)
//		return
//	}
//
//	for _, id := range ids {
//		c.ProtoBuf(http.StatusOK, &appproto.UploadResponse{ImageId: id, SuccessCount: uint32(num)})
//	}
//}
