package model

type WhisperResponse struct {
	Text string `json:"text"`
}

type DownloadReq struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

type TranscribeResp struct {
	Text string `json:"text"`
}
