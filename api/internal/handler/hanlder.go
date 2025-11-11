package handler

import (
	"github.com/berduk-dev/VideoToText-bot/api/internal/client"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
}

func New() Handler {
	return Handler{}
}

func (h *Handler) TranscribeHandle(c *gin.Context) {
	audioData, err := client.DownloadAudio(c)
	if err != nil {
		log.Println("error client.DownloadAudio:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transcribedText, err := client.TranscribeAudio(audioData)
	if err != nil {
		log.Println("error client.TranscribeAudio:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text": transcribedText,
	})
}
