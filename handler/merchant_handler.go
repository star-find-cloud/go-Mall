// Package handler implements the HTTP handlers for the Star Mall API.
//
// @title Star Mall API
// @version 1.0
// @description This is the API server for Star Mall merchant management system.
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.starmall.com/support
// @contact.email support@starmall.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath /api/v1
// @schemes http https
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	applogger "github.com/star-find-cloud/star-mall/pkg/logger"
	appproto "github.com/star-find-cloud/star-mall/protobuf/pb"

	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

// MerchantHandler handles HTTP requests related to merchant operations.
type MerchantHandler struct {
	MerchantService service.MerchantService
	UserService     service.UserService
}

func NewMerchantHandler(merchantService service.MerchantService, userService service.UserService) *MerchantHandler {
	return &MerchantHandler{
		MerchantService: merchantService,
		UserService:     userService,
	}
}

// MerchantRegisterRequest represents the request body for merchant registration.
type MerchantRegisterRequest struct {
	// Name is the merchant's business name
	// Required: true
	// Min length: 2
	// Max length: 32
	Name string `json:"name,omitempty" binding:"required,min=2,max=32"`

	// Phone is the merchant's contact phone number in E.164 format
	// Required: true
	// Example: +8613800138000
	Phone string `json:"phone,omitempty" binding:"required"`

	// Email is the merchant's contact email address
	// Required: true
	// Example: merchant@example.com
	Email string `json:"email,omitempty" binding:"required"`

	// Password for merchant account
	// Required: true
	// Min length: 8
	// Format: password
	Password string `json:"password,omitempty" binding:"required,min=8"`

	// RealName is the legal name of the merchant's representative
	// Required: true
	// Min length: 2
	// Max length: 32
	RealName string `json:"realName,omitempty" binding:"required,min=2,max=32"`

	// RealID is the ID card number of the merchant's representative
	// Required: true
	// Length: 18
	RealID string `json:"realID,omitempty" binding:"required,min=18,max=18"`

	// BusinessType represents the categories of business the merchant operates in
	// Required: true
	// Example: [1, 3, 5]
	BusinessType []int64 `json:"businessType,omitempty" binding:"required"`

	// 用户收到的邮箱验证码, 用于验证用户身份
	VerificationCode string `json:"verificationCode,omitempty" binding:"required"`
}

// MerchantRegisterResponse represents the response for a successful merchant registration.
type MerchantRegisterResponse struct {
	// HTTP status code
	Code int `json:"code" example:"201"`

	// 商家ID
	MerchantID int64 `json:"merchantId" example:"12345"`

	// 商家对应的 用户ID
	UserID int64 `json:"userId" example:"67890"`

	// Email 邮箱地址
	Email string `json:"email" example:"merchant@example.com"`

	// Phone 商家店铺电话
	Phone string `json:"phone" example:"+8613800138000"`

	// UploadToken 后续传输图片所需的 临时用户凭证
	UploadToken string `json:"uploadToken" example:"<KEY>"`
}

// MerchantUpdateRequest 商家更新请求
// @Summary 商家更新请求
// @Description 商家更新请求
type MerchantUpdateRequest struct {
	// Name 商家店铺名称
	// Required: true
	// Min length: 2
	// Max length: 32
	Name string `json:"name,omitempty" binding:"required,min=2,max=32"`

	// Phone 商家店铺电话
	Phone string `json:"phone,omitempty" binding:"required,e164"`

	// Email 商家店铺邮箱
	Email string `json:"email,omitempty" binding:"required"`

	// CateID 商家店铺分类
	CateID int64 `json:"cateId,omitempty" binding:"required"`

	// BusinessType 营业类型
	BusinessType []int64 `json:"businessType,omitempty" binding:"required"`
}

// Register 注册商家
// @Summary 注册商家
// @Description 注册商家
// @Tags merchant
// @Accept json
// @Produce json
// @Param request body MerchantRegisterRequest true "Merchant registration information"
// @Success 201 {object} MerchantRegisterResponse "Successfully registered merchant"
// @Failure 400 {object} MerchantRegisterResponse "Invalid request parameters"
// @Failure 409 {object} MerchantRegisterResponse "Email or phone already registered"
// @Failure 500 {object} MerchantRegisterResponse "Internal server error"
// @Router /api/v1/merchants/register [post]
func (h *MerchantHandler) Register(c *gin.Context) {
	var req = MerchantRegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		applogger.HttpLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// 创建用户记录
	_, userID, err := h.UserService.Register(c.Request.Context(), &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		RoleID:   _const.MerchantRole,
		Name:     req.Name,
	})
	if err != nil {
		// 记录错误但不影响商户注册结果
		applogger.AppLogger.Errorf("create user record for merchant failed: %v", err)
	}

	// 创建商户记录
	merchantID, token, err := h.MerchantService.Register(c.Request.Context(), &domain.Merchant{
		UserID:       userID,
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Password:     req.Password,
		RealName:     req.RealName,
		RealID:       req.RealID,
		BusinessType: req.BusinessType,
	})

	if err != nil {
		applogger.AppLogger.Errorf("merchant register failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "register failed", err)
		return
	}

	// 上传营业执照图片
	var imageProto = &appproto.ImageProto{}
	if err = c.ShouldBindWith(imageProto, binding.ProtoBuf); err != nil {
		c.ProtoBuf(http.StatusBadRequest, &appproto.ErrorResponse{
			Code:    appproto.ErrorCode_INVALID_CHUNK,
			Message: err.Error(),
		})
		applogger.AppLogger.Errorf("failed to unmarshal request: %v", err)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, MerchantRegisterResponse{
		Code:        http.StatusCreated,
		MerchantID:  merchantID,
		UserID:      userID,
		Email:       req.Email,
		Phone:       req.Phone,
		UploadToken: token,
	})
}

