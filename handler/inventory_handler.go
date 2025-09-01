package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// InventoryCreateRequest 库存创建请求体
// @Describe: 库存创建请求体
type InventoryCreateRequest struct {
	// @Describe: 商品ID
	ProductID int64 `json:"product_id"`

	// @Describe: 库存数量
	Count int64 `json:"count"`

	// @Describe: 低库存阈值
	LowStockThreshold int64 `json:"low_stock_threshold"`
}

// Create 创建库存
// @Summary 创建库存
// @Description 创建库存
// @Accept json
// @Produce json
// @Param inventory body InventoryCreateRequest true "inventory"
// @Success 201 {object} string "inventory created successfully"
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/inventory/create [put]
func (h *InventoryHandler) Create(c *gin.Context) {
	// 解析jwt
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	// 验证jwt
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}
	// 验证角色是否是商家
	if customClaims.Roles != _const.MerchantRole {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}

	// 解析请求参数
	var req = &InventoryCreateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var inventory = &domain.Inventory{
		ProductID:         req.ProductID,
		AvailableStock:    req.Count,
		LowStockThreshold: req.LowStockThreshold,
	}

	err := h.service.Create(c, inventory, customClaims.UserID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "创建库存失败", err)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "创建库存成功")
	return
}

// Update 更新库存
// @Summary 更新库存
// @Description 更新库存
// @Accept json
// @Produce json
// @Param inventory body InventoryCreateRequest true "inventory"
// @Success 201 {object} string "inventory created successfully"
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/inventory/update [patch]
func (h *InventoryHandler) Update(c *gin.Context) {
	// 解析jwt
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	// 验证jwt
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}
	// 验证角色是否是商家
	if customClaims.Roles != _const.MerchantRole {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}

	// 解析请求参数
	var req = &InventoryCreateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var inventory = &domain.Inventory{
		ProductID:         req.ProductID,
		AvailableStock:    req.Count,
		LowStockThreshold: req.LowStockThreshold,
	}

	err := h.service.Update(c, inventory, customClaims.UserID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "创建库存失败", err)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "创建库存成功")
	return
}

// GetInventory 获取库存
// @Summary GetInventory 获取库存
// @Description 获取库存
// @Accept json
// @Produce json
// @Param id path int64 true "inventory id"
// @Success 200 {object} domain.Inventory "inventory"
// @Failure 401 {object} string "invalid token claims"
// @Failure 404 {object} string "inventory not found"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/inventory/{id} [get]
func (h *InventoryHandler) GetByID(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	var id = c.GetInt64("id")
	inventory, err := h.service.GetByID(c, id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "获取库存失败", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, inventory)
	return
}

// GetByMerchant 通过商家id获取库存
// @Summary GetByMerchant 通过商家id获取库存
// @Description 通过商家id获取库存
// @Accept json
// @Produce json
// @Param id path int64 true "inventory id"
// @Success 200 {object} []domain.Inventory "inventory 数组"
// @Failure 401 {object} string "invalid token claims"
// @Failure 404 {object} string "inventory not found"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/inventory/search/{MerchantID} [get]
func (h *InventoryHandler) GetByMerchant(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}
	var id = c.GetInt64("id")

	inventorys, err := h.service.GetByMerchantID(c, id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "获取库存失败", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, inventorys)
	return
}
