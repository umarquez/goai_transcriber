package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarquez/goai_transcriber/internal/usecase"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// TranscriptionController handles the transcriptions.
type TranscriptionController struct {
	TranscriptionUsecase *usecase.TranscriptionUsecase
}

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
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
		ctrl.handleError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	buf, err := ctrl.readFileContent(file)
	if err != nil {
		ctrl.handleError(c, http.StatusInternalServerError, "failed to read file")
		return
	}

	if filepath.Ext(header.Filename) == ".m4a" {
		mp3FilePath, err := ctrl.convertM4AToMP3(header.Filename, buf)
		if err != nil {
			ctrl.handleError(c, http.StatusInternalServerError, "failed to convert M4A to MP3")
			return
		}
		defer os.Remove(mp3FilePath) // Clean up the temporary MP3 file after processing

		fileContent, err := os.Open(mp3FilePath)
		if err != nil {
			ctrl.handleError(c, http.StatusInternalServerError, "failed to read converted MP3 file")
			return
		}
		defer fileContent.Close()

		buf.Reset()
		if _, err := io.Copy(buf, fileContent); err != nil {
			ctrl.handleError(c, http.StatusInternalServerError, "failed to copy converted MP3 file")
			return
		}
	}

	transcription, err := ctrl.TranscriptionUsecase.Transcribe(header.Filename, buf)
	if err != nil {
		ctrl.handleError(c, http.StatusInternalServerError, "failed to transcribe audio file")
		return
	}

	c.JSON(http.StatusOK, TranscriptionResponse{Transcription: transcription.Text})
}

func (ctrl *TranscriptionController) handleError(c *gin.Context, status int, message string) {
	log.Error(message)
	c.JSON(status, ErrorResponse{Error: message})
}

func (ctrl *TranscriptionController) readFileContent(file multipart.File) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return nil, err
	}
	return &buf, nil
}

func (ctrl *TranscriptionController) convertM4AToMP3(inputPath string, inputContent *bytes.Buffer) (string, error) {
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
