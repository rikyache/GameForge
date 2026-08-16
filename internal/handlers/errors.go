package handlers

import (
	"errors"
	"net/http"
	"testsmth/internal/apperrors"
)

func handleError(w http.ResponseWriter, err error) {
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
