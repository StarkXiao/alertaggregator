package retention

import (
	"alertaggregator/internal/model"
	"alertaggregator/internal/store"
	"time"
)

func Purge(s *store.Store, before time.Time) error {
	return s.Update(func(d *model.Snapshot) error {
		out := d.Events[:0]
		for _, e := range d.Events {
			if !e.Processed || e.OccurredAt.After(before) {
				out = append(out, e)
			}
		}
		d.Events = out
		return nil
	})
}
