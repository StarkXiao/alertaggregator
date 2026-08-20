package aggregate

import (
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestLateEventAfterResolutionIsIgnored(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	m := &metrics.Metrics{}
	e := &Engine{Store: s, Metrics: m, Notifier: &notify.Notifier{Window: time.Minute}}
	resolvedAt := time.Now()
	alert := model.Alert{ID: "a1", Fingerprint: "f1", Service: "api", Environment: "prod", Level: "error", Title: "db timeout", Status: model.StatusResolved, Count: 2, FirstSeen: resolvedAt.Add(-time.Hour), LastSeen: resolvedAt.Add(-time.Minute), ResolvedAt: &resolvedAt}
	if err := s.Update(func(d *model.Snapshot) error {
		d.Alerts = []model.Alert{alert}
		d.Events = []model.Event{{ID: "late", Service: "api", Environment: "prod", Level: "error", Message: "db timeout 1", OccurredAt: resolvedAt.Add(-time.Minute)}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := e.Process(s.Events()[0]); err != nil {
		t.Fatal(err)
	}
	got := s.Alerts()[0]
	if got.Count != 2 || got.Status != model.StatusResolved || len(s.Events()) != 1 || !s.Events()[0].Processed {
		t.Fatalf("late event changed alert: %+v", got)
	}
}
