package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type PublicHandler struct {
	service service.PublicService
}

type HealthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

type MailRequest struct {
	Mail string `json:"mail"`
}

type MailResponse struct {
	Code string `json:"code"`
}

func NewPublicHandler(service *service.PublicServiceImp) *PublicHandler {
	return &PublicHandler{service: service}
}

// HealthCheck
// @Summary 健康检查
// @Description 检查后端是否正常运行
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /api/v1/health-check [get]
func (h *PublicHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:   "UP",
		Hostname: utils.GetHostName(),
	})
}

// SendVerifyCodeByEmail 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码
// @Accept json
// @Produce json
// @Param request body MailRequest true "Mail request"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/sendVerifyCodeByEmail [post]
func (h *PublicHandler) SendVerifyCodeByEmail(c *gin.Context) {
	var req = MailRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, errors.New("无效的请求"))
		applog.AppLogger.Errorf("send verify code failed: %v", err)
		return
	}

	// 验证邮箱格式
	if !utils.VerifyEmail(req.Mail) {
		utils.RespondJSON(c, http.StatusBadRequest, errors.New("无效的邮箱格式"))
		applog.AppLogger.Errorf("send verify code failed: %v", req.Mail)
		return
	}

	// 发送验证码
	code, err := h.service.SendVerificationCode(c.Request.Context(), req.Mail)
	if err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, errors.New("发送验证码失败"))
		applog.AppLogger.Errorf("send verify code failed: %v", err)
		return
	}

	// 返回成功响应
	utils.RespondJSON(c, http.StatusOK, MailResponse{Code: code})
}
