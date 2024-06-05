package repository

import (
	"bytes"
	"github.com/umarquez/goai_transcriber/internal/entity"
	"github.com/umarquez/goai_transcriber/pkg/openai"
)

type TranscriptionRepository interface {
	Transcribe(filePath string, fileContent *bytes.Buffer) (entity.Transcription, error)
}

type transcriptionRepository struct {
	openAIClient *openai.OpenAIClient
}

func NewTranscriptionRepository(client *openai.OpenAIClient) TranscriptionRepository {
	return &transcriptionRepository{openAIClient: client}
}

func (r *transcriptionRepository) Transcribe(filePath string, fileContent *bytes.Buffer) (entity.Transcription, error) {
	text, err := r.openAIClient.Transcribe(filePath, fileContent)
	if err != nil {
		return entity.Transcription{}, err
	}
	return entity.Transcription{
		ID:     filePath,
		Text:   text,
		Format: "text",
	}, nil
}
