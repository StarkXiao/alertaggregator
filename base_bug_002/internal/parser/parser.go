package parser

import (
	"alertaggregator/internal/model"
	"alertaggregator/internal/validation"
	"encoding/json"
	"io"
)

func Decode(r io.Reader) (model.Event, error) {
	var e model.Event
	if err := json.NewDecoder(r).Decode(&e); err != nil {
		return e, err
	}
	return e, validation.Event(e)
}
