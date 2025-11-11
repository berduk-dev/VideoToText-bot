package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/api/internal/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func DownloadAudio(c *gin.Context) ([]byte, error) {
	youtubeLink := c.Query("link")

	reqBody := model.DownloadReq{
		Url:    youtubeLink,
		Format: "mp3",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		"http://yt-dl-api:8001/download",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to yt-dl-api: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yt-dl-api returned not success: %d", resp.StatusCode)
	}

	audioData, _ := io.ReadAll(resp.Body)

	return audioData, nil
}
