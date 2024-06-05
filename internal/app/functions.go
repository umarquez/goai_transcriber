package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/umarquez/goai_transcriber/pkg/openai"
)

func getAudioFiles(folderPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && openai.SupportedFormats[filepath.Ext(path)] {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// readFileContent reads the content of the file at the given path and returns it as a bytes buffer
func readFileContent(filePath string) (*bytes.Buffer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// convertM4AToMP3 converts an M4A file to MP3 format using ffmpeg and stores it in a temporary file
func convertM4AToMP3(inputPath string) (string, error) {
	tempFile, err := os.CreateTemp("", "temp_audio_*.mp3")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Construct the FFmpeg command
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, tempFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the FFmpeg command and capture any errors
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, %v", err, stderr.String())
	}

	// Check if the temp file is empty
	fileInfo, err := os.Stat(tempFile.Name())
	if err != nil {
		return "", err
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("converted MP3 file is empty")
	}

	return tempFile.Name(), nil
}
