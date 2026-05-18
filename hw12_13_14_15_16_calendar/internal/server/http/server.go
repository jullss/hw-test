package internalhttp

import (
	"context"
	"net/http"
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

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		log: logger,
		app: app,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

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
