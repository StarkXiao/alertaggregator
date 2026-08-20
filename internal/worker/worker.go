package worker

import (
	"alertaggregator/internal/aggregate"
	"alertaggregator/internal/store"
	"context"
	"time"
)

type Worker struct {
	Store    *store.Store
	Engine   *aggregate.Engine
	Interval time.Duration
}

func (w *Worker) Run(ctx context.Context) {
	if w.Interval <= 0 {
		w.Interval = time.Second
	}
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	w.tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.tick()
		}
	}
}
func (w *Worker) tick() {
	for _, e := range w.Store.Events() {
		if !e.Processed {
			if err := w.Engine.Process(e); err != nil {
				w.Engine.Metrics.Errors.Add(1)
			}
		}
	}
	if err := w.Engine.Reconcile(time.Now()); err != nil {
		w.Engine.Metrics.Errors.Add(1)
	}
}
