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
	DefaultPath string         `json:"default_path"`
	AIProvider  string         `json:"ai_provider"`
	Language    LanguageConfig `json:"language"`
}

var currentConfig *Config

func init() {
	loadConfig()
}

func getConfigPath() string {
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

	if provider := c.String("set-ai-provider"); provider != "" {
		if provider != "deepseek" && provider != "qwen" {
			return fmt.Errorf("invalid AI provider, must be either 'deepseek' or 'qwen'")
		}
		currentConfig.AIProvider = provider
		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save AI provider: %v", err)
		}
		fmt.Println("AI provider updated successfully")
		return nil
	}

	if c.Bool("show") {
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
