package client

import (
	"encoding/json"
	"fmt"
	"github.com/azdaev/yt-transcribe-bot/api/internal/model"
	"os/exec"
)

func TranscribeAudio() (string, error) {
	cmd := exec.Command(
		"curl",
		"-s", // отключает информацию о процессе (silent mode)
		"https://api.vsegpt.ru/v1/audio/transcriptions",
		"-H", "Authorization: Bearer sk-or-vv-0479dbef3c94679a5be21969cd273b21b25dd71bb2088f1b54182eddd2f33ab6",
		"-F", "file=@downloadedAudio.mp3",
		"-F", "model=stt-openai/whisper-v3-turbo",
		"-F", "response_format=json",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error command - TranscribeAudio: %w", err)
	}

	var resp model.WhisperResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	return resp.Text, nil
}
