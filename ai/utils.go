package ai

import (
	"encoding/json"
	"fmt"
)

// ToolBuilder 工具构建器
type ToolBuilder struct {
	tools []Tool
}

// NewToolBuilder 创建工具构建器
func NewToolBuilder() *ToolBuilder {
	return &ToolBuilder{
		tools: make([]Tool, 0),
	}
}

// AddFunction 添加函数工具
func (tb *ToolBuilder) AddFunction(name, description string, parameters interface{}) *ToolBuilder {
	tool := Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}
	tb.tools = append(tb.tools, tool)
	return tb
}

// AddWeatherFunction 添加天气查询函数（预设）
func (tb *ToolBuilder) AddWeatherFunction() *ToolBuilder {
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，如：北京、上海",
			},
			"unit": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"celsius", "fahrenheit"},
				"description": "温度单位",
			},
		},
		"required": []string{"location"},
	}

	return tb.AddFunction("get_weather", "获取指定城市的天气信息", parameters)
}

// AddSearchFunction 添加搜索函数（预设）
func (tb *ToolBuilder) AddSearchFunction() *ToolBuilder {
	parameters := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "搜索关键词",
			},
			"max_results": map[string]interface{}{
				"type":        "integer",
				"description": "最大结果数量",
				"minimum":     1,
				"maximum":     10,
			},
		},
		"required": []string{"query"},
	}

	return tb.AddFunction("web_search", "在网络上搜索信息", parameters)
}

// Build 构建工具数组
func (tb *ToolBuilder) Build() []Tool {
	return tb.tools
}

// JSONSchema JSON Schema 构建器
type JSONSchemaBuilder struct {
	schema map[string]interface{}
}

// NewJSONSchema 创建 JSON Schema 构建器
func NewJSONSchema() *JSONSchemaBuilder {
	return &JSONSchemaBuilder{
		schema: map[string]interface{}{
			"type":       "object",
			"properties": make(map[string]interface{}),
		},
	}
}

// AddProperty 添加属性
func (jsb *JSONSchemaBuilder) AddProperty(name, propType, description string) *JSONSchemaBuilder {
	properties := jsb.schema["properties"].(map[string]interface{})
	properties[name] = map[string]interface{}{
		"type":        propType,
		"description": description,
	}
	return jsb
}

// AddStringProperty 添加字符串属性（带枚举）
func (jsb *JSONSchemaBuilder) AddStringProperty(name, description string, enum ...string) *JSONSchemaBuilder {
	properties := jsb.schema["properties"].(map[string]interface{})
	prop := map[string]interface{}{
		"type":        "string",
		"description": description,
	}
	if len(enum) > 0 {
		prop["enum"] = enum
	}
	properties[name] = prop
	return jsb
}

// AddNumberProperty 添加数字属性
func (jsb *JSONSchemaBuilder) AddNumberProperty(name, description string, min, max *float64) *JSONSchemaBuilder {
	properties := jsb.schema["properties"].(map[string]interface{})
	prop := map[string]interface{}{
		"type":        "number",
		"description": description,
	}
	if min != nil {
		prop["minimum"] = *min
	}
	if max != nil {
		prop["maximum"] = *max
	}
	properties[name] = prop
	return jsb
}

// Required 设置必需字段
func (jsb *JSONSchemaBuilder) Required(fields ...string) *JSONSchemaBuilder {
	jsb.schema["required"] = fields
	return jsb
}

// Build 构建 Schema
func (jsb *JSONSchemaBuilder) Build() map[string]interface{} {
	return jsb.schema
}

// PrettyPrint 格式化打印响应
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

// ExtractContent 提取消息内容
func ExtractContent(message *Message) string {
	if message == nil {
		return ""
	}

	switch content := message.Content.(type) {
	case string:
		return content
	case []interface{}:
		// 处理多模态内容，提取文本部分
		var textParts []string
		for _, part := range content {
			if partMap, ok := part.(map[string]interface{}); ok {
				if partMap["type"] == "text" {
					if text, ok := partMap["text"].(string); ok {
						textParts = append(textParts, text)
					}
				}
			}
		}
		return fmt.Sprintf("[%s]", join(textParts, " "))
	default:
		return fmt.Sprintf("%v", content)
	}
}

// join 字符串连接辅助函数
func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
