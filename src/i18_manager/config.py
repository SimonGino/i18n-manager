import os
import json

class Config:
    def __init__(self):
        self.config_dir = os.path.expanduser("~/.config/i18n-manager")
        self.config_file = os.path.join(self.config_dir, "config.json")
        self.default_config = {
            "default_path": ".",
            "ai_provider": "deepseek",
            "ai_providers": {
                "deepseek": {
                    "api_key": "",
                    "base_url": "https://api.deepseek.com"
                },
                "qwen": {
                    "api_key": "",
                    "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1"
                }
            }
        }
        self.load_config()

    def set_api_key(self, api_key):
        """设置 API key"""
        provider = self.config.get("ai_provider", "deepseek")
        if provider not in self.config.get("ai_providers", {}):
            self.config["ai_providers"] = self.default_config["ai_providers"]
        
        self.config["ai_providers"][provider]["api_key"] = api_key
        self.save_config()
        
    def get_api_key(self):
        """获取当前提供商的 API key"""
        provider = self.config.get("ai_provider", "deepseek")
        return self.config.get("ai_providers", {}).get(provider, {}).get("api_key", "")

    def set_ai_provider(self, provider: str):
        """设置 AI 提供商"""
        if provider not in ["deepseek", "qwen"]:
            raise ValueError("不支持的 AI 提供商")
            
        # 确保 ai_providers 结构存在
        if "ai_providers" not in self.config:
            self.config["ai_providers"] = self.default_config["ai_providers"]
            
        self.config["ai_provider"] = provider
        self.save_config()

    def get_ai_config(self):
        """获取当前 AI 提供商的配置"""
        provider = self.config.get("ai_provider", "deepseek")
        return self.config.get("ai_providers", {}).get(provider, self.default_config["ai_providers"][provider])

    def load_config(self):
        """加载配置文件"""
        if not os.path.exists(self.config_file):
            self._create_default_config()
        
        with open(self.config_file, 'r') as f:
            self.config = json.load(f)
            
        # 确保配置结构完整
        if "ai_providers" not in self.config:
            self.config["ai_providers"] = self.default_config["ai_providers"]
            self.save_config()

    def _create_default_config(self):
        """创建默认配置文件"""
        os.makedirs(self.config_dir, exist_ok=True)
        self.config = self.default_config.copy()
        self.save_config()

    def save_config(self):
        """保存配置"""
        with open(self.config_file, 'w') as f:
            json.dump(self.config, f, indent=2) 