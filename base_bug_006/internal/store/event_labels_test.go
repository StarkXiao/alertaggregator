package store

import (
	"alertaggregator/internal/model"
	"path/filepath"
	"testing"
)

func TestBug006_EventsReturnsIndependentLabels(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AddEvent(model.Event{ID: "event", Labels: map[string]string{"region": "east"}}); err != nil {
		t.Fatal(err)
	}
	events := s.Events()
	events[0].Labels["region"] = "west"
	if got := s.Events()[0].Labels["region"]; got != "east" {
		t.Fatalf("stored label changed to %q", got)
	}
}
