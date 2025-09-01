package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

type ProductCreateRequest struct {
	// @Description 商品图片id数组
	ImageID string `json:"image_id"`

	// @Description 商品名称
	Title string `json:"title"`

	// @Description 商品副标题
	SubTitle string `json:"subTitle"`

	// @Description 商品编号
	Sn string `json:"sn"`

	// @Description 商品大类ID
	ProductTypeID int `json:"productTypeID"`

	// @Description 商品分类ID
	CateID int64 `json:"cateID"`

	// @Description 商品价格
	Price float64 `json:"price"`

	// @Description 商品市场价
	MarketPrice float64 `json:"marketPrice"`

	// @Description 商品属性
	Attr string `json:"attr"`

	// @Description 关键字
	Keywords string `json:"keywords"`

	// @Description 商品描述
	Desc string `json:"desc"`

	// @Description 商品内容
	Content string `json:"content"`

	// @Description 是否预售
	IsBooking bool `json:"isBooking"`

	// @Description 预售时间
	BookingTime int64 `json:"bookingTime"`

	// @Description 商品版本
	Version string `json:"version"`

	// @Description 商品规格
	Specs string `json:"specs"`

	// @Description 品牌
	Brand string `json:"brand"`
}

// ProductCreateResponse 创建商品响应体
// @Description 创建商品响应体
type ProductCreateResponse struct {
	// @Description 状态码
	Code int `json:"code"`
	// @Description 商品ID
	ID int64 `json:"id"`
	// @Description 消息
	Message string `json:"message"`
}

// Create 创建商品
// @Summary Create 创建商品
// @Description 创建商品
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param product body ProductCreateRequest true "product"
// @Success 200 {object} ProductCreateResponse
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/products/create [put]
func (h *ProductHandler) Create(c *gin.Context) {
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
	if customClaims.Roles != _const.MerchantRole {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}
	merchantID := customClaims.UserID

	var req = &ProductCreateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var product = &domain.Product{
		MerchantID: merchantID,
		//ImageID:       req.ImageID,
		Title:         req.Title,
		SubTitle:      req.SubTitle,
		Brand:         req.Brand,
		ProductSn:     req.Sn,
		ProductTypeID: req.ProductTypeID,
		CateID:        req.CateID,
		Price:         req.Price,
		MarketPrice:   req.MarketPrice,
		Attr:          req.Attr,
		Version:       req.Version,
		Keywords:      req.Keywords,
		Desc:          req.Desc,
		Content:       req.Content,
		IsBooking:     req.IsBooking,
		BookingTime:   req.BookingTime,
		Specs:         req.Specs,
	}
	id, err := h.ProductService.Create(c.Request.Context(), product, merchantID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "create failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, ProductCreateResponse{
		Code:    200,
		ID:      id,
		Message: "product created successfully",
	})
	return
}

// GetProduct 获取商品
// @Summary GetProduct 获取商品
// @Description 获取商品
// @Accept json
// @Produce json
// @Tags 商品
// @Param id path int64 true "商品ID"
// @Success 200 {object} domain.Product "product get successfully"
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/products/getProduct/:id [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		applog.AppLogger.Errorf("解析id失败: %d", idStr)
		return
	}

	product, err := h.ProductService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "get failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, product)
	//applog.AppLogger.Info("product get successfully")
	return
}

// GetProductByMerchant 通过商家获取商品
// @Summary GetProductByMerchant 通过商家获取商品
// @Description 通过商家获取商品
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param product body ProductCreateRequest true "product"
// @Success 200 {object} domain.Product "product get successfully"
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/products/search/:MerchantID [get]
func (h *ProductHandler) GetProductByMerchant(c *gin.Context) {
	_, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	merchantID := c.GetInt64("merchantID")
	products, err := h.ProductService.GetByMerchantID(c.Request.Context(), merchantID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "get failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, products)
}

// ProductUpdateRequest 更新商品请求体
// @Description 更新商品请求体
// @Accept json
// @Produce json
type ProductUpdateRequest struct {
	// @Description 商品ID
	ID int64 `json:"id"`

	// @Description 标题
	Title string `json:"title"`

	// @Description 子标题
	SubTitle string `json:"sub_title"`

	// @Description 商品编号
	ProductSn string `json:"product_sn"`

	// @Description 商品类型ID
	ProductTypeID int `json:"product_type_id"`

	// @Description 分类ID
	CateID int64 `json:"cate_id"`

	// @Description 价格
	Price float64 `json:"price"`

	// @Description 市场价
	MarketPrice float64 `json:"market_price"`

	// @Description 属性
	Attr string `json:"attr"`

	// @Description 关键字
	Keywords string `json:"keywords"`

	// @Description 描述
	Desc string `json:"desc"`

	// @Description 内容
	Content string `json:"content"`

	// @Description 是否预售
	IsBooking bool `json:"isBooking"`

	// @Description 预售时间
	BookingTime int64 `json:"booking_time"`
}

// Update 更新商品
// @Summary Update 更新商品
// @Description 更新商品
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param product body ProductUpdateRequest true "product"
// @Success 200 {object} string "product updated successfully"
// @Failure 401 {object} string "没有权限"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/products/update [patch]
func (h *ProductHandler) Update(c *gin.Context) {
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
	if customClaims.Roles != _const.MerchantRole {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}
	if customClaims.UserID == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}

	var req = &ProductUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var product = &domain.Product{
		ID:            req.ID,
		MerchantID:    customClaims.UserID,
		Title:         req.Title,
		SubTitle:      req.SubTitle,
		ProductSn:     req.ProductSn,
		ProductTypeID: req.ProductTypeID,
		CateID:        req.CateID,
		Price:         req.Price,
		MarketPrice:   req.MarketPrice,
		Attr:          req.Attr,
		Keywords:      req.Keywords,
		Desc:          req.Desc,
		Content:       req.Content,
		IsBooking:     req.IsBooking,
		BookingTime:   req.BookingTime,
	}

	err := h.ProductService.Update(c, product, customClaims.UserID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "更新商品失败", err)
		return
	}
	utils.RespondJSON(c, http.StatusOK, "商品更新成功")
	return
}

// SearchProductRequest 搜索商品请求体
// @Description 搜索商品请求体
type SearchProductRequest struct {
	Msg    string `json:"msg"`
	Offset int    `json:"offset"`
}

// SearchProduct 搜索商品
// @Summary SearchProduct 搜索商品
// @Description 搜索商品
// @Accept json
// @Produce json
// @Tags 商品
// @Param product body SearchProductRequest true "product"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []domain.Product "product updated successfully"
// @Failure 401 {object} string "没有权限"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/products/search [post]
func (h *ProductHandler) SearchProduct(c *gin.Context) {
	var req = &SearchProductRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var products []domain.Product
	products, err := h.ProductService.Search(c, req.Msg)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "搜索商品失败", err)
		return
	}

	//fmt.Println(products)
	utils.RespondJSON(c, http.StatusOK, products)
	return
}
