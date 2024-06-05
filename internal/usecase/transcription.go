package usecase

import (
	"bytes"
	"github.com/umarquez/goai_transcriber/internal/entity"
	"github.com/umarquez/goai_transcriber/internal/repository"
)

type Transcription interface {
	Transcribe(filePath string, fileContent *bytes.Buffer) (entity.Transcription, error)
}

type transcription struct {
	TranscriptionRepo repository.Transcription
}

func (u *transcription) Transcribe(filePath string, fileContent *bytes.Buffer) (entity.Transcription, error) {
	return u.TranscriptionRepo.Transcribe(filePath, fileContent)
}
