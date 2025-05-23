package ai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SimonGino/i18n-manager/internal/config"
	"github.com/sashabaranov/go-openai"
)

type TranslationRequest struct {
	Text       string
	SourceLang string
	TargetLang string
}

func Translate(req TranslationRequest) (string, error) {
	cfg := config.GetConfig()

	// 检查API密钥是否设置
	if cfg.APIKey == "" {
		return "", fmt.Errorf("API密钥未设置。请运行:\ni18n-manager config --set-api-key YOUR_API_KEY")
	}

	// 检查API URL是否设置
	if cfg.APIURL == "" {
		return "", fmt.Errorf("API URL未设置。请运行:\ni18n-manager config --set-api-url YOUR_API_URL")
	}

	// 检查模型是否设置
	if cfg.Model == "" {
		return "", fmt.Errorf("AI模型未设置。请运行:\ni18n-manager config --set-model MODEL_NAME")
	}

	// 创建自定义配置
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	clientConfig.BaseURL = cfg.APIURL
	clientConfig.HTTPClient = &http.Client{Timeout: 30 * time.Second}

	// 创建OpenAI客户端
	client := openai.NewClientWithConfig(clientConfig)

	// 构建提示信息
	prompt := fmt.Sprintf("将以下文本从%s翻译为%s。只返回翻译后的文本，不要包含任何解释或额外内容：\n%s",
		req.SourceLang, req.TargetLang, req.Text)

	// 创建请求
	request := openai.ChatCompletionRequest{
		Model: cfg.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一位专业翻译。只返回翻译后的文本，不要包含任何解释。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3, // 较低的温度使输出更确定
	}

	// 发送请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %v\n请检查您的API密钥、配额和网络连接。", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("响应中没有翻译结果")
	}

	// 返回翻译结果
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
