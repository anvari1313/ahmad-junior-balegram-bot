package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf8"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

// logger is the global application logger.
var logger *zap.Logger

func main() {
	var err error
	logger, err = newLogger()
	if err != nil {
		panic(err) // cannot log a logger startup failure; nothing to log to yet
	}
	defer logger.Sync()

	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal("load config", zap.Error(err))
	}
	if err := cfg.validate(); err != nil {
		logger.Fatal("validate config", zap.Error(err))
	}

	fmt.Printf("%+v\n", cfg)

	initLLM(cfg)

	// Cancel the context on SIGINT/SIGTERM so polling stops gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	b, err := bot.New(cfg.Telegram.Token,
		bot.WithServerURL(cfg.Telegram.APIURL),
		bot.WithDefaultHandler(defaultHandler),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("telegram update error", zap.Error(err))
		}),
	)
	if err != nil {
		logger.Fatal("init bot", zap.Error(err))
	}

	logger.Info("bot started; polling for updates",
		zap.String("api_url", cfg.Telegram.APIURL),
		zap.String("model", cfg.ZAI.Model),
	)
	b.Start(ctx) // long polling; blocks until ctx is cancelled
	logger.Info("bot stopped")
}

// defaultHandler forwards every incoming message to the LLM and returns the
// model's reply to the user.
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID

	logger.Info("message received",
		zap.Int64("chat_id", chatID),
		zap.String("user", update.Message.From.FirstName),
		zap.Int("length", utf8.RuneCountInString(update.Message.Text)),
	)

	reply, err := askLLM(ctx, update.Message.Text)
	if err != nil {
		logger.Error("llm request failed",
			zap.Int64("chat_id", chatID),
			zap.Error(err),
		)
		reply = "Sorry, the model is unavailable right now. Please try again later."
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   reply,
	})
	if err != nil {
		logger.Error("send message failed",
			zap.Int64("chat_id", chatID),
			zap.Error(err),
		)
	}
}
