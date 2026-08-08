package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	// fmt.Printf("tool id:%v; type:%v; name:%v; arguments:%v", toolCall.ID, toolCall.Type, toolCall.Function.Name, toolCall.Function.Arguments)
	var args ReadArguments // parsed args
	jsonArgs := toolCall.Function.Arguments

	parseJsonArguments(jsonArgs, &args)

	fileContent, err := os.ReadFile(args.FilePath)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Println(string(fileContent))
}

func parseJsonArguments(jsonArgs string, args *ReadArguments) {
	err := json.Unmarshal([]byte(jsonArgs), &args)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
