package yt_dl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/api/internal/model"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

func New(baseUrl string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseUrl,
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) DownloadAudio(ctx context.Context, youtubeLink string) ([]byte, error) {
	suffix := "/download"
	reqBody := model.DownloadReq{
		Url:    youtubeLink,
		Format: "mp3",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+suffix, bytes.NewReader(body))

	// Выполняем запрос
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to yt-dl-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yt-dl-api returned not success: %d", resp.StatusCode)
	}

	audioData, _ := io.ReadAll(resp.Body)

	return audioData, nil
}
