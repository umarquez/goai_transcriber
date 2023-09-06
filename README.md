# goai_transcriber
Go based OpenAI audio files transcriber

This program automates the audio transcription process using OpenAI transcription API.

This will generate a txt files for each audio file named as `<file.wav>.txt` containing the transcription.

Checkout for costs on the OpenAI's site: [https://openai.com/pricing](https://openai.com/pricing)

# Before run:
Set the next environment variables:

```
OPENAI_TOKEN=<OpenAI API Token>
WORKING_PATH=<audio files directory>
```

# Run
```
go run ./cmd/app/main.go
```