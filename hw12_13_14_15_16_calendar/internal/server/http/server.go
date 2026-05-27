package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jullss/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	log    Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		log: logger,
		app: app,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create", s.createEventHandler)
	mux.HandleFunc("PATCH /update", s.updateEventHandler)
	mux.HandleFunc("DELETE /destroy", s.deleteEventHandler)
	mux.HandleFunc("GET /list_day", s.listDayEventHandler)
	mux.HandleFunc("GET /list_week", s.listWeekEventHandler)
	mux.HandleFunc("GET /list_month", s.listMonthEventHandler)

	handlerWithLogging := s.loggingMiddleware(mux)

	s.server = &http.Server{
		Addr:    addr,
		Handler: handlerWithLogging,
	}

	s.log.Info("http server starting", "addr", addr)

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.log.Info("http server stopping")
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Error("failed to decode request body", "err", err)
		http.Error(w, `{"error": "invalid json format"}`, http.StatusBadRequest)
		return
	}

	const timeLayout = "2006-01-02 15:04:05"

	startTime, err := time.Parse(timeLayout, req.StartTime)
	if err != nil {
		http.Error(w, `{"error": "invalid start_time format, use YYYY-MM-DD HH:MM:SS"}`, http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(timeLayout, req.EndTime)
	if err != nil {
		http.Error(w, `{"error": "invalid end_time format, use YYYY-MM-DD HH:MM:SS"}`, http.StatusBadRequest)
		return
	}

	notifyIn, err := time.ParseDuration(req.NotifyIn)
	if err != nil {
		http.Error(w, `{"error": "invalid notify_in format, use units like 1h, 30m"}`, http.StatusBadRequest)
		return
	}

	event := storage.Event{
		UserID:    req.UserID,
		Title:     req.Title,
		Desc:      req.Desc,
		StartTime: startTime,
		EndTime:   endTime,
		NotifyIn:  notifyIn,
	}

	if err := s.app.CreateEvent(r.Context(), &event); err != nil {
		s.log.Error("failed to create event in core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"status": "created"}`))
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Error("failed to decode request body", "err", err)
		http.Error(w, `{"error": "invalid json format"}`, http.StatusBadRequest)
		return
	}

	const timeLayout = "2006-01-02 15:04:05"

	startTime, err := time.Parse(timeLayout, req.StartTime)
	if err != nil {
		http.Error(w, `{"error": "invalid start_time format, use YYYY-MM-DD HH:MM:SS"}`, http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(timeLayout, req.EndTime)
	if err != nil {
		http.Error(w, `{"error": "invalid end_time format, use YYYY-MM-DD HH:MM:SS"}`, http.StatusBadRequest)
		return
	}

	notifyIn, err := time.ParseDuration(req.NotifyIn)
	if err != nil {
		http.Error(w, `{"error": "invalid notify_in format, use units like 1h, 30m"}`, http.StatusBadRequest)
		return
	}

	event := storage.Event{
		ID:        req.ID,
		UserID:    req.UserID,
		Title:     req.Title,
		Desc:      req.Desc,
		StartTime: startTime,
		EndTime:   endTime,
		NotifyIn:  notifyIn,
	}

	if err := s.app.UpdateEvent(r.Context(), &event); err != nil {
		s.log.Error("failed to update event in core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "updated"}`))
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Error("failed to decode request body", "err", err)
		http.Error(w, `{"error": "invalid json format"}`, http.StatusBadRequest)
		return
	}

	if err := s.app.DeleteEvent(r.Context(), req.ID); err != nil {
		s.log.Error("failed to delete event in core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "deleted"}`))
}

func (s *Server) listDayEventHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, `{"error": "missing date parameter"}`, http.StatusBadRequest)
		return
	}

	const dateLayout = "2006-01-02"

	parsedDate, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, use YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	events, err := s.app.ListDayEvent(r.Context(), parsedDate)
	if err != nil {
		s.log.Error("failed to get day events from core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		s.log.Error("failed to encode events to json", "err", err)
	}
}

func (s *Server) listWeekEventHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, `{"error": "missing date parameter"}`, http.StatusBadRequest)
		return
	}

	const dateLayout = "2006-01-02"

	parsedDate, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, use YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	events, err := s.app.ListWeekEvent(r.Context(), parsedDate)
	if err != nil {
		s.log.Error("failed to get week events from core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		s.log.Error("failed to encode events to json", "err", err)
	}
}

func (s *Server) listMonthEventHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, `{"error": "missing date parameter"}`, http.StatusBadRequest)
		return
	}

	const dateLayout = "2006-01-02"

	parsedDate, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, use YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	events, err := s.app.ListMonthEvent(r.Context(), parsedDate)
	if err != nil {
		s.log.Error("failed to get month events from core", "err", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		s.log.Error("failed to encode events to json", "err", err)
	}
}
