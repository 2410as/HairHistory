package logging

import (
	"context"
	"log"
)

func Warningf(_ context.Context, format string, args ...any) {
	log.Printf("WARN: "+format, args...)
}

func Errorf(_ context.Context, format string, args ...any) {
	log.Printf("ERROR: "+format, args...)
}

