# 🌍 i18n-manager

<div align="center">

[![Build Status](https://github.com/SimonGino/i18n-manager/actions/workflows/build.yml/badge.svg)](https://github.com/SimonGino/i18n-manager/actions)
[![Release](https://img.shields.io/github/v/release/SimonGino/i18n-manager?style=flat-square&logo=github&color=blue)](https://github.com/SimonGino/i18n-manager/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square&logo=bookstack)](LICENSE)
[![Stars](https://img.shields.io/github/stars/SimonGino/i18n-manager?style=flat-square&logo=github)](https://github.com/SimonGino/i18n-manager/stargazers)

[English](README.md) | [简体中文](README_CN.md)

</div>

---

一个强大的多语言属性文件管理工具，专为Java项目的国际化(i18n)设计。这个工具可以帮助您轻松管理和同步`message-application.properties`文件中的翻译。

## 特性

- 🤖 智能翻译：使用兼容OpenAI的API进行自动文本翻译
- 🔑 智能键生成：自动生成符合Java属性标准的键
- 🔄 自动同步：自动从简体中文(zh)同步到繁体中文(zh_TW)
- 📝 手动管理：支持手动添加和更新翻译
- 🔍 检查工具：验证缺失的翻译条目

## 快速安装

### 自动安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/install.sh | bash
```

### 手动安装

从[Releases](https://github.com/SimonGino/i18n-manager/releases/latest)页面下载适合您平台的二进制文件。

### 支持的平台

- Linux
  - x86_64 (amd64)
  - i386 (32位)
  - arm64
  - armv7
- macOS
  - x86_64 (amd64)
  - arm64 (Apple Silicon)
- Windows
  - x86_64 (64位)
  - i386 (32位)

## 配置

首次使用前，请配置您的API密钥：

1. 从您选择的AI提供商获取API密钥：
   - OpenAI：访问 [OpenAI API Keys](https://platform.openai.com/api-keys)
   - DeepSeek：访问 [DeepSeek Dashboard](https://platform.deepseek.com/api_keys)
   - 通义千问：访问 [DashScope Console](https://dashscope.console.aliyun.com/apiKey)
   - 或任何其他兼容OpenAI的API提供商

2. 配置API密钥、API URL和模型：
```bash
# 设置API密钥
i18n-manager config --set-api-key YOUR_API_KEY

# 设置API URL（可选，默认为OpenAI的端点）
i18n-manager config --set-api-url "https://api.openai.com/v1/chat/completions"

# 设置模型名称（可选，默认为gpt-3.5-turbo）
i18n-manager config --set-model "gpt-3.5-turbo"
```

配置文件位于 `~/.config/i18n-manager/config.json`（Windows 系统位于 `%APPDATA%\i18n-manager\config.json`）。以下是一个配置示例：

```json
{
  "api_key": "your-api-key",
  "api_url": "https://api.openai.com/v1/chat/completions",
  "model": "gpt-3.5-turbo",
  "default_path": ".",
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

### 配置选项

- `api_key`: 您的 AI 提供商 API 密钥
- `api_url`: AI 服务的 API 端点 URL
- `model`: 用于翻译的模型名称
- `default_path`: 属性文件的默认路径
- `language`: 语言配置
  - `file_pattern`: 属性文件的命名模式（如 "message-application%s.properties"）
  - `mappings`: 语言映射
    - `code`: 语言代码（如 "en"、"zh"、"zh_CN"）
    - `file`: 文件后缀（如 ""、"_zh"、"_zh_CN"）
    - `is_source`: 是否为源语言（用于翻译）

你可以使用以下命令修改这些设置：

```bash
i18n-manager config --set-api-key "your-api-key"
i18n-manager config --set-api-url "https://api.openai.com/v1/chat/completions"
i18n-manager config --set-model "gpt-3.5-turbo"
i18n-manager config --show  # 显示当前配置
```

### Azure OpenAI 配置

对于 Azure OpenAI 服务，您需要配置额外的参数：

```bash
# 设置您的 Azure OpenAI API 密钥
i18n-manager config --set-api-key "your-azure-api-key"

# 设置您的 Azure OpenAI 端点
# 格式：https://your-resource-name.openai.azure.com/openai/deployments/your-deployment-name/chat/completions
i18n-manager config --set-api-url "https://your-resource-name.openai.azure.com/openai/deployments/your-deployment-name/chat/completions"

# 设置您的部署模型名称
i18n-manager config --set-model "gpt-35-turbo"

# 设置 Azure API 版本（Azure OpenAI 必需）
i18n-manager config --set-azure-api-version "2024-02-15-preview"
```

Azure OpenAI 配置文件示例：

```json
{
  "api_key": "your-azure-api-key",
  "api_url": "https://your-resource-name.openai.azure.com/openai/deployments/your-deployment-name/chat/completions",
  "model": "gpt-35-turbo",
  "azure_api_version": "2024-02-15-preview",
  "default_path": ".",
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
        "code": "zh_TW",
        "file": "_zh_TW",
        "is_source": false
      }
    ]
  }
}
```

### 支持的模型和API端点

您可以使用任何兼容OpenAI的API端点和模型。以下是一些示例：

- OpenAI
  - URL: `https://api.openai.com/v1/chat/completions`
  - 模型: `gpt-3.5-turbo`, `gpt-4` 等

- Azure OpenAI
  - URL: `https://your-resource-name.openai.azure.com/openai/deployments/your-deployment-name/chat/completions`
  - 模型: 取决于您的部署（如 `gpt-35-turbo`, `gpt-4`）
  - API版本: `2024-02-15-preview`（或其他支持的版本）

- DeepSeek
  - URL: `https://api.deepseek.com/v1/chat/completions`
  - 模型: `deepseek-chat` 等

- 通义千问 (Qwen)
  - URL: `https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions`
  - 模型: `qwen-plus` 等

- 自托管模型 (如 LM Studio, Ollama 等)
  - URL: `http://localhost:1234/v1/chat/completions` (根据需要调整端口)
  - 模型: 取决于您的设置

## 使用方法

### 1. 智能翻译

使用自动生成的键进行翻译（默认命令）：

```bash
# 直接翻译，无需 'translate' 命令
i18n-manager "要翻译的文本"
```

使用自定义键进行翻译：

```bash
# 直接翻译并指定自定义键
i18n-manager --key custom.key.name "要翻译的文本"
i18n-manager -k custom.key.name "要翻译的文本"
```

替代语法（仍然支持）：

```bash
# 原始命令语法
i18n-manager translate "要翻译的文本"
i18n-manager translate --key custom.key.name "要翻译的文本"
```

### 2. 手动添加翻译

添加完整的多语言翻译：

```bash
i18n-manager add \
    --key "custom.key.name" \
    --zh "简体中文" \
    --en "English" \
    --zh-tw "繁體中文"
```

或添加部分翻译：

```bash
i18n-manager add \
    --key "custom.key.name" \
    --zh "简体中文" \
    --en "English"
```

### 3. 查看和检查

列出所有翻译键：

```bash
i18n-manager list
```

查看指定键的翻译：

```bash
i18n-manager list --key "error.skill.unavailable"
# 或使用简写
i18n-manager list -k "error.skill.unavailable"
```

检查缺失的翻译：

```bash
i18n-manager check
```

### 4. 配置管理

设置API密钥：

```bash
i18n-manager config --set-api-key YOUR_API_KEY
```

查看当前配置：

```bash
i18n-manager config --show
```

## 键命名约定

生成的键遵循以下约定：

- 使用小写字母、数字和点(.)
- 使用点(.)作为层级分隔符
- 使用常见前缀进行分类：
  - `error.` - 错误消息
  - `success.` - 成功消息
  - `info.` - 信息提示
  - `label.` - UI标签
  - `button.` - 按钮文本
  - `title.` - 页面/章节标题
  - `msg.` - 一般消息
  - `validation.` - 验证消息

示例：

- `validation.username.notEmpty`
- `success.data.saved`
- `button.submit.text`

## 文件结构

工具管理指定目录中的以下文件：

- `message-application.properties` - 英文翻译
- `message-application_zh.properties` - 简体中文翻译
- `message-application_zh_TW.properties` - 繁体中文翻译

## 卸载

要卸载该工具：

### 自动卸载（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/uninstall.sh -o uninstall.sh
chmod +x uninstall.sh
./uninstall.sh
```

### 手动卸载

1. 删除可执行文件

```bash
# Linux/macOS
rm $(which i18n-manager)

# Windows
del C:\path\to\i18n-manager.exe
```

2. 删除配置文件（可选）

```bash
# Linux/macOS
rm -rf ~/.config/i18n-manager

# Windows
rd /s /q %APPDATA%\i18n-manager
```

## 开发

1. 克隆项目：

```bash
git clone https://github.com/yourusername/i18n-manager.git
```

2. 创建cmd目录结构：

```bash
mkdir -p cmd/i18n-manager
```

3. 安装依赖：

```bash
go mod download
```

4. 运行测试：

```bash
go test ./...
```

5. 构建项目：

```bash
go build -o i18n-manager cmd/i18n-manager/main.go
```

## 贡献

欢迎提交Pull Request和问题报告！

## 许可证

MIT License