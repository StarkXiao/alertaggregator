package api

import (
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/store"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestBug007_AlertsLimitZeroReturnsNoItems(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Update(func(snapshot *model.Snapshot) error {
		snapshot.Alerts = []model.Alert{{ID: "one"}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	h := (&Server{Store: s, Metrics: &metrics.Metrics{}}).Routes()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/alerts?limit=0", nil))
	if rec.Body.String() != "[]\n" {
		t.Fatalf("body=%q", rec.Body.String())
	}
}
