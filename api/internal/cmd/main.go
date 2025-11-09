package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.POST("/transcribe", func(c *gin.Context) {
		youtubeLink := c.Query("link")
		type downloadReq struct {
			Url    string `json:"url"`
			Format string `json:"format"`
		}

		reqBody := downloadReq{
			Url:    youtubeLink,
			Format: "mp3",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			log.Println("failed to marshal request: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := http.Post(
			"http://yt-dl-api:8001/download",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			log.Println("failed to send request to yt-dl-api: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("yt-dl-api returned not success: ", resp.StatusCode)
			c.Status(http.StatusInternalServerError)
			return
		}

		out, err := os.Create("downloadedAudio.mp3")
		if err != nil {
			log.Println("failed to create file: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Println("failed to save file: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
