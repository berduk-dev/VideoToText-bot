package main

import (
	"github.com/azdaev/yt-transcribe-bot/api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiHandler := handler.New()

	r.POST("/transcribe", apiHandler.TranscribeHandle)

	_ = r.Run()
}
