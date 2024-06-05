package controller

// TranscriptionResponse represents the response for a successful transcription request.
type TranscriptionResponse struct {
	Transcription string `json:"transcription"`
}

// ErrorResponse represents the response for an error.
type ErrorResponse struct {
	Error string `json:"error"`
}
