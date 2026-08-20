package notify

import (
	"alertaggregator/internal/model"
	"testing"
	"time"
)

func TestBug010_FirstNotificationUsesInitialBackoff(t *testing.T) {
	now := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	snapshot := &model.Snapshot{}
	alert := &model.Alert{ID: "alert", Title: "failed"}
	(&Notifier{Window: time.Minute}).Record(snapshot, alert, "triggered", now)
	want := now.Add(2 * time.Minute)
	if alert.NextNotifyAt == nil || !alert.NextNotifyAt.Equal(want) {
		t.Fatalf("next notification=%v want=%v", alert.NextNotifyAt, want)
	}
}
