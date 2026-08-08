package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	var model string
	// when submitting, there is no access to .env files, so the default model provided by
	// CodeCrafters will be used to run the tests
	if os.Getenv("LOCAL") == "true" {
		model = os.Getenv("LOCAL_MODEL")
	} else {
		model = "anthropic/claude-haiku-4.5"
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")

	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	// contains context for the model (all previously sent messages: prompts, responses, tool call
	// content)
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}

	var tools = []openai.ChatCompletionToolUnionParam{
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "read_file",
					Description: openai.String("Read and return the contents of a file"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]any{
							"file_path": map[string]string{
								"type":        "string",
								"description": "The path to the file to read",
							},
						},
						"required": []string{"file_path"},
					},
				},
				Type: "function",
			},
		},
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	// Agent loop
	for {
		res, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:     model,
				MaxTokens: openai.Int(4096),
				Messages:  messages,
				Tools:     tools,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(res.Choices) == 0 {
			panic("No choices in response")
		}

		// You can use print statements as follows for debugging, they'll be visible when running tests.
		// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		// TODO: I feel like this is not the way to do this, I think there should be the option to
		// basically have it as JSON, or rather a normal struct variable with strings, that then gets
		// converted to the expected form of the API call
		msg := openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content:   res.Choices[0].Message.ToAssistantMessageParam().Content,
				ToolCalls: res.Choices[0].Message.ToAssistantMessageParam().ToolCalls,
			},
		}

		messages = append(messages, msg)

		var toolCalls []openai.ChatCompletionMessageToolCallUnionParam

		if toolCalls = msg.OfAssistant.ToolCalls; toolCalls == nil {
			fmt.Println(msg.OfAssistant.Content.OfString.Value)
			break
		}

		for i := range len(toolCalls) {
			content := ExecuteTool(toolCalls[i])

			tool_msg := openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					ToolCallID: *toolCalls[i].GetID(),
					Content: openai.ChatCompletionToolMessageParamContentUnion{
						OfString: openai.String(content),
					},
				},
			}

			messages = append(messages, tool_msg)
		}
	}

	os.Exit(0)
}
