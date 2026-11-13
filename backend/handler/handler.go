package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/670710731/Pizzahub/backend/auth"
	"github.com/670710731/Pizzahub/backend/mw"
	"github.com/670710731/Pizzahub/backend/store"
)

func NewRouter(database *store.Store, tokens *auth.TokenService, logger *slog.Logger, allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()
	authHandlers := &authHandler{store: database, tokens: tokens, logger: logger}
	reviewHandlers := &reviewHandler{store: database, logger: logger}
	adminHandlers := &adminHandler{store: database, logger: logger}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/v1/auth/register", authHandlers.register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandlers.login)
	mux.Handle("POST /api/v1/auth/change-password", mw.RequireAuth(tokens, false, http.HandlerFunc(authHandlers.changePassword)))
	mux.HandleFunc("GET /api/v1/pizzas", listPizzas(database, logger))
	mux.HandleFunc("GET /api/v1/pizzas/{slug}", getPizza(database, logger))
	mux.HandleFunc("GET /api/v1/pizzas/{slug}/reviews", reviewHandlers.list)
	mux.Handle("POST /api/v1/pizzas/{slug}/reviews", mw.RequireAuth(tokens, false, http.HandlerFunc(reviewHandlers.create)))
	mux.HandleFunc("GET /api/v1/categories", adminHandlers.listCategories)

	mux.Handle("GET /api/v1/admin/categories", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.listCategories)))
	mux.Handle("POST /api/v1/admin/categories", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.createCategory)))
	mux.Handle("PUT /api/v1/admin/categories/{id}", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.updateCategory)))
	mux.Handle("DELETE /api/v1/admin/categories/{id}", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.deleteCategory)))
	mux.Handle("POST /api/v1/admin/pizzas", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.savePizza)))
	mux.Handle("PUT /api/v1/admin/pizzas/{id}", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.savePizza)))
	mux.Handle("DELETE /api/v1/admin/pizzas/{id}", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.deletePizza)))
	mux.Handle("GET /api/v1/admin/reviews", mw.RequireAuth(tokens, true, http.HandlerFunc(adminHandlers.listReviews)))
	mux.Handle("PATCH /api/v1/admin/reviews/{id}", mw.RequireAuth(tokens, true, http.HandlerFunc(reviewHandlers.moderate)))

	return mw.Logging(logger, mw.CORS(allowedOrigins, mux))
}

func listPizzas(database *store.Store, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize, ok := pagination(w, r)
		if !ok {
			return
		}
		items, total, err := database.ListPizzas(r.Context(), page, pageSize, r.URL.Query().Get("category"), r.URL.Query().Get("sort"))
		if err != nil {
			if errors.Is(err, store.ErrInvalidInput) {
				writeError(w, http.StatusBadRequest, "sort must be newest, name, or rating")
				return
			}
			logger.Error("list pizzas", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(w, http.StatusOK, paginated(items, page, pageSize, total))
	}
}

func getPizza(database *store.Store, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pizza, err := database.PizzaBySlug(r.Context(), r.PathValue("slug"))
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "pizza not found")
			return
		}
		if err != nil {
			logger.Error("get pizza", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(w, http.StatusOK, pizza)
	}
}

type pageResponse[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

func paginated[T any](items []T, page, pageSize int, total int64) pageResponse[T] {
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	return pageResponse[T]{Items: items, Page: page, PageSize: pageSize, TotalItems: total, TotalPages: totalPages}
}

func pagination(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	page, pageSize := 1, 20
	var err error
	if raw := r.URL.Query().Get("page"); raw != "" {
		page, err = strconv.Atoi(raw)
		if err != nil || page < 1 {
			writeError(w, http.StatusBadRequest, "page must be a positive integer")
			return 0, 0, false
		}
	}
	if raw := r.URL.Query().Get("page_size"); raw != "" {
		pageSize, err = strconv.Atoi(raw)
		if err != nil || pageSize < 1 || pageSize > 100 {
			writeError(w, http.StatusBadRequest, "page_size must be between 1 and 100")
			return 0, 0, false
		}
	}
	if page > 1_000_000 {
		writeError(w, http.StatusBadRequest, "page is too large")
		return 0, 0, false
	}
	return page, pageSize, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON value")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("write JSON response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
