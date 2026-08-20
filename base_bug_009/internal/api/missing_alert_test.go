package api

import (
	"alertaggregator/internal/aggregate"
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/store"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestBug009_MissingAlertAckReturnsNotFound(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	m := &metrics.Metrics{}
	e := &aggregate.Engine{Store: s, Notifier: &notify.Notifier{Window: time.Minute}, Metrics: m}
	h := (&Server{Store: s, Engine: e, Metrics: m}).Routes()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/alerts/missing/ack", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
