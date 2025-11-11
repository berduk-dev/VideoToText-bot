package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/berduk-dev/VideoToText-bot/api/internal/model"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func TranscribeAudio(audioData []byte) (string, error) {
	url := "https://api.vsegpt.ru/v1/audio/transcriptions"

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

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("error api request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VSEGPT_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
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