// GetMerchant 获取商户信息
// @Summary 根据商户 ID 获取商户信息
// @Description 根据商户 ID 获取商户信息, 返回商户信息
// @Tags merchant
// @Accept json
// @Produce json
// @Param merchantID path int true "Merchant ID" minimum(1)
// @Success 200 {object} domain.Merchant "Merchant details retrieved successfully"
// @Failure 400 {object} MerchantRegisterResponse "Invalid merchant ID"
// @Failure 404 {object} MerchantRegisterResponse "Merchant not found"
// @Failure 500 {object} MerchantRegisterResponse "Internal server error"
// @Router /api/v1/merchants/getMerchant/{merchantID} [get]
func (h *MerchantHandler) GetMerchant(c *gin.Context) {
	merchantID, err := utils.ParsePathParamInt64(c, "merchantID")
	if err != nil {
		applogger.HttpLogger.Errorf("invalid merchantID: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid merchantID", err)
		return
	}

	merchant, err := h.MerchantService.GetByID(c.Request.Context(), merchantID)
	if err != nil {
		applogger.AppLogger.Errorf("get merchant failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "get merchant failed", err)
		return
	}

	if merchant == nil {
		utils.RespondError(c, http.StatusNotFound, "merchant not found", nil)
		return
	}

	utils.RespondJSON(c, http.StatusOK, merchant)
}

// Select godoc
// @Summary 根据商家名称搜索商家
// @Description 根据商家名称搜索商家, 返回匹配的商家列表
// @Tags merchant
// @Accept json
// @Produce json
// @Param name query string true "Merchant name to search for (minimum 2 characters)" minLength(2)
// @Success 200 {array} domain.Merchant "List of merchants matching the search criteria"
// @Success 204 {object} MerchantRegisterResponse "No merchants found"
// @Failure 400 {object} MerchantRegisterResponse "Invalid merchant name"
// @Failure 500 {object} MerchantRegisterResponse "Internal server error"
// @Router /api/v1/merchants/search/{merchant_name} [get]
func (h *MerchantHandler) Select(c *gin.Context) {
	merchantName := c.Query("merchant_name")
	if len(merchantName) < 2 {
		utils.RespondError(c, http.StatusBadRequest, "merchant name must be at least 2 characters", nil)
		return
	}

	merchants, err := h.MerchantService.GetByName(c.Request.Context(), merchantName)
	if err != nil {
		applogger.AppLogger.Errorf("get merchant failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "search merchants failed", err)
		return
	}

	if merchants == nil {
		utils.RespondJSON(c, http.StatusNoContent, gin.H{
			"code":    http.StatusNoContent,
			"message": "no merchants found",
		})
		return
	}

	utils.RespondJSON(c, http.StatusOK, merchants)
}

// Update 修改商家信息
// @Summary 修改商家信息
// @Description 修改商家信息
// @Tags merchant
// @Accept json
// @Produce json
// @Param request body MerchantUpdateRequest true "Merchant update request"
// @Success 200 {object} string "merchant updated successfully"
// @Failure 400 {object} string "invalid request"
// @Failure 401 {object} string "unauthorized"
// @Failure 500 {object} string "internal server error"
// @Router /api/v1/merchants/update [patch]
func (h *MerchantHandler) Update(c *gin.Context) {
	// 获取商家 token
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", nil)
		return
	}

	// 解析请求体
	var (
		req      = &MerchantUpdateRequest{}
		merchant = &domain.Merchant{}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		applogger.HttpLogger.Errorf("update failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	merchant.ID = customClaims.UserID
	merchant.Name = req.Name
	merchant.Phone = req.Phone
	merchant.Email = req.Email
	merchant.CateID = req.CateID
	merchant.BusinessType = req.BusinessType

	// 调用 service 层更新商家信息
	err := h.MerchantService.Update(c.Request.Context(), merchant)
	if err != nil {
		applogger.AppLogger.Errorf("update failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "update failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "merchant updated successfully")
}

// Delete 删除商家
// @Summary 删除商家
// @Description 删除商家
// @Tags merchant
// @Accept json
// @Produce json
// @Success 200 {object} string "merchant deleted successfully"
// @Failure 401 {object} string "unauthorized"
// @Failure 500 {object} string "internal server error"
// @Router /api/v1/merchants/delete [delete]
func (h *MerchantHandler) Delete(c *gin.Context) {
	// 获取商家 token
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", nil)
		return
	}

	err := h.MerchantService.Delete(c.Request.Context(), customClaims.UserID)
	if err != nil {
		applogger.AppLogger.Errorf("delete failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "delete failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "merchant deleted successfully")
}
