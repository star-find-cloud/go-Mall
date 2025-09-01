package ds

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	"github.com/star-find-cloud/star-mall/domain"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"io"
	"net/http"
	"strings"
)

type DeepSeekAdapter struct {
	client *http.Client
}

// DeepSeek 结构体
// @Description: ds结构体
type DeepSeek struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      bool    `json:"stream"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// DeepSeekResponse 结构体
type DeepSeekResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage"`
}

// Choice 消息选择项
type Choice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Usage Token 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 消息内容
type Message struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}

func NewDeepseekClient() *DeepSeekAdapter {
	client := &http.Client{}

	return &DeepSeekAdapter{client: client}
}

func (a *DeepSeekAdapter) GenerateUserMsg(products []domain.Product, msg interface{}) string {
	userMsg := fmt.Sprintf("用户近期常看的商品标签为: %v, 这些标签的近期热门商品为: %v, 你需要根据用户常看标签, 从近期热门商品中选择5个进行推荐", msg, products)
	return userMsg
}

// GenerateRequest 生成请求
func (a *DeepSeekAdapter) GenerateRequest() *http.Client {
	c := &http.Client{}
	return c
}

// GenerateReplyStreamChat 生成回复流式聊天回答
func (a *DeepSeekAdapter) GenerateReplyStreamChat(ctx context.Context, msg string, temp float64, model string) (<-chan string, error) {
	reqBody := DeepSeek{
		Model: model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: "你是一个购物平台商品推荐助手, 你需要根据用户提供的信息来推荐商品"}, {Role: "user", Content: msg},
		},
		Stream:      true,
		MaxTokens:   2048,
		Temperature: temp,
	}
	jsonBody, _ := json.Marshal(reqBody)

	var c = conf.GetConfig()
	apiKey := c.AI.Deepseek.ApiKey
	if apiKey == "" {
		applog.AppLogger.Errorf("DeepseekClient faild, err: %v", "apiKey is empty")
		return nil, fmt.Errorf("apiKey is empty")
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		applog.AppLogger.Errorf("DeepseekClient failed, err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-6e97edeed6694c65ad99fd1c3653f9ee")

	// 发送请求
	resp, err := a.client.Do(req)
	if err != nil {
		applog.AppLogger.Errorf("DeepseekClient failed, err: %v", err)
		return nil, err
	}

	//if resp.Header.Get("Content-Type") != "text/event-stream" {
	//	body, _ := io.ReadAll(resp.Body)
	//	applog.AppLogger.Errorf("Invalid content type: %s, body: %s",
	//		resp.Header.Get("Content-Type"), string(body))
	//	return nil, fmt.Errorf("invalid content type")
	//}
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		applog.AppLogger.Errorf("DeepseekClient bad status: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// 创建内容通道
	contentChan := make(chan string)

	// 启动协程处理流式响应
	// 启动协程处理流式响应
	go func() {
		defer resp.Body.Close()
		defer close(contentChan)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)
		//var jsonBuffer string

		for scanner.Scan() {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				applog.AppLogger.Info("Context canceled, stopping stream processing")
				return
			default:
			}

			line := scanner.Text()
			// 过滤心跳包和空行
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// 处理数据
			if line == "data: [DONE]" {
				applog.AppLogger.Debug("Stream completed with [DONE] message")
				return
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			//// 剥离 sse 前缀
			//if strings.HasPrefix(line, "date:") {
			//	line = strings.TrimPrefix(line, "date:")
			//}
			jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

			// 解析 JSON
			var response DeepSeekResponse
			if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
				applog.AppLogger.Errorf("Failed to parse JSON: %v, data: %s", err, jsonData)
				continue
			}
			//jsonBuffer += line
			//
			//// 检查是否为 json
			//if !json.Valid([]byte(jsonBuffer)) {
			//	continue
			//}
			//
			//// 解析 json
			//var chunk DeepSeekResponse
			//if err := json.Unmarshal([]byte(jsonBuffer), &chunk); err != nil {
			//	applog.AppLogger.Errorf("JSON parsing error: %v, raw: %s", err, jsonBuffer)
			//	jsonBuffer = ""
			//	continue
			//}
			//jsonBuffer = ""

			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				if choice.Delta.Content != "" {
					//applog.AppLogger.Infof("Content: %s", choice.Delta.Content)
					contentChan <- choice.Delta.Content
				}
				if choice.FinishReason == "stop" {
					applog.AppLogger.Debugln("Stream completed with finish_reason: stop")
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			applog.AppLogger.Errorf("Stream reading error: %v", err)
		}
	}()

	return contentChan, nil
}

