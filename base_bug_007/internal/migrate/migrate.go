package migrate

import (
	"alertaggregator/internal/model"
	"alertaggregator/internal/store"
	"fmt"
)

func Apply(s *store.Store) error {
	return s.Update(func(d *model.Snapshot) error {
		if d.SchemaVersion > 2 {
			return fmt.Errorf("unsupported schema version %d", d.SchemaVersion)
		}
		if d.SchemaVersion < 2 {
			d.SchemaVersion = 2
		}
		return nil
	})
}
