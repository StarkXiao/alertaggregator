package health

import "alertaggregator/internal/store"

func Ready(s *store.Store) bool { return s != nil && s.Health() == nil }
