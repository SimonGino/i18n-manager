# 🌍 i18n-manager

<div align="center">

[![Build Status](https://github.com/SimonGino/i18n-manager/actions/workflows/build.yml/badge.svg)](https://github.com/SimonGino/i18n-manager/actions)
[![Release](https://img.shields.io/github/v/release/SimonGino/i18n-manager?style=flat-square&logo=github&color=blue)](https://github.com/SimonGino/i18n-manager/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square&logo=bookstack)](LICENSE)
[![Stars](https://img.shields.io/github/stars/SimonGino/i18n-manager?style=flat-square&logo=github)](https://github.com/SimonGino/i18n-manager/stargazers)

[English](README.md) | [简体中文](README_CN.md)

</div>

---

A powerful multilingual properties file management tool designed for Java project internationalization (i18n). This tool helps you easily manage and synchronize translations in `message-application.properties` files.

## Features

- 🤖 Smart Translation: Automated text translation using DeepSeek AI or Qwen (Tongyi) AI
- 🔑 Smart Key Generation: Automatically generates keys compliant with Java properties standards
- 🔄 Auto Sync: Automatic synchronization from Simplified Chinese (zh) to Traditional Chinese (zh_TW)
- 📝 Manual Management: Support for manual addition and update of translations
- 🔍 Check Tool: Verification of missing translation entries

## Quick Installation

### Automatic Installation (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/install.sh | bash
```

### Manual Installation

Download the binary for your platform from the [Releases](https://github.com/SimonGino/i18n-manager/releases/latest) page.

### Supported Platforms

- Linux
  - x86_64 (amd64)
  - i386 (32-bit)
  - arm64
  - armv7
- macOS
  - x86_64 (amd64)
  - arm64 (Apple Silicon)
- Windows
  - x86_64 (64-bit)
  - i386 (32-bit)

## Configuration

Before first use, configure your API key:

1. Get your API key:
   - For DeepSeek: Visit [DeepSeek Dashboard](https://platform.deepseek.com/api_keys)
   - For Qwen: Visit [DashScope Console](https://dashscope.console.aliyun.com/apiKey)

2. Configure the API key and provider:
```bash
# Set AI provider (deepseek or qwen)
i18n-manager config --set-ai-provider qwen

# Set API key
i18n-manager config --set-api-key YOUR_API_KEY
```

The configuration file is located at `~/.config/i18n-manager/config.json` (or `%APPDATA%\i18n-manager\config.json` on Windows). Here's an example configuration:

```json
{
  "api_key": "your-api-key",
  "default_path": ".",
  "ai_provider": "deepseek",
  "language": {
    "file_pattern": "message-application%s.properties",
    "mappings": [
      {
        "code": "en",
        "file": "",
        "is_source": false
      },
      {
        "code": "zh",
        "file": "_zh",
        "is_source": true
      },
      {
        "code": "zh_CN",
        "file": "_zh_CN",
        "is_source": true
      },
      {
        "code": "zh_TW",
        "file": "_zh_TW",
        "is_source": false
      }
    ]
  }
}
```

### Configuration Options

- `api_key`: Your AI provider API key
- `default_path`: Default path for properties files
- `ai_provider`: AI provider for translation (currently supports "deepseek" or "qwen")
- `language`: Language configuration
  - `file_pattern`: Pattern for properties files (e.g., "message-application%s.properties")
  - `mappings`: Language mappings
    - `code`: Language code (e.g., "en", "zh", "zh_CN")
    - `file`: File suffix (e.g., "", "_zh", "_zh_CN")
    - `is_source`: Whether this is a source language for translation

You can modify these settings using the following commands:

```bash
i18n-manager config --set-api-key "your-api-key"
i18n-manager config --set-ai-provider "deepseek"
i18n-manager config --show  # Show current configuration
```

## Usage

### 1. Smart Translation

Translate with auto-generated key:

```bash
i18n-manager translate "Text to translate"
```

Translate with custom key:

```bash
i18n-manager translate --key "custom.key.name" "Text to translate"
or
i18n-manager translate -k "custom.key.name" "Text to translate"
```

### 2. Manual Translation Addition

Add complete multilingual translations:

```bash
i18n-manager add \
    --key "custom.key.name" \
    --zh "简体中文" \
    --en "English" \
    --zh-tw "繁體中文"
```

Or add partial translations:

```bash
i18n-manager add \
    --key "custom.key.name" \
    --zh "简体中文" \
    --en "English"
```

### 3. View and Check

List all translation keys:

```bash
i18n-manager list
```

List translations for a specific key:

```bash
i18n-manager list --key "error.skill.unavailable"
# or use the short flag
i18n-manager list -k "error.skill.unavailable"
```

Check for missing translations:

```bash
i18n-manager check
```

### 4. Configuration Management

Set API key:

```bash
i18n-manager config --set-api-key YOUR_API_KEY
```

View current configuration:

```bash
i18n-manager config --show
```

## Configuration File

Configuration files are located at:

- Linux/macOS: `~/.config/i18n-manager/config.json`
- Windows: `%APPDATA%\i18n-manager\config.json`

Contains the following settings:

- `api_key`: API key
- `default_path`: Default working directory
- `default_source_lang`: Default source language
- `default_target_langs`: Default target language list

## Key Naming Convention

Generated keys follow these conventions:

- Use lowercase letters, numbers, and dots (.)
- Use dots (.) as hierarchy separators
- Use common prefixes for categorization:
  - `error.` - Error messages
  - `success.` - Success messages
  - `info.` - Information prompts
  - `label.` - UI labels
  - `button.` - Button text
  - `title.` - Page/section titles
  - `msg.` - General messages
  - `validation.` - Validation messages

Examples:

- `validation.username.notEmpty`
- `success.data.saved`
- `button.submit.text`

## File Structure

The tool manages the following files in the specified directory:

- `message-application.properties` - English translations
- `message-application_zh.properties` - Simplified Chinese translations
- `message-application_zh_TW.properties` - Traditional Chinese translations

## Uninstallation

To uninstall the tool:

### Automatic Uninstallation (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/uninstall.sh -o uninstall.sh
chmod +x uninstall.sh
./uninstall.sh
```

### Manual Uninstallation

1. Remove executable

```bash
# Linux/macOS
rm $(which i18n-manager)

# Windows
del C:\path\to\i18n-manager.exe
```

2. Remove configuration files (optional)

```bash
# Linux/macOS
rm -rf ~/.config/i18n-manager

# Windows
rd /s /q %APPDATA%\i18n-manager
```

## Development

1. Clone the project:

```bash
git clone https://github.com/yourusername/i18n-manager.git
```

2. Create the cmd directory structure:

```bash
mkdir -p cmd/i18n-manager
```

3. Install dependencies:

```bash
go mod download
```

4. Run tests:

```bash
go test ./...
```

## Contributing

Pull requests and issue reports are welcome!

## License

MIT License