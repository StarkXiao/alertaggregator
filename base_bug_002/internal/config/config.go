package config

import "time"

type Config struct {
	Address, DatabasePath                      string
	WorkerInterval, NotifyWindow, ResolveAfter time.Duration
}

func Default() Config {
	return Config{Address: ":8080", DatabasePath: "./data/alerts.json", WorkerInterval: time.Second, NotifyWindow: 5 * time.Minute, ResolveAfter: 30 * time.Minute}
}
