package httpx

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const sampleToken = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestRedactPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "public share token",
			path: "/api/public/shares/" + sampleToken,
			want: "/api/public/shares/[REDACTED]",
		},
		{
			name: "share revoke token",
			path: "/api/shares/" + sampleToken,
			want: "/api/shares/[REDACTED]",
		},
		{
			name: "trailing segment after token",
			path: "/api/shares/" + sampleToken + "/extra",
			want: "/api/shares/[REDACTED]/extra",
		},
		{
			name: "share collection without token",
			path: "/api/shares/",
			want: "/api/shares/",
		},
		{
			name: "treatment uuid is not secret",
			path: "/api/treatments/0b2a9d1e-7c33-4a9f-bb0f-1f2e3d4c5b6a",
			want: "/api/treatments/0b2a9d1e-7c33-4a9f-bb0f-1f2e3d4c5b6a",
		},
		{
			name: "unrelated path",
			path: "/api/auth/me",
			want: "/api/auth/me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactPath(tt.path); got != tt.want {
				t.Errorf("RedactPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestLoggingRedactsShareToken(t *testing.T) {
	paths := []string{
		"/api/public/shares/" + sampleToken,
		"/api/shares/" + sampleToken,
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			var logged bytes.Buffer
			restore := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&logged, nil)))
			defer slog.SetDefault(restore)

			handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))

			output := logged.String()
			if strings.Contains(output, sampleToken) {
				t.Fatalf("log output leaked the token: %s", output)
			}
			if !strings.Contains(output, redactedSegment) {
				t.Errorf("log output is missing %q: %s", redactedSegment, output)
			}
		})
	}
}

func TestRecovererRedactsShareToken(t *testing.T) {
	var logged bytes.Buffer
	restore := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logged, nil)))
	defer slog.SetDefault(restore)

	handler := Recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))
	path := "/api/public/shares/" + sampleToken
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))

	if output := logged.String(); strings.Contains(output, sampleToken) {
		t.Fatalf("log output leaked the token: %s", output)
	}
}
