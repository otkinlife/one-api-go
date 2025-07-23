package ai

import (
	"encoding/json"
	"os"
	"time"
)

// ConfigFile 配置文件结构
type ConfigFile struct {
	OpenAI struct {
		BaseURL    string            `json:"base_url"`
		APIKey     string            `json:"api_key"`
		Timeout    string            `json:"timeout"`
		RetryCount int               `json:"retry_count"`
		Headers    map[string]string `json:"headers"`
	} `json:"openai"`

	Models struct {
		Chat     string `json:"chat"`
		Vision   string `json:"vision"`
		Function string `json:"function"`
	} `json:"models"`

	Defaults struct {
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"max_tokens"`
		TopP        float64 `json:"top_p"`
	} `json:"defaults"`
}

// LoadConfigFromFile 从JSON文件加载配置
func LoadConfigFromFile(filename string) (*Config, *ConfigFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	var configFile ConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return nil, nil, err
	}

	// 解析超时时间
	timeout, err := time.ParseDuration(configFile.OpenAI.Timeout)
	if err != nil {
		timeout = 30 * time.Second // 默认超时时间
	}

	config := &Config{
		BaseURL:    configFile.OpenAI.BaseURL,
		APIKey:     configFile.OpenAI.APIKey,
		Timeout:    timeout,
		RetryCount: configFile.OpenAI.RetryCount,
		Headers:    configFile.OpenAI.Headers,
	}

	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}

	return config, &configFile, nil
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *Config {
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	config := &Config{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Timeout:    30 * time.Second,
		RetryCount: 3,
		Headers:    make(map[string]string),
	}

	return config
}
