package main

import (
	"context"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/client"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/service"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/telegram/handler"
	"log"
	"time"

	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	timeout = 5 * time.Minute
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	botClient, err := client.New(os.Getenv("GO_API_URL"), timeout)
	if err != nil {
		log.Fatalf("failed creating client: %v", err)
	}

	botService := service.New(botClient)
	botHandler := handler.New(&botService)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		link := update.Message.Text
		if strings.Contains(link, "youtube.com") {
			botHandler.HandleYouTubeLink(context.Background(), bot, update.Message.Chat.ID, link)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправь ссылку на YouTube 🎥")
			_, _ = bot.Send(msg)
		}
	}
}
