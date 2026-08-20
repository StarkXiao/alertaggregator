package store

import (
	"alertaggregator/internal/model"
	"errors"
	"path/filepath"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "db.json"))
	if e != nil {
		t.Fatal(e)
	}
	if e = s.AddEvent(model.Event{ID: "1"}); e != nil {
		t.Fatal(e)
	}
	if len(s.Events()) != 1 {
		t.Fatal("missing")
	}
}

func TestUpdateRollsBackOnCallbackError(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "db.json"))
	if err != nil {
		t.Fatal(err)
	}
	want := errors.New("stop")
	err = s.Update(func(data *model.Snapshot) error {
		data.Events = append(data.Events, model.Event{ID: "discard"})
		return want
	})
	if !errors.Is(err, want) || len(s.Events()) != 0 {
		t.Fatal("failed update changed live state")
	}
}

func TestUpdateRollsBackOnFlushError(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "db.json"))
	if err != nil {
		t.Fatal(err)
	}
	s.path = t.TempDir()
	if err := s.AddEvent(model.Event{ID: "discard"}); err == nil {
		t.Fatal("expected flush error")
	}
	if len(s.Events()) != 0 {
		t.Fatal("flush error changed live state")
	}
}
