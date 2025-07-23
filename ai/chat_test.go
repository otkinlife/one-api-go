package ai

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// TestClient 测试客户端
var testClient *Client
var testConfig *ConfigFile

func init() {
	// 优先尝试从配置文件加载
	config, configFile, err := LoadConfigFromFile("config.json")
	if err != nil {
		// 如果配置文件不存在，从环境变量加载
		config = LoadConfigFromEnv()
		if config.APIKey == "" {
			fmt.Println("Warning: No API key found in config.json or environment variables")
			// 使用测试用的默认配置
			config = NewConfig("https://api.openai.com", "test-key")
		}
	} else {
		testConfig = configFile
	}

	testClient = NewClient(config)
}

// TestBasicConversation 测试基础对话
func TestBasicConversation(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		System("你是一个测试助手，请简洁回答问题。").
		User("请说'测试成功'").
		Build()

	req := NewRequest(getTestModel("chat")).
		Messages(messages).
		Temperature(0.1).
		MaxTokens(50).
		Build()

	resp, err := testClient.ChatCompletion(ctx, req)
	if err != nil {
		t.Logf("API调用失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}

	if len(resp.Choices) == 0 {
		t.Fatal("没有返回选择项")
	}

	content := ExtractContent(resp.Choices[0].Message)
	t.Logf("基础对话测试通过")
	t.Logf("用户: 请说'测试成功'")
	t.Logf("助手: %s", content)
	t.Logf("使用Token: %d", resp.Usage.TotalTokens)
}

// TestMultiModalMessage 测试多模态消息构建
func TestMultiModalMessage(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		UserWithImages(
			"这是一个测试图片，请描述你看到了什么。",
			"https://ddd.com/imgs/0.jpg",
		).
		Build()

	req := NewRequest(getTestModel("vision")).
		Messages(messages).
		MaxTokens(100).
		Build()

	resp, err := testClient.ChatCompletion(ctx, req)
	if err != nil {
		t.Logf("多模态测试失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}

	t.Logf("多模态测试通过")
	t.Logf("图片描述: %s", ExtractContent(resp.Choices[0].Message))
}

// TestFunctionCalling 测试函数调用
func TestFunctionCalling(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		User("请帮我查询北京的天气").
		Build()

	tools := NewToolBuilder().
		AddWeatherFunction().
		Build()

	req := NewRequest(getTestModel("function")).
		Messages(messages).
		Tools(tools...).
		ToolChoice("auto").
		Build()

	resp, err := testClient.ChatCompletion(ctx, req)
	if err != nil {
		t.Logf("函数调用测试失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}

	t.Logf("函数调用测试通过")
	t.Logf("助手回复: %s", ExtractContent(resp.Choices[0].Message))

	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		t.Logf("函数调用详情:")
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			t.Logf("- 函数: %s", toolCall.Function.Name)
			t.Logf("  参数: %s", toolCall.Function.Arguments)
		}
	}
}

// TestStreamingResponse 测试流式响应
func TestStreamingResponse(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		User("请数数字1到5").
		Build()

	req := NewRequest(getTestModel("chat")).
		Messages(messages).
		Temperature(0.1).
		MaxTokens(100).
		Build()

	stream, err := testClient.ChatCompletionStream(ctx, req)
	if err != nil {
		t.Logf("流式测试失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}
	defer stream.Close()

	t.Logf("流式响应测试开始...")

	var fullContent string
	chunkCount := 0

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Logf("流式接收错误: %v", err)
			break
		}

		chunkCount++
		if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
			content := ExtractContent(resp.Choices[0].Delta)
			if content != "" {
				fullContent += content
			}
		}
	}

	t.Logf("流式响应测试通过")
	t.Logf("接收到 %d 个数据块", chunkCount)
	t.Logf("完整内容: %s", fullContent)
}

// TestCustomParameters 测试自定义参数
func TestCustomParameters(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		User("请简单回答：什么是AI？").
		Build()

	req := NewRequest(getTestModel("chat")).
		Messages(messages).
		Temperature(0.5).
		MaxTokens(100).
		Extra("presence_penalty", 0.1).
		Extra("frequency_penalty", 0.1).
		Extra("top_p", 0.9).
		Build()

	resp, err := testClient.ChatCompletion(ctx, req)
	if err != nil {
		t.Logf("自定义参数测试失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}

	t.Logf("自定义参数测试通过")
	t.Logf("响应: %s", ExtractContent(resp.Choices[0].Message))
	t.Logf("使用的额外参数: presence_penalty=0.1, frequency_penalty=0.1, top_p=0.9")
}

// TestJSONMode 测试JSON模式
func TestJSONMode(t *testing.T) {
	ctx := context.Background()

	messages := NewMessageBuilder().
		System("你是一个数据分析助手，请以JSON格式回复。").
		User("用JSON格式描述一个人的基本信息，包括姓名、年龄、职业").
		Build()

	req := NewRequest(getTestModel("chat")).
		Messages(messages).
		MaxTokens(200).
		Extra("response_format", map[string]string{"type": "json_object"}).
		Build()

	resp, err := testClient.ChatCompletion(ctx, req)
	if err != nil {
		t.Logf("JSON模式测试失败 (这在没有有效API密钥时是正常的): %v", err)
		return
	}

	t.Logf("JSON模式测试通过")
	t.Logf("JSON响应: %s", ExtractContent(resp.Choices[0].Message))
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	// 测试无效请求验证
	req := &ChatRequest{
		Model:       "",                 // 空模型
		Messages:    []Message{},        // 空消息
		Temperature: &[]float64{3.0}[0], // 无效温度
		TopP:        &[]float64{1.5}[0], // 无效TopP
	}

	err := ValidateRequest(req)
	if err == nil {
		t.Fatal("应该检测到验证错误")
	}

	t.Logf("配置验证测试通过: %v", err)

	// 测试有效请求
	validReq := &ChatRequest{
		Model:       "gpt-4.1",
		Messages:    []Message{{Role: "user", Content: "test"}},
		Temperature: &[]float64{0.7}[0],
		TopP:        &[]float64{0.9}[0],
	}

	err = ValidateRequest(validReq)
	if err != nil {
		t.Fatalf("有效请求不应该有验证错误: %v", err)
	}
}

// TestMessageBuilder 测试消息构建器
func TestMessageBuilder(t *testing.T) {
	builder := NewMessageBuilder()

	messages := builder.
		System("你是一个测试助手").
		User("第一个用户消息").
		Assistant("第一个助手回复").
		User("第二个用户消息").
		Build()

	if len(messages) != 4 {
		t.Fatalf("期望4条消息，实际得到%d条", len(messages))
	}

	expectedRoles := []string{"system", "user", "assistant", "user"}
	for i, msg := range messages {
		if msg.Role != expectedRoles[i] {
			t.Errorf("消息%d角色错误，期望%s，实际%s", i, expectedRoles[i], msg.Role)
		}
	}

	t.Logf("消息构建器测试通过")

	// 测试清空功能
	builder.Clear()
	newMessages := builder.User("清空后的消息").Build()
	if len(newMessages) != 1 {
		t.Errorf("清空后期望1条消息，实际%d条", len(newMessages))
	}
}

// TestToolBuilder 测试工具构建器
func TestToolBuilder(t *testing.T) {
	tools := NewToolBuilder().
		AddWeatherFunction().
		AddSearchFunction().
		AddFunction("custom_function", "自定义函数", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]interface{}{
					"type":        "string",
					"description": "参数1",
				},
			},
		}).
		Build()

	if len(tools) != 3 {
		t.Fatalf("期望3个工具，实际得到%d个", len(tools))
	}

	expectedNames := []string{"get_weather", "web_search", "custom_function"}
	for i, tool := range tools {
		if tool.Function.Name != expectedNames[i] {
			t.Errorf("工具%d名称错误，期望%s，实际%s", i, expectedNames[i], tool.Function.Name)
		}
	}

	t.Logf("工具构建器测试通过")
}

// TestDynamicConfig 测试动态配置
func TestDynamicConfig(t *testing.T) {
	originalConfig := testClient.GetConfig()

	// 创建新配置
	newConfig := NewConfig("https://new-api.com", "new-key").
		WithTimeout(60 * time.Second).
		WithHeaders(map[string]string{
			"X-Test": "dynamic-config",
		})

	// 更新配置
	testClient.SetConfig(newConfig)

	// 验证配置更新
	currentConfig := testClient.GetConfig()
	if currentConfig.BaseURL != "https://new-api.com" {
		t.Errorf("配置更新失败，期望BaseURL为https://new-api.com，实际为%s", currentConfig.BaseURL)
	}

	if currentConfig.Headers["X-Test"] != "dynamic-config" {
		t.Errorf("自定义头部设置失败")
	}

	t.Logf("动态配置测试通过")

	// 恢复原始配置
	testClient.SetConfig(originalConfig)
}

// TestExtractContent 测试内容提取
func TestExtractContent(t *testing.T) {
	// 测试字符串内容
	msg1 := &Message{Content: "简单文本"}
	content1 := ExtractContent(msg1)
	if content1 != "简单文本" {
		t.Errorf("字符串内容提取失败，期望'简单文本'，实际'%s'", content1)
	}

	// 测试多模态内容
	msg2 := &Message{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": "文本部分1",
			},
			map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": "https://example.com/image.jpg",
				},
			},
			map[string]interface{}{
				"type": "text",
				"text": "文本部分2",
			},
		},
	}
	content2 := ExtractContent(msg2)
	expected := "[文本部分1 文本部分2]"
	if content2 != expected {
		t.Errorf("多模态内容提取失败，期望'%s'，实际'%s'", expected, content2)
	}

	t.Logf("内容提取测试通过")
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试验证错误
	validationErr := &ValidationError{
		Field:   "temperature",
		Message: "temperature must be between 0 and 2",
	}

	expectedMsg := "validation error for field 'temperature': temperature must be between 0 and 2"
	if validationErr.Error() != expectedMsg {
		t.Errorf("验证错误消息不正确，期望'%s'，实际'%s'", expectedMsg, validationErr.Error())
	}

	// 测试API错误
	apiErr := &APIError{
		Code:    400,
		Message: "Bad Request",
		Type:    "invalid_request_error",
	}

	expectedAPIMsg := "OpenAI API error (code: 400, type: invalid_request_error): Bad Request"
	if apiErr.Error() != expectedAPIMsg {
		t.Errorf("API错误消息不正确，期望'%s'，实际'%s'", expectedAPIMsg, apiErr.Error())
	}

	t.Logf("错误处理测试通过")
}

