package service

import (
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/model"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Serivce struct {
}

func New() Serivce {
	return Serivce{}
}

func (s *Serivce) GetTranscription(link string) (string, error) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		return "", fmt.Errorf("API_URL is empty")
	}

	fullURL := fmt.Sprintf("%s/transcribe?link=%s", apiURL, url.QueryEscape(link))

	resp, err := http.Post(fullURL, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api returned %d: %s", resp.StatusCode, string(body))
	}

	var result model.TranscribeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}

	return result.Text, nil
}
