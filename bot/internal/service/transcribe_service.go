package service

import (
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/client"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/model"
	"net/url"
)

type Service struct {
	Client *client.Client
}

func New(client *client.Client) Service {
	return Service{
		Client: client,
	}
}

func (s *Service) GetTranscription(link string) (string, error) {
	fullURL := fmt.Sprintf("%s/transcribe?link=%s", s.Client.BaseURL, url.QueryEscape(link))

	resp, err := s.Client.Request(fullURL)
	if err != nil {
		return "", fmt.Errorf("error s.Client.Request: %w", err)
	}

	var result model.TranscribeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}
	defer resp.Body.Close()

	return result.Text, nil
}
