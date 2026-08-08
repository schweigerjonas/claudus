package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
)

func ExecuteTool(toolCall openai.ChatCompletionMessageToolCallUnionParam) string {
	if toolType := *toolCall.GetType(); toolType != "function" {
		fmt.Fprintf(os.Stderr, "tool type not support: %v\n", toolType)
	}

	switch toolCall.GetFunction().Name {
	case "read_file":
		return ExecuteReadTool(toolCall)
	default:
		log.Fatalf("error: requested tool not supported\n")
	}

	return ""
}

func ExecuteReadTool(toolCall openai.ChatCompletionMessageToolCallUnionParam) string {
	var args ReadArguments // parsed args
	jsonArgs := toolCall.GetFunction().Arguments

	parseJsonArguments(jsonArgs, &args)

	content, err := os.ReadFile(args.FilePath)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	return string(content)
}

func parseJsonArguments(jsonArgs string, args *ReadArguments) {
	err := json.Unmarshal([]byte(jsonArgs), &args)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
