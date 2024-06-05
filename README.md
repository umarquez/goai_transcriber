
# GoAI Transcriber

GoAI Transcriber is a tool that uses OpenAI's Whisper model to transcribe audio files. It supports various audio formats including `.m4a` and converts them to `.mp3` before processing if necessary.

## Features

- Transcribes audio files using OpenAI's Whisper model.
- Supports the following audio formats: `.mp3`, `.mp4`, `.mpeg`, `.mpga`, `.m4a`, `.wav`, `.webm`.
- Automatically converts `.m4a` files to `.mp3` for processing.
- Provides a REST API for uploading and transcribing audio files.
- Includes Swagger documentation for easy API exploration.

## Project Structure

```
.
├── cmd
│   ├── api
│   │   └── main.go
│   └── app
│       └── main.go
├── deployment
│   ├── Dockerfile
│   ├── Dockerfile.api
│   ├── docker-compose.yml
│   ├── docker-compose.api.yml
│   ├── terraform
│   │   ├── main.tf
│   │   └── variables.tf
├── internal
│   ├── api
│   │   └── api.go
│   ├── app
│   │   ├── app.go
│   │   └── functions.go
│   ├── controller
│   │   └── transcription.go
│   ├── entity
│   │   └── transcription.go
│   ├── repository
│   │   └── transcription.go
│   └── usecase
│       └── transcription.go
├── pkg
│   └── openai
│       └── api.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── README.md
└── .env
```

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/umarquez/goai_transcriber.git
   cd goai_transcriber
   ```

2. **Set up environment variables:**
   Create a `.env` file and add your OpenAI token:
   ```sh
   APP_OPENAI_TOKEN=your_openai_token
   APP_WORKING_PATH=./audios
   ```

## Running the Application

### Running the Application Locally

#### Standard Application

1. **Build and run the application:**
   ```sh
   go build -o bin/app cmd/app/main.go
   ./bin/app
   ```

2. **Transcribe audio files:**
   - The application reads the content of the `./audios` directory.
   - If there are `.m4a` files, they are converted to `.mp3` due to an error from the OpenAI API processing `.m4a` files.
   - The transcription result is written to the same directory with a `.txt` extension.

#### API Version

1. **Generate Swagger documentation:**
   ```sh
   swag init --parseDependency --parseInternal -g cmd/api/main.go -o ./docs
   ```

2. **Build and run the API:**
   ```sh
   go build -o bin/api cmd/api/main.go
   ./bin/api
   ```

3. **Transcribe audio files via API:**
   Use a tool like `curl` or Postman to send a `POST` request to the `/transcribe` endpoint with your audio file.

   Example `curl` command:
   ```sh
   curl -X POST "http://localhost:8080/transcribe" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "file=@path/to/your/audiofile.m4a"
   ```

### Running the Application Using Docker

#### Standard Application

1. **Build and run the application using Docker:**
   ```sh
   docker-compose -f deployment/docker-compose.yml up --build
   ```

2. **Transcribe audio files:**
   Place your audio files in the `audios` directory and the application will automatically process and transcribe them.

#### API Version

1. **Build and run the API using Docker:**
   ```sh
   docker-compose -f deployment/docker-compose.api.yml up --build
   ```

2. **Transcribe audio files via API:**
   Use a tool like `curl` or Postman to send a `POST` request to the `/transcribe` endpoint with your audio file.

   Example `curl` command:
   ```sh
   curl -X POST "http://localhost:8080/transcribe" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "file=@path/to/your/audiofile.m4a"
   ```

## API Documentation

The API is documented using Swagger. Once the application is running, you can access the documentation at:
```
http://localhost:8080/swagger/index.html
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
