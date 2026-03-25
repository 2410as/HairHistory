package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/annasakai/hairhistorymemo/apps/main/app/usecase/request"
)

type Handler struct {
	deps Deps
}

func NewRouter(deps Deps) http.Handler {
	h := &Handler{deps: deps}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/users", h.CreateUser)
	mux.HandleFunc("/api/users/", h.UsersHistories)
	mux.HandleFunc("/api/histories/", h.HistoriesByID)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func notFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, map[string]any{
		"error": "not found",
	})
}

func writeError(w http.ResponseWriter, status int, err error) {
	msg := "internal error"
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}
	writeJSON(w, status, map[string]any{
		"error": msg,
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// CreateUser handles POST /api/users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, err := request.NewCreateUser(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.deps.CreateUser.Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusNotImplemented, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"userId": resp.UserID,
	})
}

// UsersHistories handles:
// - GET  /api/users/{userId}/histories
// - POST /api/users/{userId}/histories
func (h *Handler) UsersHistories(w http.ResponseWriter, r *http.Request) {
	// Path example: /api/users/{userId}/histories
	p := strings.TrimPrefix(r.URL.Path, "/api/users/")
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) != 2 || parts[1] != "histories" || parts[0] == "" {
		notFound(w)
		return
	}
	userID := parts[0]

	switch r.Method {
	case http.MethodGet:
		req, err := request.NewListHistories(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		_ = userID // route validation only
		resp, err := h.deps.ListHistories.Execute(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusNotImplemented, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"histories": resp.Histories,
		})
	case http.MethodPost:
		req, err := request.NewCreateHistory(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		_ = userID // route validation only
		resp, err := h.deps.CreateHistory.Execute(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusNotImplemented, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"history": resp.History,
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HistoriesByID handles:
// - PUT    /api/histories/{historyId}
// - DELETE /api/histories/{historyId}
func (h *Handler) HistoriesByID(w http.ResponseWriter, r *http.Request) {
	// Path example: /api/histories/{historyId}
	historyID := strings.TrimPrefix(r.URL.Path, "/api/histories/")
	historyID = strings.Trim(historyID, "/")
	if historyID == "" {
		notFound(w)
		return
	}

	switch r.Method {
	case http.MethodPut:
		req, err := request.NewUpdateHistory(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		_ = historyID // route validation only
		resp, err := h.deps.UpdateHistory.Execute(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusNotImplemented, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"history": resp.History,
		})
	case http.MethodDelete:
		req, err := request.NewDeleteHistory(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		_ = historyID // route validation only
		resp, err := h.deps.DeleteHistory.Execute(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusNotImplemented, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": resp.OK})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

