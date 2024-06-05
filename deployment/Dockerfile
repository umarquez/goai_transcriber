FROM golang:1.22

# Install FFmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/app/main.go

# Command to run the executable
CMD ["./main"]
