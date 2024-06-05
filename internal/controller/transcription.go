package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/umarquez/goai_transcriber/internal/usecase"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// TranscriptionResponse represents the response for a successful transcription request.
type TranscriptionResponse struct {
	Transcription string `json:"transcription"`
}

// ErrorResponse represents the response for an error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TranscriptionController handles the transcriptions.
type TranscriptionController struct {
	TranscriptionUsecase *usecase.TranscriptionUsecase
}

// NewTranscriptionController creates a new TranscriptionController.
func NewTranscriptionController(usecase *usecase.TranscriptionUsecase) *TranscriptionController {
	return &TranscriptionController{TranscriptionUsecase: usecase}
}

// Transcribe godoc
// @Summary Transcribe an audio file
// @Description Upload an audio file and get the transcription
// @Tags transcriptions
// @Accept mpfd
// @Produce json
// @Param file formData file true "Audio file"
// @Success 200 {object} TranscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transcribe [post]
func (ctrl *TranscriptionController) Transcribe(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "file is required"})
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read file"})
		return
	}

	filePath := header.Filename
	if filepath.Ext(filePath) == ".m4a" {
		mp3FilePath, err := convertM4AToMP3(filePath, &buf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to convert M4A to MP3"})
			return
		}
		defer os.Remove(mp3FilePath) // Clean up the temporary MP3 file after processing

		filePath = mp3FilePath
		fileContent, err := os.Open(mp3FilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read converted MP3 file"})
			return
		}
		defer fileContent.Close()

		buf.Reset()
		if _, err := io.Copy(&buf, fileContent); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to copy converted MP3 file"})
			return
		}
	}

	transcription, err := ctrl.TranscriptionUsecase.Transcribe(filePath, &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to transcribe audio file"})
		return
	}

	c.JSON(http.StatusOK, TranscriptionResponse{Transcription: transcription.Text})
}

// convertM4AToMP3 converts an M4A file to MP3 format using FFmpeg.
func convertM4AToMP3(inputPath string, inputContent *bytes.Buffer) (string, error) {
	tmpFile, err := os.CreateTemp("", "input-*.m4a")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, inputContent); err != nil {
		return "", err
	}

	mp3FilePath := tmpFile.Name() + ".mp3"
	cmd := exec.Command("ffmpeg", "-i", tmpFile.Name(), mp3FilePath)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return mp3FilePath, nil
}
