package worker

import (
	"alertaggregator/internal/aggregate"
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/store"
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestBug005_CancelledWorkerDoesNotProcessEvents(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AddEvent(model.Event{ID: "event", Service: "api", Environment: "prod", Level: "error", Message: "failed", OccurredAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	engine := &aggregate.Engine{Store: s, Notifier: &notify.Notifier{Window: time.Minute}, Metrics: &metrics.Metrics{}, ResolveAfter: time.Hour}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	(&Worker{Store: s, Engine: engine}).Run(ctx)
	if len(s.Alerts()) != 0 {
		t.Fatal("cancelled worker created an alert")
	}
}
