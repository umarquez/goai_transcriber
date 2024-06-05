package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/umarquez/goai_transcriber/internal/app"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	token := viper.GetString("OPENAI_TOKEN")
	if token == "" {
		fmt.Println("Error: APP_OPENAI_TOKEN environment variable is not set")
		return
	}

	workingPath := viper.GetString("WORKING_PATH")
	if workingPath == "" {
		fmt.Println("Error: APP_WORKING_PATH environment variable is not set")
		return
	}

	app := app.NewApp(token, workingPath)
	app.Run()
}
