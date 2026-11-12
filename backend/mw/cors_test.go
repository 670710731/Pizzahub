package mw

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORSAllowsAdminPreflightMethods(t *testing.T) {
	handler := CORS([]string{"https://pizzahub.example"}, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/admin/pizzas/1", nil)
	request.Header.Set("Origin", "https://pizzahub.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodPut)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://pizzahub.example" {
		t.Fatalf("allow-origin = %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "PUT") {
		t.Fatalf("allow-methods = %q, want PUT", got)
	}
}
