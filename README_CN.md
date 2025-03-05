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

- 🤖 智能翻译：使用DeepSeek AI或通义千问AI进行自动文本翻译
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

首次使用前，需要配置您的API密钥：

1. 获取API密钥：
   - DeepSeek：访问 [DeepSeek Dashboard](https://platform.deepseek.com/api_keys)
   - 通义千问：访问 [DashScope控制台](https://dashscope.console.aliyun.com/apiKey)

2. 配置API密钥和提供商：
```bash
# 设置AI提供商 (deepseek 或 qwen)
i18n-manager config --set-ai-provider qwen

# 设置API密钥
i18n-manager config --set-api-key YOUR_API_KEY
```

## 使用方法

### 1. 智能翻译

使用自动生成的键进行翻译：

```bash
i18n-manager translate "要翻译的文本"
```

使用自定义键进行翻译：

```bash
i18n-manager translate "要翻译的文本" --key "custom.key.name"
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

## 配置文件

配置文件位于：

- Linux/macOS: `~/.config/i18n-manager/config.json`
- Windows: `%APPDATA%\i18n-manager\config.json`

包含以下设置：

- `api_key`：API密钥
- `default_path`：默认工作目录
- `default_source_lang`：默认源语言
- `default_target_langs`：默认目标语言列表

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

2. 安装依赖：

```bash
go mod download
```

3. 运行测试：

```bash
go test ./...
```

## 贡献

欢迎提交Pull Request和问题报告！

## 许可证

MIT License