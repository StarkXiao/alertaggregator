package api

import (
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/store"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestBug002_EventWithoutLabelsIsAccepted(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	h := (&Server{Store: s, Metrics: &metrics.Metrics{}}).Routes()
	req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(`{"service":"api","environment":"prod","level":"error","message":"failed"}`))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status=%d", rec.Code)
	}
}
