package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

const (
	openaiTranscriptionsURL = "https://api.openai.com/v1/audio/transcriptions"
	openaiModel             = "whisper-1"
	responseFormat          = "json"
	maxFileSize             = 25 * 1024 * 1024 // 25MB
)

var SupportedFormats = map[string]bool{
	".mp3":  true,
	".mp4":  true,
	".mpeg": true,
	".mpga": true,
	".m4a":  true,
	".wav":  true,
	".webm": true,
}

type apiResponse struct {
	Text string `json:"text"`
}

type OpenAIClient struct {
	Token string
}

func NewOpenAIClient(token string) *OpenAIClient {
	return &OpenAIClient{Token: token}
}

func (c *OpenAIClient) Transcribe(filePath string, fileContent *bytes.Buffer) (string, error) {
	if err := c.validateFile(filePath, fileContent); err != nil {
		return "", err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part1, fileContent)
	if err != nil {
		return "", err
	}

	_ = writer.WriteField("model", openaiModel)
	_ = writer.WriteField("response_format", responseFormat)
	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, openaiTranscriptionsURL, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.Token))
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

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v\n%s", res.Status, body)
	}

	result := apiResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

func (c *OpenAIClient) validateFile(filePath string, fileContent *bytes.Buffer) error {
	if fileContent.Len() > maxFileSize {
		return fmt.Errorf("file size exceeds the 25MB limit")
	}

	ext := filepath.Ext(filePath)
	if !SupportedFormats[ext] {
		return fmt.Errorf("unsupported file format: %v", ext)
	}

	return nil
}
