package handler

import (
	"github.com/berduk-dev/VideoToText-bot/api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) TranscribeHandle(c *gin.Context) {
	link := c.Query("link")

	transcribedText, err := h.Service.TranscribeAudio(c, link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text": transcribedText,
	})
}
