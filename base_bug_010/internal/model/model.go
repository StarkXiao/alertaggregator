package model

import "time"

type AlertStatus string

const (
	StatusOpen         AlertStatus = "open"
	StatusAcknowledged AlertStatus = "acknowledged"
	StatusResolved     AlertStatus = "resolved"
)

type Event struct {
	ID, Service, Environment, Level, Message string
	Labels                                   map[string]string `json:"labels"`
	OccurredAt                               time.Time         `json:"occurred_at"`
	Fingerprint                              string            `json:"-"`
	Processed                                bool              `json:"-"`
}
type Alert struct {
	ID             string      `json:"id"`
	Fingerprint    string      `json:"fingerprint"`
	Service        string      `json:"service"`
	Environment    string      `json:"environment"`
	Level          string      `json:"level"`
	Title          string      `json:"title"`
	Status         AlertStatus `json:"status"`
	Count          int         `json:"count"`
	FirstSeen      time.Time   `json:"first_seen"`
	LastSeen       time.Time   `json:"last_seen"`
	ResolvedAt     *time.Time  `json:"resolved_at,omitempty"`
	NotifyCount    int         `json:"notify_count"`
	LastNotifiedAt *time.Time  `json:"last_notified_at,omitempty"`
	NextNotifyAt   *time.Time  `json:"next_notify_at,omitempty"`
	SampleEventIDs []string    `json:"sample_event_ids"`
	AcknowledgedBy string      `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time  `json:"acknowledged_at,omitempty"`
}
type Notification struct {
	ID, AlertID, Kind, Channel, Message string
	CreatedAt                           time.Time
	Delivered                           bool
}
type Snapshot struct {
	SchemaVersion int            `json:"schema_version"`
	Events        []Event        `json:"events"`
	Alerts        []Alert        `json:"alerts"`
	Notifications []Notification `json:"notifications"`
}
