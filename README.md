# I18n Manager Tool

![Build Status](https://github.com/SimonGino/i18n-manager/actions/workflows/build.yml/badge.svg)

一个强大的多语言属性文件管理工具，专为 Java 项目的国际化(i18n)设计。该工具可以帮助您轻松管理和同步 `message-application.properties` 文件中的多语言翻译。

## 快速安装

### 自动安装（推荐）
```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/install.sh | bash
```

### 手动安装
从 [Releases](https://github.com/SimonGino/i18n-manager/releases/latest) 页面下载对应平台的二进制文件。

支持的平台：
- Linux (amd64)
- macOS (amd64, arm64)

## 初始配置

首次使用前，需要设置 API 密钥：

```bash
i18n-manager config --set-api-key YOUR_API_KEY
```

查看当前配置：
```bash
i18n-manager config --show
```

## 特性

- 🚀 智能翻译：自动生成合适的翻译键并翻译到多种语言
- 🔄 自动同步：自动同步简体中文(zh)到中国大陆(zh_CN)的翻译
- 📝 手动添加：支持手动添加新的翻译条目
- 🔍 翻译检查：检查所有语言文件中缺失的翻译
- 📋 键值列表：查看所有已存在的翻译键
- 🔒 安全确认：更新已存在的翻译时需要确认
- ⚙️ 配置管理：支持通过配置文件管理 API 密钥等设置

## 支持的语言

- English (en)
- 简体中文 (zh/zh_CN)
- 繁體中文 (zh_TW)

## 使用方法

### 基础参数

工具支持以下全局参数：

- `--path`：properties 文件所在目录路径（默认为当前目录）
- `--api-key`：临时使用的翻译 API 密钥（可选，优先级高于配置文件）

### 可用命令

1. 配置管理
```bash
i18n-manager config --set-api-key YOUR_API_KEY  # 设置 API 密钥
i18n-manager config --show                      # 显示当前配置
```

2. 智能翻译
```bash
i18n-manager translate "需要翻译的文本"
```

3. 手动添加翻译
```bash
i18n-manager add <translation.key> [--en "English text"] [--zh "中文文本"] [--zh_TW "繁體中文文本"]
```

4. 检查缺失的翻译
```bash
i18n-manager check
```

5. 列出所有翻译键
```bash
i18n-manager list
```

### 示例

1. 使用自定义路径：
```bash
i18n-manager --path ./i18n translate "Hello World"
```

2. 添加多语言翻译：
```bash
i18n-manager add welcome.message --en "Welcome" --zh "欢迎" --zh_TW "歡迎"
```

## 文件结构

工具管理以下属性文件：

```
message-application.properties        # 英文（默认）
message-application_zh.properties     # 简体中文
message-application_zh_CN.properties  # 简体中文（自动同步自 zh）
message-application_zh_TW.properties  # 繁体中文
```

## 配置文件

配置文件位于：
- Linux/macOS: `~/.config/i18n-manager/config.json`

包含以下配置项：
- `api_key`: API 密钥
- `default_path`: 默认工作目录
- `default_source_lang`: 默认源语言
- `default_target_langs`: 默认目标语言列表

## 最佳实践

1. 在添加新翻译之前，建议先使用 `check` 命令检查现有翻译是否完整
2. 使用 `translate` 命令可以确保翻译键的一致性和翻译质量
3. 定期检查和更新翻译以保持所有语言版本的同步
4. 妥善保管 API 密钥，不要将其提交到版本控制系统

## 卸载

如果需要卸载工具：
```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n-manager/main/uninstall.sh | bash
```

卸载脚本会询问是否同时删除配置文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License