package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// we keep original repose writer  and remember the statusCode

type reposeWriter struct {
	http.ResponseWriter
	statusCode int
}

// This function runs whenever someone sets a status code

func (rw *reposeWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// This is the middleware function
// It wraps the real handler and adds logging

func loggigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// We take the real writer w and put it inside our own struct.
		ww := &reposeWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// this what show what the method and url and what the ip and which browser it using
		slog.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_add", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// this is acutally execute and when call the path
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// this will give same path and status code and time it took
		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
