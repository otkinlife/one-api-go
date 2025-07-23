package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client OpenAI API 客户端
type Client struct {
	config *Config
	http   *http.Client
}

// NewClient 创建新的客户端
func NewClient(config *Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		config: config,
		http:   &http.Client{Timeout: config.Timeout},
	}
}

// ChatCompletion 创建聊天补全
func (c *Client) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Stream != nil && *req.Stream {
		return nil, fmt.Errorf("use ChatCompletionStream for streaming requests")
	}

	return c.doChatRequest(ctx, req)
}

// ChatCompletionStream 创建流式聊天补全
func (c *Client) ChatCompletionStream(ctx context.Context, req *ChatRequest) (*StreamReader, error) {
	req.Stream = &[]bool{true}[0]
	return c.doChatStreamRequest(ctx, req)
}

// SetConfig 更新配置
func (c *Client) SetConfig(config *Config) {
	c.config = config
	if config.Timeout > 0 {
		c.http.Timeout = config.Timeout
	}
}

// GetConfig 获取当前配置
func (c *Client) GetConfig() *Config {
	return c.config
}
