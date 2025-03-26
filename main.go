package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Config struct {
	Bot struct {
		Token string `toml:"Token"`
	} `toml:"Bot"`
}

func loadConfig(filename string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig("cfg.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	token := config.Bot.Token
	if token == "" {
		log.Fatal("Bot token is not set in config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handleAllUpdates),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	b.Start(ctx)
}

// такая вообше должна быть логика ?
func handleAllUpdates(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		if update.Message.Text == "/start" {
			handleStartCommand(ctx, b, update)
			return
		}
	} else if update.InlineQuery != nil {
		handleInlineQuery(ctx, b, update)
	}
}

func handleStartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello world",
	})
}

func handleInlineQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	results := []models.InlineQueryResult{
		&models.InlineQueryResultArticle{
			ID:    "1",
			Title: "Start",
			InputMessageContent: &models.InputTextMessageContent{
				MessageText: "Hello world",
			},
		},
	}

	b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
	})
}
