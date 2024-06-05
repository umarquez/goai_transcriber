package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/umarquez/goai_transcriber/pkg/openai"
)

type App struct {
	OpenAIClient *openai.OpenAIClient
	SourceDir    string
}

func NewApp(apiToken, sourceDir string) *App {
	return &App{
		OpenAIClient: openai.NewOpenAIClient(apiToken),
		SourceDir:    sourceDir,
	}
}

func (a *App) Run() {
	audioFiles, err := getAudioFiles(a.SourceDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range audioFiles {
		content, err := readFileContent(file)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var tempFilePath string
		if filepath.Ext(file) == ".m4a" {
			tempFilePath, err = convertM4AToMP3(file)
			if err != nil {
				fmt.Println("Error converting M4A to MP3:", err)
				continue
			}
			defer os.Remove(tempFilePath)

			content, err = readFileContent(tempFilePath)
			if err != nil {
				fmt.Println("Error reading converted MP3 file:", err)
				continue
			}
		}

		textResult, err := a.OpenAIClient.Transcribe(tempFilePath, content)
		if err != nil {
			panic(err)
		}

		outputFilePath := spew.Sprintf("%v.txt", filepath.Join(a.SourceDir, filepath.Base(file)))
		err = os.WriteFile(outputFilePath, []byte(textResult), fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
