package metrics

import (
	"strconv"
	"sync/atomic"
)

type Metrics struct{ Events, Grouped, Notifications, Errors atomic.Uint64 }

func (m *Metrics) Text() string {
	return "aggregator_events_total " + strconv.FormatUint(m.Events.Load(), 10) + "\naggregator_grouped_total " + strconv.FormatUint(m.Grouped.Load(), 10) + "\naggregator_notifications_total " + strconv.FormatUint(m.Notifications.Load(), 10) + "\naggregator_errors_total " + strconv.FormatUint(m.Errors.Load(), 10) + "\n"
}
