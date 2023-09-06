package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"

	"github.com/umarquez/goai_transcriber/pkg/openai"
)

const workingPath = "WORKING_PATH"

func Run() {
	sourceDir := os.Getenv(workingPath)
	wavFiles, err := getWavFiles(sourceDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the found files
	for _, file := range wavFiles {
		//fmt.Println(filepath.Base(file))
		textResult, err := openai.Transcribe(file)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(spew.Sprintf("%v.txt", file), []byte(textResult), fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func getWavFiles(folderPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".wav" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
