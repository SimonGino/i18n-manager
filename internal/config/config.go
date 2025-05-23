package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type LanguageConfig struct {
	FilePattern string        `json:"file_pattern"` // 文件名模式，如 "message-application%s.properties"
	Default     string        `json:"default"`      // 默认语言文件（不带后缀的文件名）
	Mappings    []LangMapping `json:"mappings"`
}

type LangMapping struct {
	Code     string `json:"code"`      // 语言代码，如 "en", "zh", "zh_CN", "zh_TW"
	File     string `json:"file"`      // 对应的文件后缀，如 "", "_zh", "_zh_CN", "_zh_TW"
	IsSource bool   `json:"is_source"` // 是否为源语言
}

type Config struct {
	APIKey      string         `json:"api_key"`
	APIURL      string         `json:"api_url"`
	Model       string         `json:"model"`
	DefaultPath string         `json:"default_path"`
	Language    LanguageConfig `json:"language"`
}

var currentConfig *Config

func init() {
	loadConfig()
}

func getConfigPath() string {
	// 优先使用 ~/.config 目录
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, ".config")
		// 检查 ~/.config 目录是否存在
		if _, err := os.Stat(configDir); err == nil {
			return filepath.Join(configDir, "i18n-manager")
		}
	}

	// 如果 ~/.config 不可用，则使用系统默认的配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", ".config", "i18n-manager")
	}
	return filepath.Join(configDir, "i18n-manager")
}

func getConfigFilePath() string {
	return filepath.Join(getConfigPath(), "config.json")
}

func loadConfig() {
	configPath := getConfigPath()
	if err := os.MkdirAll(configPath, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		return
	}

	configFile := getConfigFilePath()
	data, err := os.ReadFile(configFile)
	if err != nil {
		// 默认配置
		currentConfig = &Config{
			DefaultPath: ".",
			APIURL:      "https://api.openai.com/v1/chat/completions", // 默认使用OpenAI API完整路径
			Model:       "gpt-3.5-turbo",                              // 默认模型
			Language: LanguageConfig{
				FilePattern: "message-application%s.properties",
				Default:     "", // 英文文件没有后缀
				Mappings: []LangMapping{
					{
						Code:     "en",
						File:     "",
						IsSource: false,
					},
					{
						Code:     "zh",
						File:     "_zh",
						IsSource: true,
					},
					{
						Code:     "zh_TW",
						File:     "_zh_TW",
						IsSource: false,
					},
				},
			},
		}
		return
	}

	currentConfig = &Config{}
	if err := json.Unmarshal(data, currentConfig); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}
}

func saveConfig() error {
	data, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	configFile := getConfigFilePath()
	return os.WriteFile(configFile, data, 0644)
}

func HandleConfig(c *cli.Context) error {
	if apiKey := c.String("set-api-key"); apiKey != "" {
		currentConfig.APIKey = apiKey
		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save API key: %v", err)
		}
		fmt.Println("API key updated successfully")
		return nil
	}

	if apiURL := c.String("set-api-url"); apiURL != "" {
		currentConfig.APIURL = apiURL
		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save API URL: %v", err)
		}
		fmt.Println("API URL updated successfully")
		return nil
	}

	if model := c.String("set-model"); model != "" {
		currentConfig.Model = model
		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save model: %v", err)
		}
		fmt.Println("Model updated successfully")
		return nil
	}

	if c.Bool("show") {
		configPath := getConfigFilePath()
		fmt.Printf("Config file path: %s\n\n", configPath)
		data, err := json.MarshalIndent(currentConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("error formatting config: %v", err)
		}
		fmt.Println(string(data))
		return nil
	}

	return fmt.Errorf("no valid config action specified")
}

func GetConfig() *Config {
	return currentConfig
}

// 获取源语言配置
func GetSourceLang() *LangMapping {
	for _, mapping := range currentConfig.Language.Mappings {
		if mapping.IsSource {
			return &mapping
		}
	}
	return nil
}

// 获取目标语言配置列表
func GetTargetLangs() []LangMapping {
	var targets []LangMapping
	for _, mapping := range currentConfig.Language.Mappings {
		if !mapping.IsSource {
			targets = append(targets, mapping)
		}
	}
	return targets
}

// 根据语言代码获取文件名
func GetPropertiesFilePath(lang string) string {
	var suffix string
	for _, mapping := range currentConfig.Language.Mappings {
		if mapping.Code == lang {
			suffix = mapping.File
			break
		}
	}
	return fmt.Sprintf(currentConfig.Language.FilePattern, suffix)
}
