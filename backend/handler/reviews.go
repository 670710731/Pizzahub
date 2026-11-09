package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/670710731/Pizzahub/backend/mw"
	"github.com/670710731/Pizzahub/backend/store"
)

type reviewHandler struct {
	store  *store.Store
	logger *slog.Logger
}

type reviewInput struct {
	Rating int    `json:"rating"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (h *reviewHandler) list(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := pagination(w, r)
	if !ok {
		return
	}
	reviews, total, err := h.store.ListReviews(r.Context(), r.PathValue("slug"), page, pageSize)
	if err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "pizza not found")
		return
	}
	if err != nil {
		h.logger.Error("list reviews", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, paginated(reviews, page, pageSize, total))
}

func (h *reviewHandler) create(w http.ResponseWriter, r *http.Request) {
	identity, ok := mw.CurrentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var input reviewInput
	if !decodeJSON(w, r, &input) {
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		writeError(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}
	if len(input.Title) > 150 || len(input.Body) > 5000 {
		writeError(w, http.StatusBadRequest, "review title or body is too long")
		return
	}
	review, err := h.store.SaveReview(r.Context(), r.PathValue("slug"), identity.UserID, input.Rating, input.Title, input.Body)
	if err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "pizza not found")
		return
	}
	if err != nil {
		h.logger.Error("save review", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, review)
}

func (h *reviewHandler) moderate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid review id")
		return
	}
	var input struct {
		Status string `json:"status"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	if input.Status != "pending" && input.Status != "published" && input.Status != "rejected" {
		writeError(w, http.StatusBadRequest, "status must be pending, published, or rejected")
		return
	}
	if err := h.store.ModerateReview(r.Context(), id, input.Status); err == store.ErrNotFound {
		writeError(w, http.StatusNotFound, "review not found")
		return
	} else if err != nil {
		h.logger.Error("moderate review", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
