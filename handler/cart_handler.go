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

type CartHandler struct {
	CartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{CartService: cartService}
}

// CartCreateResponse 创建用户购物车响应体
type CartCreateResponse struct {
	// @Description 购物车ID
	ID int64 `json:"id"`
	// @Description 消息
	Msg string `json:"msg"`
}

// Create 创建用户购物车
// @Summary 创建用户购物车
// @Description 创建用户购物车
// @Tags 购物车
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} CartCreateResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/cart/create [post]
func (h *CartHandler) Create(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	var cart = &domain.Cart{
		UserID: customClaims.UserID,
	}
	id, err := h.CartService.Create(c, cart)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "create cart failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, &CartCreateResponse{
		ID:  id,
		Msg: "create cart success",
	})
	return
}

// GetByID 通过购物车ID获取用户购物车
// @Summary 通过购物车ID获取用户购物车
// @Description 通过购物车ID获取用户购物车
// @Tags 购物车
// @Accept  json
// @Produce  json
// @Param id path int64 true "购物车ID"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/cart/{id} [get]
func (h *CartHandler) GetByID(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	id := c.GetInt64("id")
	cart, err := h.CartService.GetByID(c, id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "get cart failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, cart)
}

// GetByUserID 通过用户ID获取用户购物车
// @Summary 通过用户ID获取用户购物车
// @Description 通过用户ID获取用户购物车
// @Tags 购物车
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param id path int64 true "用户ID"
// @Success 200 {object} domain.Cart
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/cart/user [get]
func (h *CartHandler) GetByUserID(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	id := c.GetInt64("id")
	cart, err := h.CartService.GetByUserID(c, id)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "get cart failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, cart)
	return
}

// AddProductRequest 向购物车添加商品请求体
type AddProductRequest struct {
	// @Description 购物车ID
	CartID int64 `json:"cartId" binding:"required"`
	// @Description 商品ID
	ProductID int64 `json:"productId" binding:"required"`
	// @Description 商品标题
	ProductTitle string `json:"productTitle" binding:"required"`
	// @Description 商品加入购物车时价格
	CreatePrice float64 `json:"createPrice" binding:"required"`
	// @Description 商品图片Oss路径
	ProductImageOss string `json:"productImageOss" binding:"required"`
	// @Description 商品数量
	Quantity int64 `json:"quantity" binding:"required"`
	// @Description 商品规格
	Specs map[string]interface{} `json:"specs" binding:"required"`
}

// AddProduct 向购物车添加商品
// @Summary 向购物车添加商品
// @Description 向购物车添加商品
// @Tags 购物车
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param id body AddProductRequest true "购物车ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/cart/add [post]
func (h *CartHandler) AddProduct(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	if customClaims.Roles != _const.UserRole {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	var req = &AddProductRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	//var cartItems = make([]*domain.CartItemVO, 0)
	//for _, item := range req {
	//	cartItems = append(cartItems, &domain.CartItemVO{
	//		CartID:          item.CartID,
	//		ProductID:       item.ProductID,
	//		ProductTitle:    item.ProductTitle,
	//		CreatePrice:     item.CreatePrice,
	//		NowPrice:        item.CreatePrice,
	//		ProductImageOss: item.ProductImageOss,
	//		Quantity:        item.Quantity,
	//		Specs:           item.Specs,
	//	})
	//}
	crat, err := h.CartService.GetByUserID(c, customClaims.UserID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "get cart failed", err)
		return
	}

	cartItem := &domain.CartItemVO{
		CartID:       crat.ID,
		ProductID:    req.ProductID,
		ProductTitle: req.ProductTitle,
		CreatePrice:  req.CreatePrice,
		NowPrice:     req.CreatePrice,
		Quantity:     req.Quantity,
		Specs:        req.Specs,
	}

	err = h.CartService.AddProduct(c, cartItem)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "add product to cart failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "add product to cart success")
	return
}
