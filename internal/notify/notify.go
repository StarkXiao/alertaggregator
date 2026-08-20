package notify

import (
	"alertaggregator/internal/id"
	"alertaggregator/internal/model"
	"fmt"
	"time"
)

type Notifier struct {
	Window time.Duration
}

func (n *Notifier) Record(d *model.Snapshot, a *model.Alert, kind string, now time.Time) {
	d.Notifications = append(d.Notifications, model.Notification{ID: id.New("ntf"), AlertID: a.ID, Kind: kind, Channel: "log", Message: fmt.Sprintf("%s %s", kind, a.Title), CreatedAt: now, Delivered: true})
	a.NotifyCount++
	a.LastNotifiedAt = &now
	next := now.Add(n.Window * time.Duration(1<<min(a.NotifyCount+1, 4)))
	a.NextNotifyAt = &next
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
