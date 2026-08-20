package aggregate

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/store"
)

func TestProcessPersistsFingerprint(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	m := &metrics.Metrics{}
	e := &Engine{Store: s, Metrics: m, ResolveAfter: time.Hour, Notifier: &notify.Notifier{Window: time.Minute}}
	event := model.Event{ID: "evt_1", Service: "payments", Environment: "prod", Level: "error", Message: "db timeout 42", OccurredAt: time.Now()}
	if err := s.AddEvent(event); err != nil {
		t.Fatal(err)
	}
	if err := e.Process(event); err != nil {
		t.Fatal(err)
	}
	stored := s.Events()[0]
	if stored.Fingerprint == "" || !stored.Processed {
		t.Fatalf("event not fully processed: %+v", stored)
	}
}

func TestOlderEventDoesNotRegressLastSeen(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	m := &metrics.Metrics{}
	e := &Engine{Store: s, Metrics: m, ResolveAfter: time.Hour, Notifier: &notify.Notifier{Window: time.Minute}}
	recent := model.Event{ID: "new", Service: "api", Environment: "prod", Level: "error", Message: "timeout 1", OccurredAt: time.Now()}
	_ = s.AddEvent(recent)
	if err := e.Process(recent); err != nil {
		t.Fatal(err)
	}
	old := recent
	old.ID = "old"
	old.OccurredAt = recent.OccurredAt.Add(-time.Hour)
	_ = s.AddEvent(old)
	if err := e.Process(old); err != nil {
		t.Fatal(err)
	}
	if !s.Alerts()[0].LastSeen.Equal(recent.OccurredAt) {
		t.Fatal("last seen regressed")
	}
}

func TestResolvedAlertCannotBeAcknowledged(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "alerts.json"))
	m := &metrics.Metrics{}
	e := &Engine{Store: s, Metrics: m, Notifier: &notify.Notifier{Window: time.Minute}}
	event := model.Event{ID: "e", Service: "api", Environment: "prod", Level: "error", Message: "failed", OccurredAt: time.Now()}
	_ = s.AddEvent(event)
	_ = e.Process(event)
	id := s.Alerts()[0].ID
	_ = e.Resolve(id)
	if !errors.Is(e.Ack(id), ErrInvalidState) {
		t.Fatal("resolved alert was acknowledged")
	}
}
