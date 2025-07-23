package ai

// MessageBuilder 消息构建器
type MessageBuilder struct {
	messages []Message
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		messages: make([]Message, 0),
	}
}

// System 添加系统消息
func (b *MessageBuilder) System(content string) *MessageBuilder {
	b.messages = append(b.messages, Message{
		Role:    "system",
		Content: content,
	})
	return b
}

// User 添加用户消息
func (b *MessageBuilder) User(content string) *MessageBuilder {
	b.messages = append(b.messages, Message{
		Role:    "user",
		Content: content,
	})
	return b
}

// UserWithImages 添加用户消息（带图片）
func (b *MessageBuilder) UserWithImages(text string, imageUrls ...string) *MessageBuilder {
	var content []map[string]interface{}

	// 添加文本内容
	if text != "" {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": text,
		})
	}

	// 添加图片内容
	for _, url := range imageUrls {
		content = append(content, map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]interface{}{
				"url": url,
			},
		})
	}

	b.messages = append(b.messages, Message{
		Role:    "user",
		Content: content,
	})
	return b
}

// Assistant 添加助手消息
func (b *MessageBuilder) Assistant(content string) *MessageBuilder {
	b.messages = append(b.messages, Message{
		Role:    "assistant",
		Content: content,
	})
	return b
}

// AssistantWithTools 添加带工具调用的助手消息
func (b *MessageBuilder) AssistantWithTools(content string, toolCalls ...ToolCall) *MessageBuilder {
	b.messages = append(b.messages, Message{
		Role:      "assistant",
		Content:   content,
		ToolCalls: toolCalls,
	})
	return b
}

// Tool 添加工具消息
func (b *MessageBuilder) Tool(toolCallId, content string) *MessageBuilder {
	b.messages = append(b.messages, Message{
		Role:    "tool",
		Content: content,
		Name:    toolCallId,
	})
	return b
}

// Build 构建消息数组
func (b *MessageBuilder) Build() []Message {
	return b.messages
}

// Clear 清空消息
func (b *MessageBuilder) Clear() *MessageBuilder {
	b.messages = b.messages[:0]
	return b
}
