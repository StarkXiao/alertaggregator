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

func TestProcessIsIdempotentForProcessedEvent(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	m := &metrics.Metrics{}
	e := &Engine{Store: s, Metrics: m, Notifier: &notify.Notifier{Window: time.Minute}}
	ev := model.Event{ID: "same", Service: "api", Environment: "prod", Level: "error", Message: "failed", OccurredAt: time.Now()}
	if err := s.AddEvent(ev); err != nil {
		t.Fatal(err)
	}
	if err := e.Process(ev); err != nil {
		t.Fatal(err)
	}
	first := len(s.Alerts())
	notifications := m.Notifications.Load()
	if err := e.Process(ev); err != nil {
		t.Fatal(err)
	}
	if len(s.Alerts()) != first || m.Notifications.Load() != notifications {
		t.Fatalf("duplicate changed state: alerts=%d notifications=%d", len(s.Alerts()), m.Notifications.Load())
	}
}
