package store

import (
	"alertaggregator/internal/model"
	"path/filepath"
	"testing"
)

func TestBug003_AlertsReturnsIndependentSamples(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Update(func(snapshot *model.Snapshot) error {
		snapshot.Alerts = []model.Alert{{ID: "alert", SampleEventIDs: []string{"event-a"}}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	alerts := s.Alerts()
	alerts[0].SampleEventIDs[0] = "changed"
	if got := s.Alerts()[0].SampleEventIDs[0]; got != "event-a" {
		t.Fatalf("stored sample changed to %q", got)
	}
}
