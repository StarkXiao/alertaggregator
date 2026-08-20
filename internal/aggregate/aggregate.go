package aggregate

import (
	"errors"
	"fmt"
	"time"

	"alertaggregator/internal/fingerprint"
	"alertaggregator/internal/id"
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/normalize"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/similarity"
	"alertaggregator/internal/store"
)

var ErrNotFound = errors.New("alert not found")
var ErrInvalidState = errors.New("invalid alert state")

type Engine struct {
	Store        *store.Store
	Notifier     *notify.Notifier
	Metrics      *metrics.Metrics
	ResolveAfter time.Duration
}

func (e *Engine) Process(ev model.Event) error {
	ev.Fingerprint = fingerprint.Build(ev)
	eventTime := ev.OccurredAt
	if eventTime.IsZero() {
		eventTime = time.Now()
	}
	processedAt := time.Now()
	created, notified, ignored, duplicate := false, false, false, false
	err := e.Store.Update(func(d *model.Snapshot) error {
		for i := range d.Events {
			if d.Events[i].ID == ev.ID && d.Events[i].Processed {
				duplicate = true
				return nil
			}
		}
		var alert *model.Alert
		for i := range d.Alerts {
			candidate := &d.Alerts[i]
			if candidate.Fingerprint == ev.Fingerprint || sameError(candidate, ev) {
				alert = candidate
				break
			}
		}
		if alert == nil {
			d.Alerts = append(d.Alerts, model.Alert{ID: id.New("alt"), Fingerprint: ev.Fingerprint, Service: ev.Service, Environment: ev.Environment, Level: ev.Level, Title: normalize.Message(ev.Message), Status: model.StatusOpen, Count: 1, FirstSeen: eventTime, LastSeen: eventTime, SampleEventIDs: []string{ev.ID}})
			alert = &d.Alerts[len(d.Alerts)-1]
			created = true
		} else {
			if alert.Status == model.StatusResolved && alert.ResolvedAt != nil && !eventTime.After(*alert.ResolvedAt) {
				ignored = true
				for i := range d.Events {
					if d.Events[i].ID == ev.ID {
						d.Events[i].Fingerprint = ev.Fingerprint
						d.Events[i].Processed = true
					}
				}
				return nil
			}
			alert.Count++
			if eventTime.Before(alert.FirstSeen) {
				alert.FirstSeen = eventTime
			}
			if eventTime.After(alert.LastSeen) {
				alert.LastSeen = eventTime
			}
			if alert.Status == model.StatusResolved && (alert.ResolvedAt == nil || eventTime.After(*alert.ResolvedAt)) {
				alert.Status = model.StatusOpen
				alert.ResolvedAt = nil
				alert.NotifyCount = 0
				alert.LastNotifiedAt = nil
				alert.NextNotifyAt = nil
				alert.AcknowledgedAt = nil
				alert.AcknowledgedBy = ""
			}
		}
		kind := ""
		if alert.Status != model.StatusResolved && alert.NotifyCount == 0 {
			kind = "triggered"
		} else if alert.Status != model.StatusResolved && alert.NextNotifyAt != nil && processedAt.After(*alert.NextNotifyAt) {
			kind = "reminder"
		}
		if kind != "" {
			e.Notifier.Record(d, alert, kind, processedAt)
			notified = true
		}
		for i := range d.Events {
			if d.Events[i].ID == ev.ID {
				d.Events[i].Fingerprint = ev.Fingerprint
				d.Events[i].Processed = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if ignored || duplicate {
		return nil
	}
	if created {
		e.Metrics.Grouped.Add(1)
	}
	if notified {
		e.Metrics.Notifications.Add(1)
	}
	return nil
}

func sameError(a *model.Alert, ev model.Event) bool {
	return a.Service == ev.Service && a.Environment == ev.Environment && a.Level == ev.Level && similarity.Similar(a.Title, normalize.Message(ev.Message))
}

func (e *Engine) Reconcile(now time.Time) error {
	notifications := 0
	err := e.Store.Update(func(d *model.Snapshot) error {
		for i := range d.Alerts {
			a := &d.Alerts[i]
			if a.Status != model.StatusResolved && now.Sub(a.LastSeen) >= e.ResolveAfter {
				a.Status = model.StatusResolved
				resolvedAt := now
				a.ResolvedAt = &resolvedAt
				e.Notifier.Record(d, a, "resolved", now)
				notifications++
			}
		}
		return nil
	})
	if err == nil && notifications > 0 {
		e.Metrics.Notifications.Add(uint64(notifications))
	}
	return err
}

func (e *Engine) Ack(idv string) error {
	return e.Store.Update(func(d *model.Snapshot) error {
		for i := range d.Alerts {
			a := &d.Alerts[i]
			if a.ID != idv {
				continue
			}
			if a.Status == model.StatusResolved {
				return ErrInvalidState
			}
			now := time.Now()
			a.Status = model.StatusAcknowledged
			a.AcknowledgedAt = &now
			a.AcknowledgedBy = "api"
			return nil
		}
		return fmt.Errorf("acknowledge alert: %v", ErrNotFound)
	})
}

func (e *Engine) Resolve(idv string) error {
	notified := false
	err := e.Store.Update(func(d *model.Snapshot) error {
		for i := range d.Alerts {
			a := &d.Alerts[i]
			if a.ID != idv {
				continue
			}
			if a.Status == model.StatusResolved {
				return nil
			}
			now := time.Now()
			a.Status = model.StatusResolved
			a.ResolvedAt = &now
			e.Notifier.Record(d, a, "resolved", now)
			notified = true
			return nil
		}
		return ErrNotFound
	})
	if err == nil && notified {
		e.Metrics.Notifications.Add(1)
	}
	return err
}
