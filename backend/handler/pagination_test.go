package handler

import (
	"net/http/httptest"
	"testing"
)

func TestPaginationDefaultsAndLimits(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/v1/pizzas", nil)
	recorder := httptest.NewRecorder()
	page, pageSize, ok := pagination(recorder, request)
	if !ok || page != 1 || pageSize != 20 {
		t.Fatalf("unexpected defaults: page=%d page_size=%d ok=%v", page, pageSize, ok)
	}

	request = httptest.NewRequest("GET", "/api/v1/pizzas?page=2&page_size=100", nil)
	recorder = httptest.NewRecorder()
	page, pageSize, ok = pagination(recorder, request)
	if !ok || page != 2 || pageSize != 100 {
		t.Fatalf("unexpected pagination: page=%d page_size=%d ok=%v", page, pageSize, ok)
	}
}

func TestPaginationRejectsInvalidValues(t *testing.T) {
	for _, target := range []string{
		"/api/v1/pizzas?page=0",
		"/api/v1/pizzas?page_size=101",
		"/api/v1/pizzas?page=nope",
	} {
		request := httptest.NewRequest("GET", target, nil)
		recorder := httptest.NewRecorder()
		if _, _, ok := pagination(recorder, request); ok {
			t.Errorf("expected %q to be rejected", target)
		}
		if recorder.Code != 400 {
			t.Errorf("expected 400 for %q, got %d", target, recorder.Code)
		}
	}
}
