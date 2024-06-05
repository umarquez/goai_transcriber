package usecase

import (
	"bytes"
	"github.com/umarquez/goai_transcriber/internal/entity"
	"github.com/umarquez/goai_transcriber/internal/repository"
)

type TranscriptionUsecase struct {
	TranscriptionRepo repository.TranscriptionRepository
}

func (u *TranscriptionUsecase) Transcribe(filePath string, fileContent *bytes.Buffer) (entity.Transcription, error) {
	return u.TranscriptionRepo.Transcribe(filePath, fileContent)
}
