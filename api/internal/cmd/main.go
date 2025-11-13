package main

import (
	"github.com/berduk-dev/VideoToText-bot/api/internal/client/whisper"
	ytdl "github.com/berduk-dev/VideoToText-bot/api/internal/client/yt-dl"
	"github.com/berduk-dev/VideoToText-bot/api/internal/handler"
	"github.com/berduk-dev/VideoToText-bot/api/internal/service"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

const (
	timeout = 5 * time.Minute
)

func main() {
	r := gin.Default()

	whisperClient := whisper.New(os.Getenv("WHISPER_API_KEY"), os.Getenv("WHISPER_URL"), timeout)
	ytdlClient := ytdl.New(os.Getenv("YTDL_URL"), timeout)
	apiService := service.New(whisperClient, ytdlClient)
	apiHandler := handler.New(apiService)

	r.POST("/transcribe", apiHandler.TranscribeHandle)

	_ = r.Run()
}
