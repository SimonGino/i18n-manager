# 🌍 i18n-manager

<div align="center">

[![Build Status](https://github.com/SimonGino/i18n-manager/actions/workflows/build.yml/badge.svg)](https://github.com/SimonGino/i18n-manager/actions)
[![Release](https://img.shields.io/github/v/release/SimonGino/i18n-manager?style=flat-square&logo=github&color=blue)](https://github.com/SimonGino/i18n-manager/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square&logo=bookstack)](LICENSE)
[![Stars](https://img.shields.io/github/stars/SimonGino/i18n-manager?style=flat-square&logo=github)](https://github.com/SimonGino/i18n-manager/stargazers)

[English](README.md) | [简体中文](README_CN.md)

</div>

---

## 特性

- 🤖 智能翻译：利用 DeepSeek AI 自动翻译文本
- 🔑 智能生成 key：自动生成符合 Java properties 规范的 key
- 🔄 自动同步：自动同步中文简体(zh)到繁体(zh_TW)
- 📝 手动管理：支持手动添加和更新翻译
- 🔍 检查工具：检查缺失的翻译条目

## 快速安装

### 自动安装（推荐）
```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/install.sh | bash
```

### 手动安装
从 [Releases](https://github.com/SimonGino/i18n-manager/releases/latest) 页面下载对应平台的二进制文件。

### 支持的平台
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64) - 即将支持

## 配置

首次使用前，需要配置 DeepSeek API key：

1. 前往 [DeepSeek Dashboard](https://platform.deepseek.com/api_keys) 获取您的 API 密钥。
2. 配置 API key：
```bash
i18n-manager config --set-api-key YOUR_API_KEY
```

## 使用方法

### 1. 智能翻译

自动翻译并生成 key：
```bash
i18n-manager translate "需要翻译的文本"
```

使用自定义 key 翻译：
```bash
i18n-manager translate "需要翻译的文本" --key "custom.key.name"
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

也可以只添加部分语言：
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

设置 API key：
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

包含以下配置项：
- `api_key`: API 密钥
- `default_path`: 默认工作目录
- `default_source_lang`: 默认源语言
- `default_target_langs`: 默认目标语言列表

## Key 命名规范

生成的 key 遵循以下规范：
- 使用小写字母、数字和点号(.)
- 使用点号(.)作为层级分隔符
- 使用常用前缀分类：
  - `error.` - 错误消息
  - `success.` - 成功消息
  - `info.` - 信息提示
  - `label.` - UI 标签
  - `button.` - 按钮文本
  - `title.` - 页面/区域标题
  - `msg.` - 一般消息
  - `validation.` - 验证消息

示例：
- `validation.username.notEmpty`
- `success.data.saved`
- `button.submit.text`

## 文件结构

工具会在指定目录下管理以下文件：
- `message-application.properties` - 英文翻译
- `message-application_zh.properties` - 简体中文翻译
- `message-application_zh_TW.properties` - 繁体中文翻译

## 卸载

如果需要卸载工具：

### 自动卸载（推荐）
```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/uninstall.sh | bash
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
pdm install
```

3. 运行测试：
```bash
pdm run test
```

## 贡献

欢迎提交 Pull Request 或创建 Issue！

## 许可

MIT License