package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserService service.UserService
}
type LoginRequest struct {
	ID       string `json:"userID,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Role   int64  `json:"role"`
	UserID int64  `json:"userId"`
}

type UserRegisterRequest struct {
	Name             string `json:"name,omitempty" binding:"required,min=2,max=32"`
	Password         string `json:"password,omitempty" binding:"required,min=8"`
	Email            string `json:"email,omitempty" binding:"required"`
	Phone            string `json:"phone,omitempty" binding:"required"`
	Sex              int    `json:"sex,omitempty" binding:"required"`
	VerificationCode string `json:"verificationCode,omitempty"`
}

// UpdateRequest 用户信息更新请求
type UpdateRequest struct {
	// @Description 用户名
	// @Required false
	// @Example "张三"
	Name string `json:"name,omitempty" binding:"omitempty,min=2,max=32"`

	// @Description 邮箱
	// @Required false
	// @Example "example@mail.com"
	Email string `json:"email,omitempty" binding:"omitempty,email"`

	// @Description 手机号
	// @Required false
	// @Example "+8613800138000"
	Phone string `json:"phone,omitempty" binding:"omitempty,e164"`

	// @Description 性别(203:未知 201:男 202:女)
	// @Required false
	// @Enum 201 202 203
	// @Example 1
	Sex int `json:"sex,omitempty" binding:"omitempty,oneof=0 1 2"`
}

// UpdateEmailRequest 用户邮箱更新请求
type UpdateEmailRequest struct {
	// @Description 邮箱
	Email string `json:"email" binding:"required,email"`

	// @Description 验证码
	// @Required true
	// @Example "12AB56"
	VerifyCode string `json:"verifyCode" binding:"required"`
}

type UpdatePasswordRequest struct {
	// @Description 原密码
	// @Required true
	// @Example "<PASSWORD>"
	OldPassword string `json:"oldPassword" binding:"required,min=6"`

	// @Description 新密码
	// @Required true
	// @Example "<PASSWORD>"
	NewPassword string `json:"newPassword" binding:"required,min=6"`

	// @Description 邮箱
	// @Required true
	// @Example "<EMAIL>"
	Email string `json:"email" binding:"required,email"`

	// @Description 验证码
	// @Required true
	// @Example "12AB56"
	VerifyCode string `json:"verifyCode" binding:"required"`
}

type ForgetPasswordRequest struct {
	// @Description 邮箱
	// @Required true
	// @Example "<EMAIL>"
	Email string `json:"email" binding:"required,email"`

	// @Description 验证码
	// @Required true
	// @Example "123456"
	VerifyCode string `json:"verifyCode" binding:"required"`

	// @Description 新密码
	// @Required true
	// @Example "<PASSWORD>"
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type DeleteRequest struct {
	// @Description 邮箱
	Email string `json:"email" binding:"required,email"`

	// @Description 验证码
	VerifyCode string `json:"verifyCode" binding:"required"`
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetUser 获取用户模型
// @Summary GetUser
// @Description 获取用户信息
// @Accept json
// @Produce json
// @Tags 用户
// @Param userID path int true "User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} domain.User "成功返回用户信息"
// @Failure 400 {object} utils.ResponseError
// @Failure 404 {object} utils.ResponseError
// @Router /api/v1/user/get [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", errors.New("invalid token claims"))
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	userID := customClaims.UserID

	user, err := h.UserService.GetByID(context.Background(), userID)
	if user == nil {
		utils.RespondError(c, http.StatusNotFound, "user not found", err)
		return
	}
	if err != nil {
		logger.AppLogger.Errorf("failed to get user: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, user)
	return
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户通过 email 或 userID 登录
// @Accept json
// @Produce json
// @Tags 用户
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ResponseError
// @Router /api/v1/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req = &LoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("login failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	if req.ID == "" && req.Email == "" {
		utils.RespondError(c, http.StatusBadRequest, "either merchantID or email is required", nil)
		return
	}

	var (
		token string
		role  int64
		err   error
		id    int64
	)
	if req.ID != "" {
		id, _ = strconv.ParseInt(req.ID, 10, 64)
		token, role, err = h.UserService.LoginByID(c.Request.Context(), id, req.Password)
	} else {
		token, role, err = h.UserService.LoginByEmail(c.Request.Context(), req.Email, req.Password)
	}
	if err != nil {
		logger.AppLogger.Errorf("login failed: %v", err)
		utils.RespondError(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	//fmt.Printf("token: %s, role: %d, id: %d", token, role, id)
	utils.RespondJSON(c, http.StatusOK, LoginResponse{
		token,
		role,
		id,
	})
	return
}

// Register 用户注册接口
// @Summary 用户注册接口
// @Description 用户注册
// @Accept json
// @Produce json
// @Tags 用户
// @Param request body UserRegisterRequest true "Register request"
// @Success 201 {object} LoginResponse
// @Failure 400 {object} utils.ResponseError
// @Router /api/v1/user/register [put]
func (h *UserHandler) Register(c *gin.Context) {
	var req = UserRegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	get, err := h.UserService.CheckEmailVerificationCode(c.Request.Context(), req.Email, req.VerificationCode)
	if err != nil {
		logger.AppLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}
	if !get {
		utils.RespondError(c, http.StatusBadRequest, "invalid verification code", nil)
		return
	}

	token, id, err := h.UserService.Register(c.Request.Context(), &domain.User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Sex:      req.Sex,
		RoleID:   _const.UserRole,
	})
	if err != nil {
		logger.AppLogger.Errorf("register failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "register failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, LoginResponse{
		token,
		_const.UserRole,
		id,
	})
	return
}

// Update
// @Summary Update user
// @Description 更新用户信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Tags 用户
// @Param request body UpdateRequest true "Update user request"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} string "user updated successfully"
// @Failure 400 {object} utils.ResponseError
// @Failure 401 {object} utils.ResponseError
// @Failure 500 {object} utils.ResponseError
// @Router /api/v1/user/update [patch]
func (h *UserHandler) Update(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", errors.New("invalid token claims"))
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	userID := customClaims.UserID

	// 解析请求体
	var req = &UpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("update failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// 调用 service 层更新用户信息
	err := h.UserService.Update(
		c.Request.Context(),
		req.Name,
		req.Phone,
		req.Email,
		userID,
		req.Sex,
	)
	if err != nil {
		logger.AppLogger.Errorf("update failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "update failed", err)
		return
	}

	// 返回成功响应
	utils.RespondJSON(c, http.StatusOK, "user updated successfully")
	return
}

// UpdateEmail 修改邮箱
func (h *UserHandler) UpdateEmail(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", errors.New("invalid token claims"))
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	userID := customClaims.UserID

	req := &UpdateEmailRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("update email failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	get, err := h.UserService.CheckEmailVerificationCode(c.Request.Context(), req.Email, req.VerifyCode)
	if err != nil {
		logger.AppLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}
	if !get {
		utils.RespondError(c, http.StatusBadRequest, "invalid verification code", nil)
		return
	}

	err = h.UserService.UpdateEmail(c.Request.Context(), req.Email, req.VerifyCode, userID)
	if err != nil {
		logger.AppLogger.Errorf("update email failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "update email failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "email updated successfully")
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Description 修改用户密码
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Tags 用户
// @Param request body UpdatePasswordRequest true "Update password request"
// @Success 200 {string} string "password updated successfully"
// @Failure 401 {object} utils.ResponseError
// @Failure 500 {object} utils.ResponseError
// @Router /api/v1/user/update/password [patch]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", errors.New("invalid token claims"))
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	userID := customClaims.UserID

	// 解析请求体
	var req = &UpdatePasswordRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("update password failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	get, err := h.UserService.CheckEmailVerificationCode(c.Request.Context(), req.Email, req.VerifyCode)
	if err != nil {
		logger.AppLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}
	if !get {
		utils.RespondError(c, http.StatusBadRequest, "invalid verification code", nil)
		return
	}

	// 调用 service 层更新用户信息
	err = h.UserService.UpdatePassword(c.Request.Context(), userID, req.NewPassword, req.OldPassword)
	if err != nil {
		logger.AppLogger.Errorf("update password failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "update password failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "password updated successfully")
	return
}

// ForgetPassword 用户忘记密码, 重置密码接口
// @Summary 忘记密码
// @Description 用户忘记密码, 重置密码接口
// @Accept json
// @Produce json
// @Tags 用户
// @Param request body ForgetPasswordRequest true "Forget password request"
// @Success 200 {string} string "password updated successfully"
// @Failure 400 {object} utils.ResponseError
// @Failure 401 {object} utils.ResponseError
// @Failure 500 {object} utils.ResponseError
// @Router /api/v1/user/forgetPassword [patch]
func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req = ForgetPasswordRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("forget password failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	err := h.UserService.ForgetPassword(c, req.Email, req.VerifyCode, req.NewPassword)
	if err != nil {
		logger.AppLogger.Errorf("forget password failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "forget password failed", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "password updated successfully")
}

// Delete 用户删除账户
// @Summary Delete user
// @Description 删除当前登录用户的账户
// @Accept json
// @Produce json
// @Tags 用户
// @Param request body DeleteRequest true "Delete request"
// @Security ApiKeyAuth
// @Success 200 {string} string "user deleted successfully"
// @Failure 401 {object} utils.ResponseError
// @Failure 500 {object} utils.ResponseError
// @Router /api/v1/user/delete [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", errors.New("invalid token claims"))
		return
	}

	// 将 claims 转换为 CustomClaims 类型
	customClaims, ok := claims.(*appjwt.CustomClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token claims", errors.New("invalid token claims"))
		return
	}

	req := &DeleteRequest{}
	get, err := h.UserService.CheckEmailVerificationCode(c.Request.Context(), req.Email, req.VerifyCode)
	if err != nil {
		logger.AppLogger.Errorf("Register failed: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}
	if !get {
		utils.RespondError(c, http.StatusBadRequest, "invalid verification code", nil)
		return
	}

	userID := customClaims.UserID
	err = h.UserService.Delete(c.Request.Context(), userID)
	if err != nil {
		logger.AppLogger.Errorf("delete failed: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "delete failed", err)
		return
	}
	utils.RespondJSON(c, http.StatusOK, "user deleted successfully")
	return
}
