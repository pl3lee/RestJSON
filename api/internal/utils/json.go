package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	RespondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func DecodeResponse[T any](resp *http.Response, result *T) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
