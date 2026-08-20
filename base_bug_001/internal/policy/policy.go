package policy

import "time"

func ShouldNotify(last time.Time, now time.Time, window time.Duration) bool {
	return last.IsZero() || now.Sub(last) >= window
}
func Backoff(n int) time.Duration {
	if n > 8 {
		n = 8
	}
	return time.Duration(1<<n) * time.Minute
}
