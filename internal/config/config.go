package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type Config struct {
	APIKey             string   `json:"api_key"`
	DefaultPath        string   `json:"default_path"`
	DefaultSourceLang  string   `json:"default_source_lang"`
	DefaultTargetLangs []string `json:"default_target_langs"`
	AIProvider         string   `json:"ai_provider"`
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
		currentConfig = &Config{
			DefaultSourceLang:  "zh",
			DefaultTargetLangs: []string{"en", "zh_TW"},
			DefaultPath:        ".",
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
