package render

import (
	"encoding/json"
	"log"
	"net/http"
)

const publicInternalServerErrorMessage = "サーバーエラーが発生しました。しばらくしてから再度お試しください。"

func ErrorJSON(w http.ResponseWriter, message string, status int) {
	if status >= http.StatusInternalServerError {
		log.Printf("internal error response status=%d message=%q", status, message)
		message = publicInternalServerErrorMessage
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": message,
	})
}
