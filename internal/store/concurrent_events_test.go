package store

import (
	"alertaggregator/internal/model"
	"path/filepath"
	"sync"
	"testing"
)

func TestBug001_ConcurrentEventsReadIsSynchronized(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "alerts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	group.Add(2)
	go func() {
		defer group.Done()
		for i := 0; i < 200; i++ {
			if err := s.AddEvent(model.Event{ID: string(rune(i))}); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	go func() {
		defer group.Done()
		for i := 0; i < 200; i++ {
			_ = s.Events()
		}
	}()
	group.Wait()
}
