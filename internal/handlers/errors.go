package handlers

import (
	"errors"
	"log"
	"net/http"
	"testsmth/internal/apperrors"
)

func handleError(w http.ResponseWriter, err error) {
	log.Printf("handler error: %v", err)
	switch {
	case errors.Is(err, apperrors.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, apperrors.ErrPlayerNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, apperrors.ErrGameNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, apperrors.ErrAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, apperrors.ErrInsufficientBalance):
		http.Error(w, err.Error(), http.StatusConflict)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
