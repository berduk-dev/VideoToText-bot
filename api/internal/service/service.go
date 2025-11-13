package service

import (
	"context"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/api/internal/client/whisper"
	ytdl "github.com/berduk-dev/VideoToText-bot/api/internal/client/yt-dl"
)

type Service struct {
	WhisperClient *whisper.Client
	YtdlClient    *ytdl.Client
}

func New(whisperClient *whisper.Client, ytdlClient *ytdl.Client) Service {
	return Service{
		WhisperClient: whisperClient,
		YtdlClient:    ytdlClient,
	}
}

func (s *Service) TranscribeAudio(ctx context.Context, link string) (string, error) {
	audioData, err := s.YtdlClient.DownloadAudio(ctx, link)
	if err != nil {
		return "", fmt.Errorf("error client.DownloadAudio: %w", err)
	}

	transcribedText, err := s.WhisperClient.TranscribeAudio(ctx, audioData)
	if err != nil {
		return "", fmt.Errorf("error s.WhisperClient.TranscribeAudio: %w", err)
	}

	return transcribedText, nil
}
