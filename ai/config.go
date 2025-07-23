package ai

import (
	"net/http"
	"time"
)

// Config 客户端配置
type Config struct {
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key"`
	Timeout    time.Duration     `json:"timeout"`
	Headers    map[string]string `json:"headers"`
	Proxy      string            `json:"proxy,omitempty"`
	RetryCount int               `json:"retry_count"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:    "https://api.openai.com",
		Timeout:    30 * time.Second,
		RetryCount: 3,
		Headers:    make(map[string]string),
	}
}

// NewConfig 创建配置
func NewConfig(baseURL, apiKey string) *Config {
	config := DefaultConfig()
	config.BaseURL = baseURL
	config.APIKey = apiKey
	return config
}

// WithTimeout 设置超时时间
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithHeaders 设置自定义头部
func (c *Config) WithHeaders(headers map[string]string) *Config {
	for k, v := range headers {
		c.Headers[k] = v
	}
	return c
}

// WithProxy 设置代理
func (c *Config) WithProxy(proxy string) *Config {
	c.Proxy = proxy
	return c
}

// WithRetry 设置重试次数
func (c *Config) WithRetry(count int) *Config {
	c.RetryCount = count
	return c
}

// ToHTTPClient 转换为 HTTP 客户端配置
func (c *Config) ToHTTPClient() *http.Client {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	// 如果有代理设置，可以在这里配置
	// if c.Proxy != "" { ... }

	return client
}
