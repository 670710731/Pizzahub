package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/670710731/Pizzahub/backend/store"
)

type adminHandler struct {
	store  *store.Store
	logger *slog.Logger
}

func (h *adminHandler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.ListCategories(r.Context())
	if err != nil {
		h.logger.Error("list categories", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

func (h *adminHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var input store.Category
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(input.Slug)
	if input.Name == "" || input.Slug == "" {
		writeError(w, http.StatusBadRequest, "name and slug are required")
		return
	}
	category, err := h.store.CreateCategory(r.Context(), input)
	if errors.Is(err, store.ErrConflict) {
		writeError(w, http.StatusConflict, "category slug already exists")
		return
	}
	if err != nil {
		h.logger.Error("create category", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, category)
}

func (h *adminHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	var input store.Category
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(input.Slug)
	if input.Name == "" || input.Slug == "" {
		writeError(w, http.StatusBadRequest, "name and slug are required")
		return
	}
	err = h.store.UpdateCategory(r.Context(), id, input)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	if errors.Is(err, store.ErrConflict) {
		writeError(w, http.StatusConflict, "category slug already exists")
		return
	}
	if err != nil {
		h.logger.Error("update category", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *adminHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.store.DeleteCategory(r.Context(), id); errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	} else if err != nil {
		h.logger.Error("delete category", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *adminHandler) listReviews(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := pagination(w, r)
	if !ok {
		return
	}
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending"
	}
	reviews, total, err := h.store.ListAdminReviews(r.Context(), status, page, pageSize)
	if errors.Is(err, store.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "status must be pending, published, or rejected")
		return
	}
	if err != nil {
		h.logger.Error("list admin reviews", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, paginated(reviews, page, pageSize, total))
}

func (h *adminHandler) savePizza(w http.ResponseWriter, r *http.Request) {
	var id int64
	if r.Method == http.MethodPut {
		parsed, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "invalid pizza id")
			return
		}
		id = parsed
	}
	var input store.PizzaInput
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.TrimSpace(input.Slug)
	if input.Name == "" || input.Slug == "" {
		writeError(w, http.StatusBadRequest, "name and slug are required")
		return
	}
	if len(input.Name) > 150 || len(input.Slug) > 160 || len(input.Description) > 10000 {
		writeError(w, http.StatusBadRequest, "pizza fields exceed maximum length")
		return
	}
	for _, image := range input.Images {
		if strings.TrimSpace(image.URL) == "" || image.SortOrder < 0 {
			writeError(w, http.StatusBadRequest, "each image requires image_url and non-negative sort_order")
			return
		}
	}
	for _, variant := range input.Variants {
		if strings.TrimSpace(variant.Name) == "" || strings.TrimSpace(variant.SizeLabel) == "" || variant.Price == "" || variant.SortOrder < 0 {
			writeError(w, http.StatusBadRequest, "variants require name, size_label, price, and non-negative sort_order")
			return
		}
	}
	savedID, err := h.store.SavePizza(r.Context(), id, input)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "pizza not found")
		return
	}
	if errors.Is(err, store.ErrConflict) {
		writeError(w, http.StatusConflict, "pizza slug or variant size already exists")
		return
	}
	if errors.Is(err, store.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid pizza status, category, image order, or variant value")
		return
	}
	if err != nil {
		h.logger.Error("save pizza", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	status := http.StatusCreated
	if id != 0 {
		status = http.StatusOK
	}
	writeJSON(w, status, map[string]int64{"id": savedID})
}

func (h *adminHandler) deletePizza(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid pizza id")
		return
	}
	if err := h.store.DeletePizza(r.Context(), id); errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "pizza not found")
		return
	} else if err != nil {
		h.logger.Error("delete pizza", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
