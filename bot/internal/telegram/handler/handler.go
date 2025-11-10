package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Service interface {
	GetTranscription(link string) (string, error)
}

type Handler struct {
	Service
}

func New(service Service) Handler {
	return Handler{
		Service: service,
	}
}

func (h *Handler) HandleYouTubeLink(bot *tgbotapi.BotAPI, chatID int64, link string) {
	msg := tgbotapi.NewMessage(chatID, "⏳ Расшифровываю аудио, подожди немного...")
	_, _ = bot.Send(msg)

	text, err := h.GetTranscription(link)
	if err != nil {
		log.Println("error getTranscription:", err)
		msg := tgbotapi.NewMessage(chatID, "❌ Произошла ошибка. Попробуйте позже!")
		_, _ = bot.Send(msg)
		return
	}
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("🗣️ Расшифровка:\n%s", text))

	_, _ = bot.Send(msg)
}
