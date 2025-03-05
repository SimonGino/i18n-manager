package manager

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SimonGino/i18n-manager/internal/ai"
	"github.com/SimonGino/i18n-manager/internal/config"
	"github.com/urfave/cli/v2"
)

type Translation struct {
	Key    string
	Values map[string]string
}

func HandleTranslate(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("please provide the text to translate")
	}

	text := c.Args().First()
	key := c.String("key")
	if key == "" {
		key = generateKey(text)
	}

	translations := make(map[string]string)

	// 获取源语言配置
	sourceLang := config.GetSourceLang()
	if sourceLang == nil {
		return fmt.Errorf("no source language configured")
	}

	// 获取目标语言配置
	targetLangs := config.GetTargetLangs()
	if len(targetLangs) == 0 {
		return fmt.Errorf("no target languages configured")
	}

	// 保存源语言文本
	translations[sourceLang.Code] = text

	// 翻译到目标语言
	for _, targetLang := range targetLangs {
		translated, err := ai.Translate(ai.TranslationRequest{
			Text:       text,
			SourceLang: sourceLang.Code,
			TargetLang: targetLang.Code,
		})
		if err != nil {
			return fmt.Errorf("error translating to %s: %v", targetLang.Code, err)
		}
		translations[targetLang.Code] = translated
	}

	if err := saveTranslations(key, translations); err != nil {
		return fmt.Errorf("error saving translations: %v", err)
	}

	fmt.Printf("Successfully added translations for key: %s\n", key)
	return nil
}

func HandleAdd(c *cli.Context) error {
	key := c.String("key")
	if key == "" {
		return fmt.Errorf("key is required")
	}

	translations := make(map[string]string)

	// 从命令行参数获取各语言的翻译
	for _, mapping := range config.GetConfig().Language.Mappings {
		if value := c.String(mapping.Code); value != "" {
			translations[mapping.Code] = value
		}
	}

	if len(translations) == 0 {
		return fmt.Errorf("at least one translation is required")
	}

	if err := saveTranslations(key, translations); err != nil {
		return fmt.Errorf("error saving translations: %v", err)
	}

	fmt.Printf("Successfully added translations for key: %s\n", key)
	return nil
}

func decodeUnicode(s string) string {
	var result string
	for len(s) > 0 {
		if strings.HasPrefix(s, "\\u") && len(s) >= 6 {
			code, err := strconv.ParseInt(s[2:6], 16, 32)
			if err == nil {
				result += string(rune(code))
				s = s[6:]
				continue
			}
		}
		result += string(s[0])
		s = s[1:]
	}
	return result
}

func HandleList(c *cli.Context) error {
	translations, err := loadAllTranslations()
	if err != nil {
		return fmt.Errorf("error loading translations: %v", err)
	}

	// 如果指定了key，只显示该key的翻译
	if key := c.String("key"); key != "" {
		found := false
		for _, t := range translations {
			if t.Key == key {
				fmt.Printf("Key: %s\n", t.Key)
				for lang, value := range t.Values {
					decodedValue := decodeUnicode(value)
					fmt.Printf("  %s: %s\n", lang, decodedValue)
				}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("key '%s' not found", key)
		}
		return nil
	}

	// 显示所有翻译
	for _, t := range translations {
		fmt.Printf("Key: %s\n", t.Key)
		for lang, value := range t.Values {
			decodedValue := decodeUnicode(value)
			fmt.Printf("  %s: %s\n", lang, decodedValue)
		}
		fmt.Println()
	}

	return nil
}

func HandleCheck(c *cli.Context) error {
	translations, err := loadAllTranslations()
	if err != nil {
		return fmt.Errorf("error loading translations: %v", err)
	}

	cfg := config.GetConfig()
	var missingCount int

	for _, t := range translations {
		for _, mapping := range cfg.Language.Mappings {
			if _, ok := t.Values[mapping.Code]; !ok {
				fmt.Printf("Missing translation for key '%s' in language '%s'\n", t.Key, mapping.Code)
				missingCount++
			}
		}
	}

	if missingCount == 0 {
		fmt.Println("All translations are complete!")
	} else {
		fmt.Printf("Found %d missing translations\n", missingCount)
	}

	return nil
}

func generateKey(text string) string {
	key := strings.ToLower(text)
	key = strings.ReplaceAll(key, " ", ".")
	key = strings.ReplaceAll(key, "'", "")
	key = strings.ReplaceAll(key, "\"", "")
	key = strings.ReplaceAll(key, ":", "")
	key = strings.ReplaceAll(key, "?", "")
	key = strings.ReplaceAll(key, "!", "")
	key = strings.ReplaceAll(key, ",", "")
	key = strings.ReplaceAll(key, ".", ".")

	if len(key) > 50 {
		key = key[:50]
	}

	return "msg." + key
}

func loadAllTranslations() ([]Translation, error) {
	cfg := config.GetConfig()
	translations := make(map[string]*Translation)

	// 遍历所有语言文件
	for _, mapping := range cfg.Language.Mappings {
		filename := config.GetPropertiesFilePath(mapping.Code)
		file, err := os.Open(filename)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %v", filename, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if _, ok := translations[key]; !ok {
				translations[key] = &Translation{
					Key:    key,
					Values: make(map[string]string),
				}
			}
			translations[key].Values[mapping.Code] = value
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading %s: %v", filename, err)
		}
	}

	result := make([]Translation, 0, len(translations))
	for _, t := range translations {
		result = append(result, *t)
	}
	return result, nil
}

func saveTranslations(key string, translations map[string]string) error {
	// 保存到每个语言对应的文件
	for lang, value := range translations {
		filename := config.GetPropertiesFilePath(lang)

		// 读取现有文件
		existingContent := make(map[string]string)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error opening %s: %v", filename, err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			existingContent[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			return fmt.Errorf("error reading %s: %v", filename, err)
		}

		// 更新或添加新的翻译
		existingContent[key] = value

		// 重写文件
		if err := file.Truncate(0); err != nil {
			file.Close()
			return fmt.Errorf("error truncating %s: %v", filename, err)
		}

		if _, err := file.Seek(0, 0); err != nil {
			file.Close()
			return fmt.Errorf("error seeking in %s: %v", filename, err)
		}

		writer := bufio.NewWriter(file)
		for k, v := range existingContent {
			if _, err := fmt.Fprintf(writer, "%s=%s\n", k, v); err != nil {
				file.Close()
				return fmt.Errorf("error writing to %s: %v", filename, err)
			}
		}

		if err := writer.Flush(); err != nil {
			file.Close()
			return fmt.Errorf("error flushing %s: %v", filename, err)
		}

		file.Close()
	}

	return nil
}
