package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

const requestIDHeader = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			generated, err := uuid.NewV7()
			if err != nil {
				generated = uuid.New()
			}
			id = generated.String()
		}
		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDContextKey, id)))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDContextKey).(string)
	return id
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (s *statusRecorder) WriteHeader(status int) {
	if !s.wrote {
		s.status = status
		s.wrote = true
	}
	s.ResponseWriter.WriteHeader(status)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if !s.wrote {
		s.status = http.StatusOK
		s.wrote = true
	}
	return s.ResponseWriter.Write(b)
}

const redactedSegment = "[REDACTED]"

// Paths whose first segment after the prefix is a secret, not an identifier.
var secretSegmentPrefixes = []string{
	"/api/public/shares/",
	"/api/shares/",
}

func RedactPath(path string) string {
	for _, prefix := range secretSegmentPrefixes {
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		rest := path[len(prefix):]
		if rest == "" {
			return path
		}
		if slash := strings.IndexByte(rest, '/'); slash >= 0 {
			return prefix + redactedSegment + rest[slash:]
		}
		return prefix + redactedSegment
	}
	return path
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		slog.Info("request",
			"request_id", RequestIDFromContext(r.Context()),
			"method", r.Method,
			"path", RedactPath(r.URL.Path),
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("panic recovered",
					"request_id", RequestIDFromContext(r.Context()),
					"method", r.Method,
					"path", RedactPath(r.URL.Path),
					"panic", recovered,
				)
				JSON(w, http.StatusInternalServerError, errorEnvelope{Error: errorBody{
					Code:    CodeInternal,
					Message: "internal server error",
				}})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
