package ai

// RequestBuilder 请求构建器
type RequestBuilder struct {
	request *ChatRequest
}

// NewRequest 创建请求构建器
func NewRequest(model string) *RequestBuilder {
	return &RequestBuilder{
		request: &ChatRequest{
			Model: model,
			Extra: make(map[string]interface{}),
		},
	}
}

// Messages 设置消息
func (b *RequestBuilder) Messages(messages []Message) *RequestBuilder {
	b.request.Messages = messages
	return b
}

// Temperature 设置温度
func (b *RequestBuilder) Temperature(temp float64) *RequestBuilder {
	b.request.Temperature = &temp
	return b
}

// TopP 设置TopP
func (b *RequestBuilder) TopP(topP float64) *RequestBuilder {
	b.request.TopP = &topP
	return b
}

// MaxTokens 设置最大令牌数
func (b *RequestBuilder) MaxTokens(tokens int) *RequestBuilder {
	b.request.MaxTokens = &tokens
	return b
}

// Stream 设置流式输出
func (b *RequestBuilder) Stream(stream bool) *RequestBuilder {
	b.request.Stream = &stream
	return b
}

// Stop 设置停止词
func (b *RequestBuilder) Stop(stop ...string) *RequestBuilder {
	b.request.Stop = stop
	return b
}

// Tools 设置工具
func (b *RequestBuilder) Tools(tools ...Tool) *RequestBuilder {
	b.request.Tools = tools
	return b
}

// ToolChoice 设置工具选择
func (b *RequestBuilder) ToolChoice(choice interface{}) *RequestBuilder {
	b.request.ToolChoice = choice
	return b
}

// Extra 设置额外参数
func (b *RequestBuilder) Extra(key string, value interface{}) *RequestBuilder {
	b.request.Extra[key] = value
	return b
}

// Build 构建请求
func (b *RequestBuilder) Build() *ChatRequest {
	return b.request
}
