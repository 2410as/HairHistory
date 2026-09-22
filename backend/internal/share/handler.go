package share

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/2410as/HairHistory/backend/internal/auth"
	"github.com/2410as/HairHistory/backend/internal/httpx"
	"github.com/2410as/HairHistory/backend/internal/treatment"
)

type Handler struct {
	usecase *Usecase
}

func NewHandler(usecase *Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Delete("/{token}", h.revoke)
	return r
}

func (h *Handler) PublicRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{token}", h.publicView)
	return r
}

type linkResponse struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
}

func newLinkResponse(link Link) linkResponse {
	return linkResponse{
		ID:        link.ID.String(),
		Token:     link.Token,
		ExpiresAt: link.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt: link.CreatedAt.UTC().Format(time.RFC3339),
	}
}

type createRequest struct {
	ExpiresInHours *int `json:"expiresInHours"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if r.ContentLength > 0 {
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.Error(w, r, err)
			return
		}
	}

	created, err := h.usecase.Create(r.Context(), auth.MustUserID(r.Context()), req.ExpiresInHours)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, newLinkResponse(created))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	links, err := h.usecase.List(r.Context(), auth.MustUserID(r.Context()))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	responses := make([]linkResponse, 0, len(links))
	for _, link := range links {
		responses = append(responses, newLinkResponse(link))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"shares": responses})
}

func (h *Handler) revoke(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.Revoke(r.Context(), auth.MustUserID(r.Context()), chi.URLParam(r, "token")); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.NoContent(w)
}

type publicTreatmentResponse struct {
	TreatedOn string   `json:"treatedOn"`
	Services  []string `json:"services"`
	SalonName *string  `json:"salonName"`
	Memo      *string  `json:"memo"`
	Cost      *int     `json:"cost"`
}

func (h *Handler) publicView(w http.ResponseWriter, r *http.Request) {
	treatments, err := h.usecase.PublicView(r.Context(), chi.URLParam(r, "token"))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	responses := make([]publicTreatmentResponse, 0, len(treatments))
	for _, item := range treatments {
		responses = append(responses, publicTreatmentResponse{
			TreatedOn: item.TreatedOn.Format(treatment.DateLayout),
			Services:  item.Services,
			SalonName: item.SalonName,
			Memo:      item.Memo,
			Cost:      item.Cost,
		})
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"treatments": responses})
}