// GenerateReplyStreamChat 生成回复流式聊天回答
func (a *DeepSeekAdapter) GenerateReplyStreamReasoner(ctx context.Context, msg string, temp float64, model string) (<-chan string, error) {
	reqBody := DeepSeek{
		Model: model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: msg},
		},
		Stream:      true,
		MaxTokens:   2048,
		Temperature: temp,
	}
	jsonBody, _ := json.Marshal(reqBody)

	var c = conf.GetConfig()
	apiKey := c.AI.Deepseek.ApiKey
	if apiKey == "" {
		applog.AppLogger.Errorf("DeepseekClient faild, err: %v", "apiKey is empty")
		return nil, fmt.Errorf("apiKey is empty")
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		applog.AppLogger.Errorf("DeepseekClient failed, err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-6e97edeed6694c65ad99fd1c3653f9ee")

	// 发送请求
	resp, err := a.client.Do(req)
	if err != nil {
		applog.AppLogger.Errorf("DeepseekClient failed, err: %v", err)
		return nil, err
	}

	//if resp.Header.Get("Content-Type") != "text/event-stream" {
	//	body, _ := io.ReadAll(resp.Body)
	//	applog.AppLogger.Errorf("Invalid content type: %s, body: %s",
	//		resp.Header.Get("Content-Type"), string(body))
	//	return nil, fmt.Errorf("invalid content type")
	//}
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		applog.AppLogger.Errorf("DeepseekClient bad status: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// 创建内容通道
	contentChan := make(chan string)

	// 启动协程处理流式响应
	go func() {
		defer resp.Body.Close()
		defer close(contentChan)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)
		//var jsonBuffer string

		for scanner.Scan() {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				applog.AppLogger.Info("Context canceled, stopping stream processing")
				return
			default:
			}

			line := scanner.Text()
			// 过滤心跳包和空行
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// 处理数据
			if line == "data: [DONE]" {
				applog.AppLogger.Debug("Stream completed with [DONE] message")
				return
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			//// 剥离 sse 前缀
			//if strings.HasPrefix(line, "date:") {
			//	line = strings.TrimPrefix(line, "date:")
			//}
			jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

			// 解析 JSON
			var response DeepSeekResponse
			if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
				applog.AppLogger.Errorf("Failed to parse JSON: %v, data: %s", err, jsonData)
				continue
			}
			//jsonBuffer += line
			//
			//// 检查是否为 json
			//if !json.Valid([]byte(jsonBuffer)) {
			//	continue
			//}
			//
			//// 解析 json
			//var chunk DeepSeekResponse
			//if err := json.Unmarshal([]byte(jsonBuffer), &chunk); err != nil {
			//	applog.AppLogger.Errorf("JSON parsing error: %v, raw: %s", err, jsonBuffer)
			//	jsonBuffer = ""
			//	continue
			//}
			//jsonBuffer = ""

			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				if choice.Delta.Content != "" {
					//applog.AppLogger.Infof("Content: %s", choice.Delta.Content)
					contentChan <- choice.Delta.Content
				}
				if choice.FinishReason == "stop" {
					applog.AppLogger.Debugln("Stream completed with finish_reason: stop")
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			applog.AppLogger.Errorf("Stream reading error: %v", err)
		}
	}()

	return contentChan, nil
}
