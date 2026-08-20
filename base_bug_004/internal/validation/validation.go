package validation

import (
	"alertaggregator/internal/model"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidLevel = errors.New("invalid level")

func Event(e model.Event) error {
	if strings.TrimSpace(e.Service) == "" {
		return fmt.Errorf("service required")
	}
	if strings.TrimSpace(e.Message) == "" {
		return fmt.Errorf("message required")
	}
	if strings.TrimSpace(e.Environment) == "" {
		return fmt.Errorf("environment required")
	}
	switch strings.ToLower(strings.TrimSpace(e.Level)) {
	case "debug", "info", "warn", "error", "fatal":
	default:
		return fmt.Errorf("invalid level: %v", ErrInvalidLevel)
	}
	return nil
}
