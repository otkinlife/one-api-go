package ai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// StreamReader 流式响应读取器
type StreamReader struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
	ctx     context.Context
}

// StreamResponse 流式响应
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice 流式选择
type StreamChoice struct {
	Index        int      `json:"index"`
	Delta        *Message `json:"delta"`
	FinishReason string   `json:"finish_reason,omitempty"`
}

// Recv 接收下一个流式响应
func (s *StreamReader) Recv() (*StreamResponse, error) {
	for s.scanner.Scan() {
		line := s.scanner.Text()

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				return nil, io.EOF
			}

			var response StreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				return nil, fmt.Errorf("failed to unmarshal stream response: %w", err)
			}

			return &response, nil
		}
	}

	if err := s.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

// Close 关闭流
func (s *StreamReader) Close() error {
	return s.reader.Close()
}
