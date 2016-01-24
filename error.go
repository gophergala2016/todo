package main

import (
	"log"
	"net/http"
)

// AppError helps with http error handling
type AppError struct {
	Error   error
	Message string
	Code    int
}

func internalServerError(err error) *AppError {
	return &AppError{
		Error:   err,
		Message: http.StatusText(http.StatusInternalServerError),
		Code:    http.StatusInternalServerError,
	}
}

type appHandler func(w http.ResponseWriter, r *http.Request) *AppError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf(
			"handler error: status code %d, message %s, underlying error: %#v",
			e.Code,
			e.Message,
			e.Error,
		)
		http.Error(w, e.Message, e.Code)
	}
}
