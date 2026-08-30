package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

// llm is the Z.ai chat client used by the bot.
var llm openai.Client

// llmModel is the model name used for chat completions.
var llmModel string

// initLLM builds the chat client from cfg.
func initLLM(cfg *Config) {
	llm = openai.NewClient(
		option.WithAPIKey(cfg.ZAI.APIKey),
		option.WithBaseURL(cfg.ZAI.BaseURL),
	)
	llmModel = cfg.ZAI.Model
}

// askLLM sends prompt to the configured model and returns its reply text.
func askLLM(ctx context.Context, prompt string) (string, error) {
	completion, err := llm.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: llmModel,
	})
	if err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("model returned no choices")
	}
	return completion.Choices[0].Message.Content, nil
}
