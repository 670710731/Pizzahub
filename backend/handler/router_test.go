package handler

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/670710731/Pizzahub/backend/auth"
	"github.com/670710731/Pizzahub/backend/store"
)

func TestRouterHealthAndProtectedReviewRoute(t *testing.T) {
	tokens, err := auth.NewTokenService(strings.Repeat("s", 32), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(&store.Store{}, tokens, logger, nil)

	health := httptest.NewRecorder()
	router.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", health.Code, http.StatusOK)
	}

	review := httptest.NewRecorder()
	router.ServeHTTP(review, httptest.NewRequest(http.MethodPost, "/api/v1/pizzas/margherita/reviews", strings.NewReader(`{"rating":5}`)))
	if review.Code != http.StatusUnauthorized {
		t.Fatalf("review status = %d, want %d", review.Code, http.StatusUnauthorized)
	}
}
