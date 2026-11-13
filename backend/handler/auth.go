package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/670710731/Pizzahub/backend/auth"
	"github.com/670710731/Pizzahub/backend/mw"
	"github.com/670710731/Pizzahub/backend/store"
)

type authHandler struct {
	store  *store.Store
	tokens *auth.TokenService
	logger *slog.Logger
}

type credentials struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Email = strings.TrimSpace(input.Email)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email || len(input.Email) > 254 {
		writeError(w, http.StatusBadRequest, "valid email is required")
		return
	}
	if input.DisplayName == "" || len(input.DisplayName) > 100 {
		writeError(w, http.StatusBadRequest, "display_name must be 1-100 characters")
		return
	}
	passwordHash, err := auth.HashPassword(input.Password)
	if errors.Is(err, auth.ErrInvalidPassword) {
		writeError(w, http.StatusBadRequest, "password must be 8-72 bytes")
		return
	}
	if err != nil {
		h.logger.Error("hash password", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	user, err := h.store.CreateUser(r.Context(), input.Email, passwordHash, input.DisplayName)
	if errors.Is(err, store.ErrConflict) {
		writeError(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		h.logger.Error("register user", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.writeToken(w, http.StatusCreated, user)
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if !decodeJSON(w, r, &input) {
		return
	}
	user, err := h.store.UserByEmail(r.Context(), strings.TrimSpace(input.Email))
	if errors.Is(err, store.ErrNotFound) || (err == nil && !auth.VerifyPassword(user.PasswordHash, input.Password)) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		h.logger.Error("find login user", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.writeToken(w, http.StatusOK, user)
}

func (h *authHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	identity, ok := mw.CurrentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	user, err := h.store.UserByID(r.Context(), identity.UserID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "account is unavailable")
		return
	}
	if err != nil {
		h.logger.Error("find user for password change", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !auth.VerifyPassword(user.PasswordHash, input.CurrentPassword) {
		writeError(w, http.StatusUnauthorized, "current password is incorrect")
		return
	}
	newHash, err := auth.HashPassword(input.NewPassword)
	if errors.Is(err, auth.ErrInvalidPassword) {
		writeError(w, http.StatusBadRequest, "new password must be 8-72 bytes")
		return
	}
	if err != nil {
		h.logger.Error("hash new password", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if err := h.store.UpdatePasswordHash(r.Context(), identity.UserID, user.PasswordHash, newHash); errors.Is(err, store.ErrConflict) {
		writeError(w, http.StatusConflict, "password changed concurrently; please try again")
		return
	} else if err != nil {
		h.logger.Error("update password", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *authHandler) writeToken(w http.ResponseWriter, status int, user store.User) {
	token, expiresAt, err := h.tokens.Issue(user.ID, user.Role)
	if err != nil {
		h.logger.Error("issue access token", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, status, map[string]any{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_at":   expiresAt.Format(time.RFC3339),
		"user":         user,
	})
}
