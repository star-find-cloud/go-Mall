package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/internal/service"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}
type LoginRequest struct {
	id       int    `json:"id" binding:"required,id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := utils.ParsePathParamInt(c, "id")
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

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	token, err := h.UserService.LoginByID(c.Request.Context(), req.id, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, LoginResponse{
		token,
	})
}