// BenchmarkMessageBuilder 性能测试 - 消息构建器
func BenchmarkMessageBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewMessageBuilder()
		builder.System("系统消息").
			User("用户消息").
			Assistant("助手回复").
			Build()
	}
}

// BenchmarkRequestBuilder 性能测试 - 请求构建器
func BenchmarkRequestBuilder(b *testing.B) {
	messages := NewMessageBuilder().User("测试消息").Build()

	for i := 0; i < b.N; i++ {
		NewRequest("gpt-4.1").
			Messages(messages).
			Temperature(0.7).
			MaxTokens(100).
			Build()
	}
}

// getTestModel 获取测试模型名称
func getTestModel(modelType string) string {
	if testConfig != nil {
		switch modelType {
		case "chat":
			return testConfig.Models.Chat
		case "vision":
			return testConfig.Models.Vision
		case "function":
			return testConfig.Models.Function
		}
	}
	return "gpt-4.1" // 默认模型
}

// TestMain 测试主函数
func TestMain(m *testing.M) {
	fmt.Println("开始运行OpenAI客户端测试...")
	fmt.Printf("使用配置: %s\n", testClient.GetConfig().BaseURL)

	// 运行测试
	code := m.Run()

	fmt.Println("测试完成")
	os.Exit(code)
}
