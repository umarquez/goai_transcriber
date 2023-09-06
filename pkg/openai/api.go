package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

const (
	openaiTranscriptionsURL = "https://api.openai.com/v1/audio/transcriptions"
	openaiTokenVar          = "OPENAI_TOKEN"
	openaiModel             = "whisper-1"
)

func Transcribe(filePath string) (string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part1, file)
	if err != nil {
		return "", err
	}
	_ = writer.WriteField("model", openaiModel)
	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, openaiTranscriptionsURL, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", spew.Sprintf("Bearer %v", os.Getenv(openaiTokenVar)))
	//req.Header.Add("Cookie", "__cf_bm=bVf7TDsPMy5gY4fAVWa8iK8li.MYWERnPFyrL8dOQd8-1694022164-0-ASazMZYHiGnE8ur9NNwmlUGHaAdeMu6Yt1uTlR2L9J1+fRyUk5nBNu6XcTWvKAibx/bkapeyeAo0mmqWRX8me/s=")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Text string `json:"text"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}
