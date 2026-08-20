package validation

import (
	"alertaggregator/internal/model"
	"errors"
	"testing"
)

func TestBug004_InvalidLevelPreservesSentinel(t *testing.T) {
	err := Event(model.Event{Service: "api", Environment: "prod", Level: "verbose", Message: "failed"})
	if !errors.Is(err, ErrInvalidLevel) {
		t.Fatalf("invalid level error does not preserve sentinel: %v", err)
	}
}
