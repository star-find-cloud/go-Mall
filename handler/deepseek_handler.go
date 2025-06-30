package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/domain"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/service"
	"github.com/star-find-cloud/star-mall/utils"
	"net/http"
)

type DeepseekHandler struct {
	service service.DeepseekService
}

func NewDeepseekHandler(service service.DeepseekService) *DeepseekHandler {
	return &DeepseekHandler{service: service}
}

// SuggestProductRequest 商品推荐请求体
// @Description 商品推荐请求体
type SuggestProductRequest struct {
	// @Description 商品搜索关键词(用户登录请求ai接口忽略该参数)
	Msg string `json:"msg"`
	// @Description 是否启用 r1 模型, true: 使用 v3 模型, false: 使用 r1 深度思考模型
	IsChat bool `json:"isChat"`
}

// SuggestProductByUserTags 根据用户标签推荐商品
// @Summary 根据用户标签推荐商品
// @Description 根据用户标签推荐商品(sse协议)
// @Tags deepseek
// @Accept  json
// @Produce  text/event-stream
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body SuggestProductRequest true "请求体"
// @Success 200 {string} string "data: <message>\n\n"
// @Failure 400 {object} string "deepseek 生成信息失败"
// @Failure 401 {object} string "invalid token claims"
// @Failure 404 {object} string "请求错误"
// @Failure 500 {object} string "解析请求体失败"
// @Router /api/v1/deepseek/suggest [post]
func (h *DeepseekHandler) SuggestProductByUserTags(c *gin.Context) {
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

	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not merchant", errors.New("not merchant"))
		return
	}

	userID := customClaims.UserID

	var req = &SuggestProductRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("解析请求体失败: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	var (
		err         error
		messageChan <-chan string
	)

	if req.IsChat {
		messageChan, err = h.service.SuggestProductByUserTags(c.Request.Context(), userID, domain.ModelResearcher)
	} else if !req.IsChat {
		messageChan, err = h.service.SuggestProductByUserTags(c.Request.Context(), userID, domain.ModelChat)
	} else {
		utils.RespondError(c, http.StatusNotFound, "请求错误", errors.New("请求错误"))
		return
	}
	if err != nil {
		logger.AppLogger.Errorf("deepseek 生成信息失败: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "deepseek 生成信息失败", err)
		return
	} else {
		utils.RespondJSON(c, http.StatusOK, "deepseek 生成信息成功")
	}

	// 使用 SSE 进行流式响应
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// 阻止自动断连
	c.Writer.Flush()

	// 监听客户端是否断连
	clientGone := c.Request.Context().Done()

	// 从 messageChan 中读取数据并发送给客户端
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				// 通道已关闭，发送完成事件
				_, err = fmt.Fprintf(c.Writer, "event: complete\n\n")
				if err != nil {
					logger.AppLogger.Errorf("发送完成事件失败: %v", err)
				}
				c.Writer.Flush()
				return
			}
			// 发送数据
			_, err = fmt.Fprintf(c.Writer, "%s", msg)
			if err != nil {
				logger.AppLogger.Errorf("发送数据失败: %v", err)
			}
			c.Writer.Flush()
		case <-clientGone:
			// 客户端断开连接
			logger.AppLogger.Info("Client closed connection")
			return
		}
	}
}

// SuggestProductBySearch 根据搜索关键词推荐商品
// @Summary 根据搜索关键词推荐商品
// @Description 根据搜索关键词推荐商品(sse 响应)
// @Tags deepseek
// @Accept  json
// @Produce  text/event-stream
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body SuggestProductRequest true "请求体"
// @Success 200 {string} string "data: <message>\n\n"
// @Failure 400 {object} string "deepseek 生成信息失败"
// @Failure 401 {object} string "invalid token claims"
// @Failure 404 {object} string "请求错误"
// @Failure 500 {object} string "解析请求体失败"
// @Router /api/v1/deepseek/search [post]
func (h *DeepseekHandler) SuggestProductBySearch(c *gin.Context) {
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

	if customClaims.Roles == 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	if customClaims.UserID <= 0 {
		utils.RespondError(c, http.StatusUnauthorized, "not user", errors.New("not user"))
		return
	}

	var req = &SuggestProductRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.HttpLogger.Errorf("解析请求体失败: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	if req.Msg == "" {
		utils.RespondError(c, http.StatusBadRequest, "Msg can't nil", errors.New("msg can't nil"))
		return
	}

	var (
		err         error
		messageChan <-chan string
	)

	if req.IsChat {
		messageChan, err = h.service.SearchProductSuggest(c.Request.Context(), req.Msg, domain.ModelResearcher)
	} else if !req.IsChat {
		messageChan, err = h.service.SearchProductSuggest(c.Request.Context(), req.Msg, domain.ModelChat)
	} else {
		utils.RespondError(c, http.StatusNotFound, "请求错误", errors.New("请求错误"))
		return
	}

	if err != nil {
		logger.AppLogger.Errorf("deepseek 生成信息失败: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "deepseek 生成信息失败", err)
		return
	} else {
		utils.RespondJSON(c, http.StatusOK, "deepseek 生成信息成功")
	}

	// 使用 SSE 进行流式响应
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// 阻止自动断连
	c.Writer.Flush()

	// 监听客户端是否断连
	clientGone := c.Request.Context().Done()

	// 从 messageChan 中读取数据并发送给客户端
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				// 通道已关闭，发送完成事件
				_, err = fmt.Fprintf(c.Writer, "event: complete\n\n")
				if err != nil {
					logger.AppLogger.Errorf("发送完成事件失败: %v", err)
				}
				c.Writer.Flush()
				return
			}
			// 发送数据
			_, err = fmt.Fprintf(c.Writer, "%s", msg)
			if err != nil {
				logger.AppLogger.Errorf("发送数据失败: %v", err)
			}
			c.Writer.Flush()
		case <-clientGone:
			// 客户端断开连接
			logger.AppLogger.Info("Client closed connection")
			return
		}
	}
}
