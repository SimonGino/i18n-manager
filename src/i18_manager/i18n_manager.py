import codecs
import argparse
import os
import sys
from typing import Dict
from openai import OpenAI
from dataclasses import dataclass
import json
from i18_manager.config import Config


@dataclass
class TranslationResponse:
    key: str
    translations: dict
    status: bool = True
    message: str = ""


class TranslationService:
    def __init__(self, api_key: str):
        self.client = OpenAI(
            api_key=api_key,
            base_url="https://api.deepseek.com"
        )

    def translate(self, text: str) -> TranslationResponse:
        try:
            # 系统提示词保持不变
            system_prompt = """You are a translation API that ONLY responds in JSON format.
            ALWAYS follow this exact format for ANY input:
            {
                "key": "<java_properties_key>",
                "translations": {
                    "en": "<english_translation>",
                    "zh": "<simplified_chinese>",
                    "zh_TW": "<traditional_chinese>"
                },
                "status": "success"
            }

            Key Generation Rules:
            1. Use lowercase letters, numbers, and dots (.)
            2. Use dots (.) as separators for hierarchical keys
            3. Use common prefixes for different modules/categories:
               - error. for error messages
               - success. for success messages
               - info. for information messages
               - label. for UI labels
               - button. for button texts
               - title. for page/section titles
               - msg. for general messages
               - validation. for validation messages"""

            # 调用 DeepSeek API
            response = self.client.chat.completions.create(
                model="deepseek-chat",
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": f"Translate this text: {text}"}
                ],
                stream=False
            )

            # 解析响应
            result = json.loads(response.choices[0].message.content)

            return TranslationResponse(
                key=result['key'],
                translations=result['translations'],
                status=True
            )

        except json.JSONDecodeError as e:
            return TranslationResponse(
                key="",
                translations={},
                status=False,
                message=f"JSON解析错误: {str(e)}"
            )
        except Exception as e:
            return TranslationResponse(
                key="",
                translations={},
                status=False,
                message=f"翻译服务错误: {str(e)}"
            )


