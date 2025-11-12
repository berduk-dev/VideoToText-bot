package whisper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/api/internal/model"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	ApiKey     string
	BaseURL    string
	HttpClient *http.Client
}

func New(apiKey string, baseUrl string, timeout time.Duration) *Client {
	return &Client{
		ApiKey:  apiKey,
		BaseURL: baseUrl,
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	suffix := "/v1/audio/transcriptions"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile("file", "audio.mp3")
	if err != nil {
		return "", fmt.Errorf("error CreateFormFile: %w", err)
	}

	_, err = io.Copy(fileWriter, bytes.NewReader(audioData))
	if err != nil {
		return "", fmt.Errorf("error Copy - TranscribeAudio: %w", err)
	}

	_ = writer.WriteField("model", "stt-openai/whisper-v3-turbo")
	_ = writer.WriteField("response_format", "json")

	_ = writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+suffix, &body)
	if err != nil {
		return "", fmt.Errorf("error api request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error client.Do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading resp body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vsegpt API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var transcribeResp model.TranscribeResp

	err = json.Unmarshal(respBody, &transcribeResp)
	if err != nil {
		return "", fmt.Errorf("error parsing json response: %w (body: %s)", err, string(respBody))
	}

	return transcribeResp.Text, nil
}
