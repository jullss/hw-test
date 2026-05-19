package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		interceptor := &statusInterceptor{rw: w, statusCode: http.StatusOK}

		next.ServeHTTP(interceptor, r)

		latency := time.Since(startTime)
		timeFormatted := startTime.Format("01/Jan/2000:00:00:00 -0700")

		ua := r.UserAgent()
		if ua == "" {
			ua = "-"
		}

		logLine := fmt.Sprintf("%s [%s] %s %s %s %d (latency: %v) %q",
			r.RemoteAddr,
			timeFormatted,
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			interceptor.statusCode,
			latency,
			ua,
		)

		s.log.Info(logLine)
	})
}

type statusInterceptor struct {
	rw         http.ResponseWriter
	statusCode int
}

func (i *statusInterceptor) WriteHeader(code int) {
	i.statusCode = code
	i.rw.WriteHeader(code)
}

func (i *statusInterceptor) Write(b []byte) (int, error) {
	return i.rw.Write(b)
}

func (i *statusInterceptor) Header() http.Header {
	return i.rw.Header()
}