class I18nManager:
    def __init__(self, base_path: str, translation_service: TranslationService):
        self.base_path = base_path
        self.translation_service = translation_service
        self.files = {
            "en": "message-application.properties",  # 英文（默认文件）
            "zh": "message-application_zh.properties",  # 简体中文
            "zh_CN": "message-application_zh_CN.properties",  # 简体中文（与zh相同）
            "zh_TW": "message-application_zh_TW.properties"  # 繁体中文
        }

    def _ensure_file_exists(self, filename: str) -> None:
        """确保文件存在，如果不存在则创建"""
        full_path = os.path.join(self.base_path, filename)
        if not os.path.exists(full_path):
            with open(full_path, 'w', encoding='utf-8') as f:
                pass

    def _load_existing_translations(self, filename: str) -> Dict[str, str]:
        """加载已存在的翻译"""
        translations = {}
        full_path = os.path.join(self.base_path, filename)
        if os.path.exists(full_path):
            with codecs.open(full_path, 'r', 'utf-8') as f:
                for line in f:
                    line = line.strip()
                    if line and not line.startswith('#'):
                        try:
                            key, value = line.split('=', 1)
                            translations[key.strip()] = value.strip()
                        except ValueError:
                            continue
        return translations

    def add_translation(self, key: str, translations: Dict[str, str]) -> None:
        """添加新的翻译"""
        # 首先检查所有语言文件中是否存在该key
        key_exists = False
        existing_translations = {}
        
        for lang in translations.keys():
            if lang not in self.files:
                print(f"警告: 不支持的语言 {lang}")
                continue
            
            filename = self.files[lang]
            self._ensure_file_exists(filename)
            existing = self._load_existing_translations(filename)
            
            if key in existing:
                key_exists = True
                existing_translations[lang] = existing[key]

        # 如果key存在，显示当前值并询问是否更新
        if key_exists:
            print(f"\n键 '{key}' 已存在于以下语言文件中：")
            for lang, value in existing_translations.items():
                print(f"{lang}: {value}")
                print(f"新值: {translations[lang]}")
            
            response = input("\n是否要更新这些翻译? (y/N): ")
            if response.lower() != 'y':
                print("已取消更新")
                return

        # 继续处理每个语言的翻译
        for lang, text in translations.items():
            if lang not in self.files:
                continue

            filename = self.files[lang]
            full_path = os.path.join(self.base_path, filename)

            # 如果是中文，转换为Unicode
            if lang.startswith("zh"):
                text = text.encode('unicode_escape').decode()

            # 更新或添加翻译
            if key in existing_translations:
                # 更新现有值
                lines = []
                with codecs.open(full_path, 'r', 'utf-8') as f:
                    for line in f:
                        if line.startswith(f"{key}="):
                            lines.append(f"{key}={text}\n")
                        else:
                            lines.append(line)
                
                # 确保最后一行有换行符
                if lines and not lines[-1].endswith('\n'):
                    lines[-1] += '\n'
                
                with codecs.open(full_path, 'w', 'utf-8') as f:
                    f.writelines(lines)
            else:
                # 追加新值前确保文件以换行符结束
                with codecs.open(full_path, 'r+', 'utf-8') as f:
                    content = f.read()
                    f.seek(0, 2)  # 移动到文件末尾
                    if content and not content.endswith('\n'):
                        f.write('\n')  # 如果最后一行没有换行符，添加一个
                    f.write(f"{key}={text}\n")

            print(f"已更新 {filename}")

            # 如果更新的是zh文件，自动同步到zh_CN文件
            if lang == 'zh':
                zh_cn_filename = self.files['zh_CN']
                self._ensure_file_exists(zh_cn_filename)
                zh_cn_full_path = os.path.join(self.base_path, zh_cn_filename)
                
                # 直接同步到zh_CN，无需询问
                if key in existing_translations:
                    # 更新现有值
                    lines = []
                    with codecs.open(zh_cn_full_path, 'r', 'utf-8') as f:
                        for line in f:
                            if line.startswith(f"{key}="):
                                lines.append(f"{key}={text}\n")
                            else:
                                lines.append(line)
                    
                    # 确保最后一行有换行符
                    if lines and not lines[-1].endswith('\n'):
                        lines[-1] += '\n'
                    
                    with codecs.open(zh_cn_full_path, 'w', 'utf-8') as f:
                        f.writelines(lines)
                else:
                    # 追加新值前确保文件以换行符结束
                    with codecs.open(zh_cn_full_path, 'r+', 'utf-8') as f:
                        content = f.read()
                        f.seek(0, 2)  # 移动到文件末尾
                        if content and not content.endswith('\n'):
                            f.write('\n')  # 如果最后一行没有换行符，添加一个
                        f.write(f"{key}={text}\n")
                
                print(f"已自动同步到 {zh_cn_filename}")

    def _sync_translation(self, key: str, text: str, target_lang: str) -> None:
        """同步翻译到目标语言文件"""
        filename = self.files[target_lang]
        self._ensure_file_exists(filename)
        existing = self._load_existing_translations(filename)
        existing[key] = text

        full_path = os.path.join(self.base_path, filename)
        with codecs.open(full_path, 'w', 'utf-8') as f:
            for k, v in sorted(existing.items()):
                f.write(f'{k}={v}\n')
        print(f"已同步更新 {filename}")

    def list_keys(self) -> None:
        """列出所有翻译键"""
        all_keys = set()
        for filename in self.files.values():
            translations = self._load_existing_translations(filename)
            all_keys.update(translations.keys())

        print("\n所有翻译键:")
        for key in sorted(all_keys):
            print(key)

    def check_missing(self) -> None:
        """检查缺失的翻译"""
        all_keys = set()
        file_translations = {}

        # 收集所有键和翻译
        for lang, filename in self.files.items():
            translations = self._load_existing_translations(filename)
            file_translations[lang] = translations
            all_keys.update(translations.keys())

        # 检查每个文件的缺失键
        print("\n缺失的翻译:")
        for lang, translations in file_translations.items():
            # 跳过zh_CN的检查，因为它应该与zh完全相同
            if lang == 'zh_CN':
                continue
            missing = all_keys - set(translations.keys())
            if missing:
                print(f"\n{self.files[lang]} 缺失以下键:")
                for key in sorted(missing):
                    print(f"  - {key}")

    def smart_translate(self, text: str) -> bool:
        """智能翻译并添加到文件"""
        print(f"正在翻译: {text}")

        # 调用翻译服务
        result = self.translation_service.translate(text)

        if not result.status:
            print(f"翻译失败: {result.message}")
            return False

        print(f"生成的键: {result.key}")
        print("翻译结果:")
        for lang, trans in result.translations.items():
            print(f"{lang}: {trans}")

        # 确认是否添加
        response = input("\n是否添加这些翻译? (y/N): ")
        if response.lower() != 'y':
            print("已取消添加翻译")
            return False

        # 添加翻译
        self.add_translation(result.key, result.translations)
        return True


