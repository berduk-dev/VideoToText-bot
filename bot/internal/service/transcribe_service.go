package service

import (
	"context"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/bot/internal/client"
)

type Service struct {
	Client *client.Client
}

func New(client *client.Client) Service {
	return Service{
		Client: client,
	}
}

func (s *Service) GetTranscription(ctx context.Context, link string) (string, error) {
	text, err := s.Client.TranscribeAudio(ctx, link)
	if err != nil {
		return "", fmt.Errorf("error s.Client.TranscribeAudio: %w", err)
	}

	return text, nil
}
