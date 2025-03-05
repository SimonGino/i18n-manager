package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type Input struct {
	Messages []QwenMsg `json:"messages"`
}

type QwenMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Parameters struct {
	ResultFormat string `json:"result_format"`
}

func Translate(req TranslationRequest) (string, error) {
	cfg := config.GetConfig()

	switch cfg.AIProvider {
	case "deepseek":
		return translateWithDeepSeek(req, cfg.APIKey)
	case "qwen":
		return translateWithQwen(req, cfg.APIKey)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", cfg.AIProvider)
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format")
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	return content, nil
}

func translateWithQwen(req TranslationRequest, apiKey string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text from %s to %s. Only return the translated text without any explanation or additional context:\n%s",
		req.SourceLang, req.TargetLang, req.Text)

	qwenReq := QwenRequest{
		Model: "qwen-max",
		Input: Input{
			Messages: []QwenMsg{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
		Parameters: Parameters{
			ResultFormat: "text",
		},
	}

	jsonData, err := json.Marshal(qwenReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/chat/completions", bytes.NewBuffer(jsonData))
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	output, ok := result["output"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	text, ok := output["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text format")
	}

	return text, nil
}
