# I18n Manager Tool

![Build Status](https://github.com/SimonGino/i18n_manager/actions/workflows/build.yml/badge.svg)

一个强大的多语言属性文件管理工具，专为 Java 项目的国际化(i18n)设计。该工具可以帮助您轻松管理和同步 `message-application.properties` 文件中的多语言翻译。

## 快速安装

### 自动安装（推荐）
```bash
curl -fsSL https://raw.githubusercontent.com/SimonGino/i18n_manager/main/install.sh | bash
```

### 手动安装
从 [Releases](https://github.com/SimonGino/i18n_manager/releases/latest) 页面下载对应平台的二进制文件。

支持的平台：
- Linux (amd64)
- macOS (amd64, arm64)

## 特性

- 🚀 智能翻译：自动生成合适的翻译键并翻译到多种语言
- 🔄 自动同步：自动同步简体中文(zh)到中国大陆(zh_CN)的翻译
- 📝 手动添加：支持手动添加新的翻译条目
- 🔍 翻译检查：检查所有语言文件中缺失的翻译
- 📋 键值列表：查看所有已存在的翻译键
- 🔒 安全确认：更新已存在的翻译时需要确认

## 支持的语言

- English (en)
- 简体中文 (zh/zh_CN)
- 繁體中文 (zh_TW)

## 使用方法

### 基础参数

工具支持以下全局参数：

- `--path`：properties 文件所在目录路径（默认为 ./lang）
- `--api-key`：翻译 API 的密钥（默认为 'app-xxxxxxx'）

### 可用命令

1. 智能翻译
```bash
i18n_manager translate "需要翻译的文本"
```

2. 手动添加翻译
```bash
i18n_manager add <translation.key> [--en "English text"] [--zh "中文文本"] [--zh_TW "繁體中文文本"]
```

3. 检查缺失的翻译
```bash
i18n_manager check
```

4. 列出所有翻译键
```bash
i18n_manager list
```

### 示例

1. 使用自定义路径和 API 密钥：
```bash
i18n_manager --path ./i18n --api-key "your-api-key" translate "Hello World"
```

2. 添加多语言翻译：
```bash
i18n_manager add welcome.message --en "Welcome" --zh "欢迎" --zh_TW "歡迎"
```

## 文件结构

工具管理以下属性文件：

```
message-application.properties        # 英文（默认）
message-application_zh.properties     # 简体中文
message-application_zh_CN.properties  # 简体中文（自动同步自 zh）
message-application_zh_TW.properties  # 繁体中文
```

## 最佳实践

1. 在添加新翻译之前，建议先使用 `check` 命令检查现有翻译是否完整
2. 使用 `translate` 命令可以确保翻译键的一致性和翻译质量
3. 定期检查和更新翻译以保持所有语言版本的同步

## 注意事项

- 更新已存在的翻译时会提示确认（除了 zh_CN，它会自动同步自 zh）
- 中文翻译会自动转换为 Unicode 编码以确保兼容性
- 建议定期备份您的属性文件

## 卸载

如果需要卸载工具：
```bash
rm -f ~/.local/bin/i18n_manager
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License