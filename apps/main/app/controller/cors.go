package controller

import (
	"log"
	"os"
	"strings"
)

// defaultCORSOrigin is used when HAIR_CORS_ORIGINS is unset or parses to no origins (local Next.js).
const defaultCORSOrigin = "http://localhost:3000"

// ResolveCORSOrigins returns AllowedOrigins for go-chi/cors. On implicit fallback to the
// development default, it logs once at call time; behavior matches the previous main-only helper.
func ResolveCORSOrigins() []string {
	origins, warn := parseCORSOrigins()
	if warn != "" {
		log.Printf("%s: %v", warn, origins)
	}
	return origins
}

func parseCORSOrigins() (origins []string, warn string) {
	raw := strings.TrimSpace(os.Getenv("HAIR_CORS_ORIGINS"))
	if raw == "" {
		return []string{defaultCORSOrigin},
			"HAIR_CORS_ORIGINS is not set; using development default allowed origin(s) only; set HAIR_CORS_ORIGINS in production"
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return []string{defaultCORSOrigin},
			"HAIR_CORS_ORIGINS has no valid origins after parsing; using development default allowed origin(s) only; set HAIR_CORS_ORIGINS in production"
	}
	return out, ""
}
