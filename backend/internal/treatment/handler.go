package treatment

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/auth"
	"github.com/2410as/HairHistory/backend/internal/httpx"
)

type Handler struct {
	usecase *Usecase
}

func NewHandler(usecase *Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

type Response struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	TreatedOn string   `json:"treatedOn"`
	Services  []string `json:"services"`
	SalonName *string  `json:"salonName"`
	Memo      *string  `json:"memo"`
	Cost      *int     `json:"cost"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

func NewResponse(item Treatment) Response {
	return Response{
		ID:        item.ID.String(),
		UserID:    item.UserID.String(),
		TreatedOn: item.TreatedOn.Format(DateLayout),
		Services:  item.Services,
		SalonName: item.SalonName,
		Memo:      item.Memo,
		Cost:      item.Cost,
		CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

type request struct {
	TreatedOn string   `json:"treatedOn"`
	Services  []string `json:"services"`
	SalonName *string  `json:"salonName"`
	Memo      *string  `json:"memo"`
	Cost      *int     `json:"cost"`
}

func (r request) toInput() Input {
	return Input{
		TreatedOn: r.TreatedOn,
		Services:  r.Services,
		SalonName: r.SalonName,
		Memo:      r.Memo,
		Cost:      r.Cost,
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.usecase.List(r.Context(), auth.MustUserID(r.Context()))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	responses := make([]Response, 0, len(items))
	for _, item := range items {
		responses = append(responses, NewResponse(item))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"treatments": responses})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}

	created, err := h.usecase.Create(r.Context(), auth.MustUserID(r.Context()), req.toInput())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, NewResponse(created))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	found, err := h.usecase.Get(r.Context(), auth.MustUserID(r.Context()), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, NewResponse(found))
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	var req request
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}

	updated, err := h.usecase.Update(r.Context(), auth.MustUserID(r.Context()), id, req.toInput())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, NewResponse(updated))
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}

	if err := h.usecase.Delete(r.Context(), auth.MustUserID(r.Context()), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.NoContent(w)
}

func parseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, httpx.InvalidArgument("id must be a UUID", map[string]string{"id": "must be a UUID"})
	}
	return id, nil
}
