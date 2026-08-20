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

func TestBug008_EventIDIsAssignedByServer(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	h := (&Server{Store: s, Metrics: &metrics.Metrics{}}).Routes()
	rec := httptest.NewRecorder()
	body := `{"ID":"caller-id","service":"api","environment":"prod","level":"error","message":"failed"}`
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(body)))
	if s.Events()[0].ID == "caller-id" {
		t.Fatal("caller controlled stored event ID")
	}
}
