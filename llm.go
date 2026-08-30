package main

import (
	"context"
	"fmt"
	"time"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"go.uber.org/zap"
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

	logger.Info("llm client ready",
		zap.String("base_url", cfg.ZAI.BaseURL),
		zap.String("model", cfg.ZAI.Model),
	)
}

// askLLM sends prompt to the configured model and returns its reply text.
func askLLM(ctx context.Context, prompt string) (string, error) {
	logger.Debug("llm request",
		zap.String("model", llmModel),
		zap.Int("prompt_len", len(prompt)),
	)

	start := time.Now()
	completion, err := llm.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: llmModel,
	})
	elapsed := time.Since(start)

	if err != nil {
		logger.Error("llm completion failed",
			zap.String("model", llmModel),
			zap.Duration("latency", elapsed),
			zap.Error(err),
		)
		return "", err
	}

	if len(completion.Choices) == 0 {
		logger.Warn("llm returned no choices",
			zap.String("model", llmModel),
			zap.Duration("latency", elapsed),
		)
		return "", fmt.Errorf("model returned no choices")
	}

	reply := completion.Choices[0].Message.Content

	logger.Info("llm response",
		zap.String("model", completion.Model),
		zap.Duration("latency", elapsed),
		zap.Int("reply_len", len(reply)),
		zap.Int64("prompt_tokens", completion.Usage.PromptTokens),
		zap.Int64("completion_tokens", completion.Usage.CompletionTokens),
		zap.Int64("total_tokens", completion.Usage.TotalTokens),
	)

	return reply, nil
}
