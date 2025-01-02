# i18_manager
以下是一个适合您的国际化管理工具的 README.md 建议：

```markdown
# I18n Manager Tool

一个强大的多语言属性文件管理工具，专为 Java 项目的国际化(i18n)设计。该工具可以帮助您轻松管理和同步 `message-application.properties` 文件中的多语言翻译。

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

## 安装

```bash
# Clone 仓库
git clone [repository-url]
cd i18n-manager

# 安装依赖
pip install -r requirements.txt
```

## 使用方法

### 基础参数

工具支持以下全局参数：

- `--path`：properties 文件所在目录路径（默认为 ./lang）
- `--api-key`：翻译 API 的密钥（默认为 'app-xxxxxxx'）

### 可用命令

1. 智能翻译
```bash
python i18n_manager.py translate "需要翻译的文本"
```
该命令会：
- 自动生成合适的翻译键
- 翻译到所有支持的语言
- 在添加前请求确认
- 自动同步 zh_CN 与 zh 的内容

2. 手动添加翻译
```bash
python i18n_manager.py add <translation.key> [--en "English text"] [--zh "中文文本"] [--zh_TW "繁體中文文本"]
```
注意：
- 至少需要提供一种语言的翻译
- 如果翻译键已存在，会请求确认是否覆盖
- zh_CN 会自动与 zh 同步，无需手动添加

3. 检查缺失的翻译
```bash
python i18n_manager.py check
```
该命令会：
- 检查所有语言文件中缺失的翻译键
- 显示每个文件中缺失的具体键名
- 自动跳过 zh_CN 的检查（因为它与 zh 同步）

4. 列出所有翻译键
```bash
python i18n_manager.py list
```
该命令会：
- 显示所有语言文件中存在的翻译键
- 按字母顺序排序显示

### 示例

1. 使用自定义路径和 API 密钥：
```bash
python i18n_manager.py --path ./i18n --api-key "your-api-key" translate "Hello World"
```

2. 添加多语言翻译：
```bash
python i18n_manager.py add welcome.message --en "Welcome" --zh "欢迎" --zh_TW "歡迎"
```

## 配置

工具支持以下配置选项：

- `--path`：properties 文件所在目录路径（默认为当前目录）
- `--api-key`：翻译 API 的密钥

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

## 贡献

欢迎提交 Issue 和 Pull Request！
