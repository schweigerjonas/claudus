package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(fileContent))
}

func parseJsonArguments(jsonArgs string, args *ReadArguments) {
	/**
	* format field file_path to be unmarshalled by json lib
	* this field will always be present, since it is defined as required in the tool advertisement to
	* he model
	* in the case of "file_path" being part of the value of the field, this should not be a problem,
	* since I only replace the first occurence
	 */
	jsonArgs = strings.Replace(jsonArgs, "file_path", "FilePath", 1)

	err := json.Unmarshal([]byte(jsonArgs), &args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
