package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SimonGino/i18n-manager/internal/config"
)

type TranslationRequest struct {
	Text       string
	SourceLang string
	TargetLang string
}

type DeepSeekRequest struct {
	Model    string        `json:"model"`
	Messages []DeepSeekMsg `json:"messages"`
}

type DeepSeekMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenRequest struct {
	Model    string    `json:"model"`
	Messages []QwenMsg `json:"messages"`
}

type QwenMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Parameters struct {
	ResultFormat string `json:"result_format"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type QwenResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func Translate(req TranslationRequest) (string, error) {
	cfg := config.GetConfig()

	// Check if API key is set
	if cfg.APIKey == "" {
		return "", fmt.Errorf("API key not set. Please run:\ni18n-manager config --set-api-key YOUR_API_KEY")
	}

	// Check if AI provider is set
	if cfg.AIProvider == "" {
		return "", fmt.Errorf("AI provider not set. Please run:\ni18n-manager config --set-ai-provider [deepseek|qwen]")
	}

	switch cfg.AIProvider {
	case "deepseek":
		return translateWithDeepSeek(req, cfg.APIKey)
	case "qwen":
		return translateWithQwen(req, cfg.APIKey)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s, please use 'deepseek' or 'qwen'", cfg.AIProvider)
	}
}

func translateWithDeepSeek(req TranslationRequest, apiKey string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text from %s to %s. Only return the translated text without any explanation or additional context:\n%s",
		req.SourceLang, req.TargetLang, req.Text)

	deepSeekReq := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMsg{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(deepSeekReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result DeepSeekResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no translation result in response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func translateWithQwen(req TranslationRequest, apiKey string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text from %s to %s. Only return the translated text without any explanation:\n%s",
		req.SourceLang, req.TargetLang, req.Text)

	qwenReq := QwenRequest{
		Model: "qwen-plus",
		Messages: []QwenMsg{
			{
				Role:    "system",
				Content: "You are a professional translator. Only return the translated text without any explanation.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(qwenReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed (status %d): %s\nPlease check your API key and quota.\nResponse Headers: %v",
			resp.StatusCode, string(body), resp.Header)
	}

	var result QwenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v\nResponse body: %s", err, string(body))
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no translation result in response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
