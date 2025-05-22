package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}
type LoginRequest struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"userId"`
}

type RegisterRequest struct {
	Name     string `json:"name,omitempty" binding:"required,min=2,max=32"`
	Password string `json:"password,omitempty" binding:"required,min=8"`
	Email    string `json:"email,omitempty" binding:"required"`
	Phone    string `json:"phone,omitempty" binding:"required,e164"`
	Sex      int    `json:"sex,omitempty" binding:"required"`
	Image    string `json:"image,omitempty" binding:"required"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// @Summary GetUser
// @Description 获取用户信息
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} utils.ResponseError
// @Failure 404 {object} utils.ResponseError
// @Router /api/v1/users/{id} [get]`
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := utils.ParsePathParamInt64(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid product id", err)
		return
	}

	user, err := h.UserService.GetUsrByID(context.Background(), id)
	if err != nil {
		logger.AppLogger.Errorf("failed to get user: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	if user == nil {
		utils.RespondError(c, http.StatusNotFound, "user not found", nil)
		return
	}
	utils.RespondJSON(c, http.StatusOK, user)
}

// @Summary Login
// @Description 用户通过 email 或 id 登录
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ResponseError
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("login failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	if req.ID == 0 && req.Email == "" {
		utils.RespondError(c, http.StatusBadRequest, "either id or email is required", nil)
		return
	}

	var (
		token string
		err   error
	)
	if req.ID != 0 {
		token, err = h.UserService.LoginByID(c.Request.Context(), req.ID, req.Password)
	} else {
		token, err = h.UserService.LoginByEmail(c.Request.Context(), req.Email, req.Password)
	}

	if err != nil {
		logger.AppLogger.Errorf("login failed: %v", err)
		utils.RespondError(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, LoginResponse{
		token,
		req.ID,
	})
}

// @Summary Register
// @Description 用户注册
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ResponseError
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("login failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	token, id, err := h.UserService.Register(c.Request.Context(), req.Name, req.Password, req.Email, req.Phone, req.Sex)
	if err != nil {
		logger.AppLogger.Errorf("register failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "register failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, LoginResponse{
		token,
		id,
	})
}
