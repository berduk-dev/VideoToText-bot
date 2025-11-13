package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

func New(baseURL string, timeout time.Duration) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is empty")
	}

	return &Client{
		BaseURL: baseURL,
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) Request(fullURL string) (*http.Response, error) {
	resp, err := http.Post(fullURL, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api returned %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
