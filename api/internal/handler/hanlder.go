package handler

import (
	"github.com/berduk-dev/VideoToText-bot/api/internal/client/whisper"
	"github.com/berduk-dev/VideoToText-bot/api/internal/client/yt-dl"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	WhisperClient *whisper.Client
	YtdlClient    *yt_dl.Client
}

func New(whisperClient *whisper.Client, ytdlClient *yt_dl.Client) Handler {
	return Handler{
		WhisperClient: whisperClient,
		YtdlClient:    ytdlClient,
	}
}

func (h *Handler) TranscribeHandle(c *gin.Context) {
	link := c.Query("link")

	audioData, err := h.YtdlClient.DownloadAudio(c, link)
	if err != nil {
		log.Println("error client.DownloadAudio:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transcribedText, err := h.WhisperClient.TranscribeAudio(c, audioData)
	if err != nil {
		log.Println("error client.TranscribeAudio:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text": transcribedText,
	})
}
