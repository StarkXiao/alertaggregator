package store

import (
	"alertaggregator/internal/model"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	path string
	data model.Snapshot
}
func Open(p string) (*Store, error) {
	s := &Store{path: p}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		s.data.SchemaVersion = 1
		return s, s.flush()
	}
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &s.data); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Store) flush() error {
	return s.flushData(s.data)
}
func (s *Store) flushData(data model.Snapshot) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
func (s *Store) AddEvent(e model.Event) error {
	return s.Update(func(data *model.Snapshot) error { data.Events = append(data.Events, e); return nil })
}
func (s *Store) Events() []model.Event {
	out := append([]model.Event(nil), s.data.Events...)
	for i := range out {
		out[i].Labels = copyLabels(out[i].Labels)
	}
	return out
}
func (s *Store) Alerts() []model.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := append([]model.Alert(nil), s.data.Alerts...)
	for i := range out {
		out[i].SampleEventIDs = append([]string(nil), out[i].SampleEventIDs...)
	}
	return out
}
func copyLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}
	out := make(map[string]string, len(labels))
	for key, value := range labels {
		out[key] = value
	}
	return out
}
func (s *Store) Update(fn func(*model.Snapshot) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	var next model.Snapshot
	if err := json.Unmarshal(b, &next); err != nil {
		return err
	}
	if err := fn(&next); err != nil {
		return err
	}
	if err := s.flushData(next); err != nil {
		return err
	}
	s.data = next
	return nil
}
func (s *Store) Health() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal(b, &snapshot); err != nil {
		return err
	}
	f, err := os.OpenFile(s.path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	probe, err := os.CreateTemp(filepath.Dir(s.path), ".alertaggregator-health-*")
	if err != nil {
		return err
	}
	probePath := probe.Name()
	if err := probe.Close(); err != nil {
		_ = os.Remove(probePath)
		return err
	}
	return os.Remove(probePath)
}
