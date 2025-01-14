from openai import OpenAI
from typing import Dict, Optional

class AIService:
    def __init__(self, api_key: str, base_url: str, model: str = "deepseek-chat"):
        self.client = OpenAI(
            api_key=api_key,
            base_url=base_url
        )
        self.model = model

    def chat(self, messages: list, **kwargs) -> Dict:
        try:
            response = self.client.chat.completions.create(
                model=self.model,
                messages=messages,
                **kwargs
            )
            return {
                "success": True,
                "content": response.choices[0].message.content
            }
        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

class AIServiceFactory:
    @staticmethod
    def create(provider: str, config: Dict) -> Optional[AIService]:
        if provider == "deepseek":
            return AIService(
                api_key=config["api_key"],
                base_url=config["base_url"],
                model="deepseek-chat"
            )
        elif provider == "qwen":
            return AIService(
                api_key=config["api_key"],
                base_url=config["base_url"],
                model="qwen-plus"
            )
        return None