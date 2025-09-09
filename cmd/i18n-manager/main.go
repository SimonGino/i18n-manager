package main

import (
	"log"
	"os"

	"github.com/SimonGino/i18n-manager/internal/config"
	"github.com/SimonGino/i18n-manager/internal/manager"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "i18n-manager",
		Usage:  "A powerful multilingual properties file management tool for Java project internationalization",
		Action: manager.HandleTranslate, // Default action for translate
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Custom key for translation",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "translate",
				Aliases: []string{"t"},
				Usage:   "Translate text with auto-generated or custom key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "Custom key for translation",
					},
				},
				Action: manager.HandleTranslate,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add manual translations",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "Translation key",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "zh",
						Usage: "Simplified Chinese translation",
					},
					&cli.StringFlag{
						Name:  "en",
						Usage: "English translation",
					},
					&cli.StringFlag{
						Name:  "zh-tw",
						Usage: "Traditional Chinese translation",
					},
				},
				Action: manager.HandleAdd,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List all translation keys or show translations for a specific key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "Show translations for a specific key",
					},
				},
				Action: manager.HandleList,
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check for missing translations",
				Action:  manager.HandleCheck,
			},
			{
				Name:  "config",
				Usage: "Manage configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "set-api-key",
						Usage: "Set API key",
					},
					&cli.StringFlag{
						Name:  "set-api-url",
						Usage: "Set API URL (e.g., https://api.openai.com/v1/chat/completions)",
					},
					&cli.StringFlag{
						Name:  "set-model",
						Usage: "Set model name (e.g., gpt-3.5-turbo, gpt-4, deepseek-chat, qwen-plus)",
					},
					&cli.BoolFlag{
						Name:  "show",
						Usage: "Show current configuration",
					},
				},
				Action: config.HandleConfig,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
