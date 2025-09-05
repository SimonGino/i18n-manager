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

	translations := make(map[string]string)

	// Get source language configuration
	sourceLang := config.GetSourceLang()
	if sourceLang == nil {
		return fmt.Errorf("no source language configured")
	}

	// Get target language configuration
	targetLangs := config.GetTargetLangs()
	if len(targetLangs) == 0 {
		return fmt.Errorf("no target languages configured")
	}

	// Save source language text
	translations[sourceLang.Code] = text

	// If source language is zh, also save for zh_CN
	if sourceLang.Code == "zh" {
		translations["zh_CN"] = text
	}

	// If no key provided, translate to English first for key generation
	if key == "" {
		// Get English translation for key generation
		englishText, err := ai.Translate(ai.TranslationRequest{
			Text:       text,
			SourceLang: sourceLang.Code,
			TargetLang: "en",
		})
		if err != nil {
			return fmt.Errorf("failed to generate key: %v", err)
		}
		translations["en"] = englishText
		key = generateKey(englishText)
	}

	// Translate to other target languages
	for _, targetLang := range targetLangs {
		if targetLang.Code == "en" && translations["en"] != "" {
			continue // Skip if English translation already exists
		}
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

	// Print translations to be added
	fmt.Printf("\nTranslations to be added:\n")
	fmt.Printf("Key: %s\n", key)
	for lang, value := range translations {
		fmt.Printf("%s: %s\n", lang, value)
	}

	// Ask for confirmation
	fmt.Print("\nDo you want to add these translations? (y/N): ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// 如果用户直接按回车，Scanln 会返回错误，这种情况我们视为取消操作
		fmt.Println("Translation cancelled")
		return nil
	}

	if strings.ToLower(response) != "y" {
		fmt.Println("Translation cancelled")
		return nil
	}

	if err := saveTranslations(key, translations); err != nil {
		return fmt.Errorf("error saving translations: %v", err)
	}

	fmt.Printf("Successfully added translations with key: %s\n", key)
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

	// If key is specified, show only that key's translations
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

	// Show all translations
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
	// 将文本转换为小写
	key := strings.ToLower(text)

	// 移除标点符号和特殊字符
	key = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			return r
		}
		return -1
	}, key)

	// 将空格替换为点号
	key = strings.ReplaceAll(key, " ", ".")

	// 处理连续的点号
	for strings.Contains(key, "..") {
		key = strings.ReplaceAll(key, "..", ".")
	}

	// 移除开头和结尾的点号
	key = strings.Trim(key, ".")

	// 限制长度
	if len(key) > 50 {
		key = key[:50]
		// 确保不以点号结尾
		key = strings.TrimRight(key, ".")
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

// Convert Chinese characters to Unicode escape sequences
func encodeToUnicode(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r > 127 {
			result.WriteString(fmt.Sprintf("\\u%04x", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func saveTranslations(key string, translations map[string]string) error {
	cfg := config.GetConfig()

	// Process each configured language mapping
	for _, mapping := range cfg.Language.Mappings {
		filename := config.GetPropertiesFilePath(mapping.Code)
		value, exists := translations[mapping.Code]

		// Skip if no translation provided for this language
		if !exists {
			continue
		}

		// Read existing file content while preserving order
		var lines []string
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error opening %s: %v", filename, err)
		}

		scanner := bufio.NewScanner(file)
		keyFound := false
		for scanner.Scan() {
			line := scanner.Text() // Don't trim to preserve formatting
			if line == "" || strings.HasPrefix(line, "#") {
				lines = append(lines, line)
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				lines = append(lines, line)
				continue
			}

			currentKey := strings.TrimSpace(parts[0])
			if currentKey == key {
				// Encode value if it's Chinese
				if strings.Contains(mapping.Code, "zh") {
					value = encodeToUnicode(value)
				}
				lines = append(lines, fmt.Sprintf("%s=%s", currentKey, value))
				keyFound = true
			} else {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			file.Close()
			return fmt.Errorf("error reading %s: %v", filename, err)
		}

		// If key not found, append it to the end
		if !keyFound {
			// Encode value if it's Chinese
			if strings.Contains(mapping.Code, "zh") {
				value = encodeToUnicode(value)
			}
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		}

		// Rewrite the file
		if err := file.Truncate(0); err != nil {
			file.Close()
			return fmt.Errorf("error truncating %s: %v", filename, err)
		}

		if _, err := file.Seek(0, 0); err != nil {
			file.Close()
			return fmt.Errorf("error seeking in %s: %v", filename, err)
		}

		writer := bufio.NewWriter(file)
		for i, line := range lines {
			if i > 0 || line != "" { // Skip initial empty line
				if i < len(lines)-1 {
					if _, err := fmt.Fprintln(writer, line); err != nil {
						file.Close()
						return fmt.Errorf("error writing to %s: %v", filename, err)
					}
				} else {
					// 最后一行不添加换行符
					if _, err := fmt.Fprint(writer, line); err != nil {
						file.Close()
						return fmt.Errorf("error writing to %s: %v", filename, err)
					}
				}
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
