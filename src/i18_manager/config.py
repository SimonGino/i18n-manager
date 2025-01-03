import os
import json

class Config:
    def __init__(self):
        self.config_dir = os.path.expanduser("~/.config/i18n-manager")
        self.config_file = os.path.join(self.config_dir, "config.json")
        self.load_config()

    def load_config(self):
        """加载配置文件"""
        if not os.path.exists(self.config_file):
            self._create_default_config()
        
        with open(self.config_file, 'r') as f:
            self.config = json.load(f)

    def _create_default_config(self):
        """创建默认配置文件"""
        os.makedirs(self.config_dir, exist_ok=True)
        default_config = {
            "api_key": "",
            "default_path": ".",
            "default_source_lang": "zh",
            "default_target_langs": ["en", "zh_TW"]
        }
        
        with open(self.config_file, 'w') as f:
            json.dump(default_config, f, indent=2)
        
        self.config = default_config

    def save_config(self):
        """保存配置"""
        with open(self.config_file, 'w') as f:
            json.dump(self.config, f, indent=2)

    def get_api_key(self):
        """获取 API key"""
        return self.config.get("api_key", "")

    def set_api_key(self, api_key):
        """设置 API key"""
        self.config["api_key"] = api_key
        self.save_config() 