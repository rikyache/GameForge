package handlers

import (
	"encoding/json"
	"net/http"
	"testsmth/internal/auth"
	"testsmth/internal/middleware"
	"testsmth/internal/models"
	"testsmth/internal/service"
	"time"
)

type AuthHandler struct {
	Service   *service.AuthService
	blacklist auth.TokenBlacklist
}

func NewAuthHandler(service *service.AuthService, blacklist auth.TokenBlacklist) *AuthHandler {
	return &AuthHandler{
		Service:   service,
		blacklist: blacklist,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.Register(r.Context(), req); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	}); err != nil {
		handleError(w, err)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenID, ok := middleware.TokenIDFromContext(ctx)
	if !ok {
		http.Error(w, "token id not found", http.StatusUnauthorized)
		return
	}

	expiresAt, ok := middleware.TokenExpiresAtFromContext(ctx)
	if !ok {
		http.Error(w, "token expires at not found", http.StatusUnauthorized)
		return
	}

	ttl := time.Until(expiresAt)

	if ttl < 0 {
		http.Error(w, "token expired", http.StatusUnauthorized)
		return
	}

	err := h.blacklist.Revoke(ctx, tokenID, ttl)
	if err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
