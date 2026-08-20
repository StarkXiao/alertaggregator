package main

import (
	"alertaggregator/internal/aggregate"
	"alertaggregator/internal/api"
	"alertaggregator/internal/config"
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/migrate"
	"alertaggregator/internal/notify"
	"alertaggregator/internal/store"
	"alertaggregator/internal/worker"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := config.Default()
	flag.StringVar(&c.Address, "addr", c.Address, "address")
	flag.StringVar(&c.DatabasePath, "db", c.DatabasePath, "db")
	flag.Parse()
	s, e := store.Open(c.DatabasePath)
	if e != nil {
		log.Fatal(e)
	}
	if e = migrate.Apply(s); e != nil {
		log.Fatal(e)
	}
	m := &metrics.Metrics{}
	n := &notify.Notifier{Window: c.NotifyWindow}
	g := &aggregate.Engine{Store: s, Notifier: n, Metrics: m, ResolveAfter: c.ResolveAfter}
	srv := &http.Server{Addr: c.Address, Handler: (&api.Server{Store: s, Engine: g, Metrics: m}).Routes()}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go (&worker.Worker{Store: s, Engine: g, Interval: c.WorkerInterval}).Run(ctx)
	go func() {
		if e := srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Fatal(e)
		}
	}()
	<-ctx.Done()
	x, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	_ = srv.Shutdown(x)
}
