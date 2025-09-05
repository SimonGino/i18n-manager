# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build
```bash
go build -v ./cmd/i18n-manager/main.go
```

### Test
```bash
go test -v ./...
```

### Run locally
```bash
go run cmd/i18n-manager/main.go [command] [args]
```

### Install dependencies
```bash
go mod download
```

### Lint (if available)
```bash
go vet ./...
go fmt ./...
```

## Architecture

This is a CLI tool for managing Java i18n properties files with AI-powered translation capabilities.

### Core Structure
- `cmd/i18n-manager/main.go` - CLI entry point using urfave/cli/v2 framework
- `internal/manager/` - Core translation and file management logic
- `internal/config/` - Configuration management (stores in ~/.config/i18n-manager/)
- `internal/ai/` - OpenAI-compatible API integration for translations

### Key Components
- **CLI Commands**: translate, add, list, check, config
- **Translation Flow**: Input text → AI translation → Properties file updates
- **File Format**: Java properties files with Unicode encoding for Chinese text
- **Configuration**: JSON config with API keys, URLs, models, and language mappings

### Language Configuration
The tool operates on language mappings defined in config:
- Source languages (is_source: true) are used as input for translation
- Target languages are translated to using AI
- File patterns determine properties file naming (e.g., message-application_zh.properties)

### Properties File Management
- Preserves existing file structure and formatting
- Updates existing keys in-place or appends new ones
- Encodes Chinese characters as Unicode escape sequences
- Handles multiple language files simultaneously