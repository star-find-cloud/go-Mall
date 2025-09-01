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

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrderRequest 创建订单请求参数
// @Description: 创建订单请求参数
type CreateOrderRequest struct {
	// @Description: 总金额
	TotalPrice int64 `json:"total_price"`
	// @Description: 支付金额
	PayPrice int64 `json:"pay_price"`
	// @Description: 商品ID
	ProductID int64 `json:"product_id"`
	// @Description: 商品名称
	ProductTitle string `json:"product_title"`
	// @Description: 商品单价
	UnitPrice int64 `json:"unit_price"`
	// @Description: 商品数量
	Quantity int64 `json:"quantity"`
	// @Description: 商品小计
	Subtotal int64 `json:"subtotal"`
}

// CreateOrderResponse 创建订单响应参数
// @Description: 创建订单响应参数
type CreateOrderResponse struct {
	// @Description: 订单ID
	ID int64 `json:"id"`
	// @Description: 创建时间
	CreatedAt int64 `json:"created_at"`
}

// CreateOrder 创建订单
// @Summary CreateOrder 创建订单
// @Description 创建订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "order"
// @Success 200 {object} CreateOrderResponse
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/order/create [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
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
		utils.RespondError(c, http.StatusUnauthorized, "not User", errors.New("not User"))
		return
	}

	userID := customClaims.UserID

	var req = &CreateOrderRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var order = &domain.Order{
		UserID:      userID,
		OrderStatus: "待付款",
		TotalPrice:  req.TotalPrice,
		PayPrice:    req.PayPrice,
	}

	var orderItem = &domain.OrderItem{
		OrderID:      order.ID,
		ProductID:    req.ProductID,
		ProductTitle: req.ProductTitle,
		UnitPrice:    req.UnitPrice,
		Quantity:     req.Quantity,
		Subtotal:     req.Subtotal,
	}

	id, createAt, err := h.service.Create(c.Request.Context(), order, orderItem, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "create order failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, &CreateOrderResponse{
		ID:        id,
		CreatedAt: createAt,
	})
}

// OrderGetRequest 获取订单请求参数
// @Description: 获取订单请求参数
type OrderGetRequest struct {
	// @Description: 订单ID
	ID int64 `json:"id"`
}

type OrderGetResponse struct {
	// @Description: 订单ID
	ID int64 `json:"id"`
	// @Description: 用户ID
	UserID int64 `json:"userID"`
	// @Description: 订单状态
	OrderStatus string `json:"orderStatus"`
	// @Description: 总金额
	TotalPrice int64 `json:"totalPrice"`
	// @Description: 支付金额
	PayPrice int64 `json:"payPrice"`
	// @Description: 创建时间
	CreatedAt int64 `json:"createdAt"`
	// @Description: 更新时间
	UpdatedAt int64 `json:"updatedAt"`
	// @Description: 支付方式ID
	PaymentMethodID int64 `json:"paymentID"`
	// @Description: 收货地址ID
	ShippingAddressID int64 `json:"shippingID"`
	// @Description: 订单项
	OrderItems *domain.OrderItem `json:"orderItems"`
}

// GetOrder 获取订单
// @Summary 获取订单
// @Description: 获取订单
// @Tags 订单
// @Accept  json
// @Produce  json
// @Param id path int true "订单ID"
// @Success 200 {object} OrderGetResponse
// @Failure 401 {object} string "invalid token claims"
// @Failure 500 {object} string "服务器错误"
// @Router /api/v1/order/get [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
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
		utils.RespondError(c, http.StatusUnauthorized, "not User", errors.New("not User"))
		return
	}

	var req = &OrderGetRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	order, orderItem, err := h.service.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "get order failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, &OrderGetResponse{
		ID:                order.ID,
		UserID:            order.UserID,
		OrderStatus:       order.OrderStatus,
		TotalPrice:        order.TotalPrice,
		PayPrice:          order.PayPrice,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		PaymentMethodID:   order.PaymentMethodID,
		ShippingAddressID: order.ShippingAddressID,
		OrderItems:        orderItem,
	})

	return
}
