package main

import (
	"fmt"

	"github.com/openai/openai-go/v3"
)

func ExecuteTool(toolCall openai.ChatCompletionMessageToolCallUnion) {
	if toolCall.Type != "function" {
		panic("Custom tools not supported")
	}

	if toolCall.Function.Name == "read_file" {
		ExecuteReadTool(toolCall)
	}
}

func ExecuteReadTool(toolCall openai.ChatCompletionMessageToolCallUnion) {
	fmt.Println(toolCall)
}
