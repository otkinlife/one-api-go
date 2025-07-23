package ai

import "fmt"

// APIError API 错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("OpenAI API error (code: %d, type: %s): %s", e.Code, e.Type, e.Message)
}

// RateLimitError 速率限制错误
type RateLimitError struct {
	*APIError
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit exceeded, retry after %d seconds: %s", e.RetryAfter, e.APIError.Error())
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ValidateRequest 验证请求
func ValidateRequest(req *ChatRequest) error {
	if req.Model == "" {
		return &ValidationError{Field: "model", Message: "model is required"}
	}

	if len(req.Messages) == 0 {
		return &ValidationError{Field: "messages", Message: "at least one message is required"}
	}

	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return &ValidationError{Field: "temperature", Message: "temperature must be between 0 and 2"}
	}

	if req.TopP != nil && (*req.TopP < 0 || *req.TopP > 1) {
		return &ValidationError{Field: "top_p", Message: "top_p must be between 0 and 1"}
	}

	return nil
}
