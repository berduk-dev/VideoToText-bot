package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/model"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string, timeout time.Duration) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is empty")
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) TranscribeAudio(ctx context.Context, link string) (string, error) {
	fullURL := fmt.Sprintf("%s/transcribe?link=%s", c.baseURL, url.QueryEscape(link))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to go-api: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api returned %d: %s", resp.StatusCode, string(body))
	}

	var result model.TranscribeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}
	defer resp.Body.Close()

	return result.Text, nil
}
