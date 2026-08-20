package api

import (
	"alertaggregator/internal/aggregate"
	"alertaggregator/internal/id"
	"alertaggregator/internal/metrics"
	"alertaggregator/internal/model"
	"alertaggregator/internal/store"
	"alertaggregator/internal/validation"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Store   *store.Store
	Engine  *aggregate.Engine
	Metrics *metrics.Metrics
}

func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.health)
	m.HandleFunc("/metrics", s.metric)
	m.HandleFunc("/v1/events", s.events)
	m.HandleFunc("/v1/alerts", s.alerts)
	m.HandleFunc("/v1/alerts/", s.action)
	return m
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if s.Store.Health() != nil {
		http.Error(w, "unhealthy", 500)
		return
	}
	write(w, map[string]string{"status": "ok"})
}
func (s *Server) metric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte(s.Metrics.Text()))
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method", 405)
		return
	}
	var e model.Event
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&e) != nil || decoder.Decode(&struct{}{}) != io.EOF || validation.Event(e) != nil {
		http.Error(w, "invalid", 400)
		return
	}
	e.ID = id.New("evt")
	e.Fingerprint = ""
	e.Processed = false
	e.Service = strings.TrimSpace(e.Service)
	e.Environment = strings.TrimSpace(e.Environment)
	e.Level = strings.ToLower(strings.TrimSpace(e.Level))
	e.Message = strings.TrimSpace(e.Message)
	e.Labels["source"] = "api"
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	if err := s.Store.AddEvent(e); err != nil {
		s.Metrics.Errors.Add(1)
		http.Error(w, "store", 500)
		return
	}
	s.Metrics.Events.Add(1)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	write(w, map[string]string{"id": e.ID, "status": "accepted"})
}
func (s *Server) alerts(w http.ResponseWriter, r *http.Request) {
	all := s.Store.Alerts()
	status := r.URL.Query().Get("status")
	service := r.URL.Query().Get("service")
	environment := r.URL.Query().Get("environment")
	level := r.URL.Query().Get("level")
	limit, err := queryInt(r, "limit", 100)
	if err != nil {
		http.Error(w, "invalid limit", 400)
		return
	}
	offset, err := queryInt(r, "offset", 0)
	if err != nil {
		http.Error(w, "invalid offset", 400)
		return
	}
	out := make([]model.Alert, 0, len(all))
	for _, a := range all {
		if (status == "" || string(a.Status) == status) && (service == "" || a.Service == service) && (environment == "" || a.Environment == environment) && (level == "" || a.Level == level) {
			out = append(out, a)
		}
	}
	if offset > len(out) {
		offset = len(out)
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	write(w, out[offset:end])
}

func queryInt(r *http.Request, key string, fallback int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 || value > 1000 {
		return 0, errors.New("invalid query integer")
	}
	return value, nil
}
func (s *Server) action(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) != 4 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var err error
	if p[3] == "ack" {
		err = s.Engine.Ack(p[2])
	} else if p[3] == "resolve" {
		err = s.Engine.Resolve(p[2])
	} else {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		if errors.Is(err, aggregate.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, aggregate.ErrInvalidState) {
			http.Error(w, "invalid alert state", http.StatusConflict)
			return
		}
		s.Metrics.Errors.Add(1)
		http.Error(w, "update", 500)
		return
	}
	write(w, map[string]string{"status": "updated"})
}
