package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"alertaggregator/internal/metrics"
	"alertaggregator/internal/store"
)

func TestEventsRejectsTrailingJSON(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	h := (&Server{Store: s, Metrics: &metrics.Metrics{}}).Routes()
	body := `{"service":"api","environment":"prod","level":"error","message":"failed"}{}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestEventsCannotSetInternalFields(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	h := (&Server{Store: s, Metrics: &metrics.Metrics{}}).Routes()
	body := `{"service":"api","environment":"prod","level":"error","message":"failed","Processed":true,"Fingerprint":"forged"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status=%d", rec.Code)
	}
	event := s.Events()[0]
	if event.Processed || event.Fingerprint != "" {
		t.Fatalf("internal fields accepted: %+v", event)
	}
}
