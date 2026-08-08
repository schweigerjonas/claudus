package main

import (
	"fmt"

	"github.com/openai/openai-go/v3"
)

func ExecuteTool(toolCall openai.ChatCompletionMessageToolCallUnion) {
	if toolCall.Type != "function" {
		panic("Custom tools not supported")
	}

	switch toolCall.Function.Name {
	case "read_file":
		ExecuteReadTool(toolCall)
	default:
		fmt.Println("Requested tool not supported")
	}
}

func ExecuteReadTool(toolCall openai.ChatCompletionMessageToolCallUnion) {
	fmt.Println(toolCall)
}