def main():
    parser = argparse.ArgumentParser(description='I18n Manager Tool')
    parser.add_argument('--version', action='version', version='i18n-manager v0.1.9')
    parser.add_argument('--path', default='.', help='properties文件所在目录路径')
    parser.add_argument('--api-key', help='翻译API的密钥（可选，优先级高于配置文件）')

    subparsers = parser.add_subparsers(dest='command', help='可用命令')

    # 智能翻译命令
    translate_parser = subparsers.add_parser('translate', help='智能翻译并添加')
    translate_parser.add_argument('text', help='要翻译的文本')
    translate_parser.add_argument('--key', help='指定翻译的key')


    # 添加新翻译命令
    add_parser = subparsers.add_parser('add', help='添加新的翻译')
    add_parser.add_argument('key', help='翻译键')
    add_parser.add_argument('--en', help='英文翻译')
    add_parser.add_argument('--zh', help='简体中文翻译')
    add_parser.add_argument('--zh_TW', help='繁体中文翻译')

    # list命令
    subparsers.add_parser('list', help='列出所有翻译键')

    # check命令
    subparsers.add_parser('check', help='检查缺失的翻译')

    # 添加配置管理命令
    config_parser = subparsers.add_parser('config', help='管理配置')
    config_parser.add_argument('--set-api-key', help='设置API密钥')
    config_parser.add_argument('--show', action='store_true', help='显示当前配置')

    args = parser.parse_args()

      # 加载配置
    config = Config()
    if args.path:
        config.base_path = args.path
    if args.api_key:  # 使用命令行提供的 API key
        config.api_key = args.api_key
    
    # 如果是配置命令，直接处理
    if args.command == 'config':
        if args.set_api_key:
            config.set_api_key(args.set_api_key)
            print("API密钥已更新")
        if args.show:
            print("\n当前配置:")
            print(f"API密钥: {'*' * 8 + config.get_api_key()[-4:] if config.get_api_key() else '未设置'}")
            print(f"默认路径: {config.config.get('default_path', '.')}")
        return

    # 如果没有提供命令，显示帮助信息并退出
    if not args.command:
        parser.print_help()
        print("\n提示：首次使用请先设置 API 密钥：")
        print("i18n-manager config --set-api-key YOUR_API_KEY")
        sys.exit(1)

    # 检查 API key
    api_key = args.api_key or config.get_api_key()
    if not api_key:
        print("错误：未设置API密钥。请使用以下命令设置：")
        print("i18n-manager config --set-api-key YOUR_API_KEY")
        print("\n或者使用 --api-key 参数临时指定：")
        print("i18n-manager --api-key YOUR_API_KEY <command>")
        sys.exit(1)

    translation_service = TranslationService(api_key)
    manager = I18nManager(args.path, translation_service)

    # 处理其他命令
    if args.command == 'translate':
        print(f"正在翻译: {args.text}")
        result = translation_service.translate(args.text)
        if result.status:
            # 显示生成的键和翻译结果
            original_key = result.key
            if args.key:  # 如果指定了key，使用指定的key
                print(f"使用指定的键: {args.key} (原生成键: {original_key})")
                result.key = args.key
            else:
                print(f"生成的键: {result.key}")
            
            print("\n翻译结果:")
            for lang, trans in result.translations.items():
                print(f"{lang}: {trans}")
            
            # 确认是否添加
            response = input("\n是否添加这些翻译? (y/N): ")
            if response.lower() == 'y':
                manager.add_translation(result.key, result.translations)
                print("翻译已添加")
            else:
                print("已取消添加翻译")
        else:
            print(f"翻译失败: {result.message}")

    elif args.command == 'list':
        manager.list_keys()
    elif args.command == 'check':
        manager.check_missing()
    elif args.command == 'add':
        translations = {}
        if args.en:
            translations['en'] = args.en
        if args.zh:
            translations['zh'] = args.zh
        if args.zh_TW:
            translations['zh_TW'] = args.zh_TW

        if not translations:
            print("错误：至少需要提供一种语言的翻译")
            sys.exit(1)

        manager.add_translation(args.key, translations)


if __name__ == '__main__':
    main()