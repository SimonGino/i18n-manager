package manager

import (
	"bufio"
	"fmt"
	"os"
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

	cfg := config.GetConfig()
	translations := make(map[string]string)

	// 从中文翻译到其他语言
	for _, targetLang := range cfg.DefaultTargetLangs {
		if targetLang == cfg.DefaultSourceLang {
			translations[targetLang] = text
			continue
		}

		translated, err := ai.Translate(ai.TranslationRequest{
			Text:       text,
			SourceLang: cfg.DefaultSourceLang,
			TargetLang: targetLang,
		})
		if err != nil {
			return fmt.Errorf("error translating to %s: %v", targetLang, err)
		}
		translations[targetLang] = translated
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
	if zh := c.String("zh"); zh != "" {
		translations["zh"] = zh
	}
	if en := c.String("en"); en != "" {
		translations["en"] = en
	}
	if zhTW := c.String("zh-tw"); zhTW != "" {
		translations["zh_TW"] = zhTW
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

func HandleList(c *cli.Context) error {
	translations, err := loadAllTranslations()
	if err != nil {
		return fmt.Errorf("error loading translations: %v", err)
	}

	for _, t := range translations {
		fmt.Printf("Key: %s\n", t.Key)
		for lang, value := range t.Values {
			fmt.Printf("  %s: %s\n", lang, value)
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
	allLangs := append([]string{cfg.DefaultSourceLang}, cfg.DefaultTargetLangs...)

	var missingCount int
	for _, t := range translations {
		for _, lang := range allLangs {
			if _, ok := t.Values[lang]; !ok {
				fmt.Printf("Missing translation for key '%s' in language '%s'\n", t.Key, lang)
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
	// 简单的key生成逻辑，可以根据需要扩展
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

func getPropertiesFilePath(lang string) string {
	filename := "message-application"
	if lang != "" {
		filename += "_" + lang
	}
	filename += ".properties"
	return filename
}

func loadAllTranslations() ([]Translation, error) {
	cfg := config.GetConfig()
	allLangs := append([]string{cfg.DefaultSourceLang}, cfg.DefaultTargetLangs...)

	translations := make(map[string]*Translation)

	for _, lang := range allLangs {
		filename := getPropertiesFilePath(lang)
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
			translations[key].Values[lang] = value
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
	for lang, value := range translations {
		filename := getPropertiesFilePath(lang)

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
