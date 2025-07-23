# AI Go Client 测试指南

## 配置文件

创建 `config.json` 文件：

```json
{
  "openai": {
    "base_url": "https://你的newapi服务器地址",
    "api_key": "your-api-key",
    "timeout": "60s",
    "retry_count": 3,
    "headers": {
      "User-Agent": "OpenAI-Go-Client/1.0"
    }
  },
  "models": {
    "chat": "gpt-4.1",
    "vision": "gpt-4.1", 
    "function": "gpt-4.1"
  },
  "defaults": {
    "temperature": 0.7,
    "max_tokens": 1000,
    "top_p": 0.9
  }
}