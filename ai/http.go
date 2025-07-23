package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// doChatRequest 执行聊天请求
func (c *Client) doChatRequest(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 合并额外参数
	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	json.Unmarshal(reqBytes, &reqMap)

	for k, v := range req.Extra {
		reqMap[k] = v
	}

	jsonData, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// doChatStreamRequest 执行流式聊天请求
func (c *Client) doChatStreamRequest(ctx context.Context, req *ChatRequest) (*StreamReader, error) {
	// 合并额外参数
	reqMap := make(map[string]interface{})
	reqBytes, _ := json.Marshal(req)
	json.Unmarshal(reqBytes, &reqMap)

	for k, v := range req.Extra {
		reqMap[k] = v
	}

	jsonData, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return &StreamReader{
		reader:  resp.Body,
		scanner: bufio.NewScanner(resp.Body),
		ctx:     ctx,
	}, nil
}

// setHeaders 设置请求头
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	// 添加自定义头部
	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}
}
