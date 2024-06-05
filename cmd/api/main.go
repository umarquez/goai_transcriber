package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/umarquez/goai_transcriber/internal/api"
	"github.com/umarquez/goai_transcriber/internal/controller"
	"github.com/umarquez/goai_transcriber/internal/repository"
	"github.com/umarquez/goai_transcriber/internal/usecase"
	"github.com/umarquez/goai_transcriber/pkg/openai"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	token := viper.GetString("OPENAI_TOKEN")
	if token == "" {
		fmt.Println("Error: APP_OPENAI_TOKEN environment variable is not set")
		return
	}

	openAIClient := openai.NewOpenAIClient(token)
	transcriptionRepo := repository.NewTranscriptionRepository(openAIClient)
	transcriptionUsecase := &usecase.TranscriptionUsecase{TranscriptionRepo: transcriptionRepo}
	transcriptionController := controller.NewTranscriptionController(transcriptionUsecase)

	app := api.NewApi(transcriptionController)
	err := app.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
