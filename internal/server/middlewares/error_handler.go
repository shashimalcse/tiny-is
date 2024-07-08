package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, message string) APIError {
	return APIError{Status: status, Message: message}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func ErrorMiddleware(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %v", r)
				sendErrorResponse(w, NewAPIError(http.StatusInternalServerError, "Internal server error"))
			}
		}()

		if err := next(w, r); err != nil {
			sendErrorResponse(w, err)
		}
	}
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	apiErr, ok := err.(APIError)
	if !ok {
		apiErr = NewAPIError(http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)
	json.NewEncoder(w).Encode(apiErr)
}
