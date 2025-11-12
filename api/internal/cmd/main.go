package main

import (
	"github.com/berduk-dev/VideoToText-bot/api/internal/client/whisper"
	yt_dl "github.com/berduk-dev/VideoToText-bot/api/internal/client/yt-dl"
	"github.com/berduk-dev/VideoToText-bot/api/internal/handler"
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
	ytdlClient := yt_dl.New(os.Getenv("YTDL_URL"), timeout)
	apiHandler := handler.New(whisperClient, ytdlClient)

	r.POST("/transcribe", apiHandler.TranscribeHandle)

	_ = r.Run()
}
