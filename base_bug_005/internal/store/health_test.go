package store

import (
	"path/filepath"
	"testing"
)

func TestHealthChecksDatabaseDirectory(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "nested", "db.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Health(); err != nil {
		t.Fatal(err)
	}
}
