package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

const SessionCookieName = "hh_session"

type Handler struct {
	usecase      *Usecase
	cookieSecure bool
	devLogin     bool
}

func NewHandler(usecase *Usecase, cookieSecure, devLoginEnabled bool) *Handler {
	return &Handler{usecase: usecase, cookieSecure: cookieSecure, devLogin: devLoginEnabled}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/google", h.loginWithGoogle)
	r.Post("/logout", h.logout)
	r.With(RequireUser(h.usecase)).Get("/me", h.me)
	if h.devLogin {
		r.Post("/dev-login", h.devLoginHandler)
	}
	return r
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

func newUserResponse(user User) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}

type googleLoginRequest struct {
	IDToken string `json:"idToken"`
}

func (h *Handler) loginWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req googleLoginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}

	result, err := h.usecase.LoginWithGoogle(r.Context(), req.IDToken)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	h.setSessionCookie(w, result.SessionToken, result.SessionExpiresAt)
	httpx.JSON(w, http.StatusOK, newUserResponse(result.User))
}

type devLoginRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) devLoginHandler(w http.ResponseWriter, r *http.Request) {
	if !h.devLogin {
		httpx.Error(w, r, httpx.NotFound("endpoint not found"))
		return
	}

	var req devLoginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}

	result, err := h.usecase.LoginForDevelopment(r.Context(), req.Email, req.Name)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	h.setSessionCookie(w, result.SessionToken, result.SessionExpiresAt)
	httpx.JSON(w, http.StatusOK, newUserResponse(result.User))
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		httpx.Error(w, r, httpx.Unauthenticated("authentication required"))
		return
	}
	httpx.JSON(w, http.StatusOK, newUserResponse(user))
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(SessionCookieName); err == nil {
		if err := h.usecase.Logout(r.Context(), cookie.Value); err != nil {
			httpx.Error(w, r, err)
			return
		}
	}
	h.clearSessionCookie(w)
	httpx.NoContent(w)
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
