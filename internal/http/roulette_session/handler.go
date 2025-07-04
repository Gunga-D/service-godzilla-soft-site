package roulette_session

import (
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette/postgres"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	r *postgres.Repo
}

func NewHandler(repo *postgres.Repo) *Handler {
	return &Handler{
		r: repo,
	}
}

func (h *Handler) HandleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := h.r.CreateSession(r.Context())
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		api.ReturnOK(sessionId, w)
	}
}

func (h *Handler) HandleFetchItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Context().Value("uuid").(uuid.UUID)
		items, err := h.r.FetchItems(r.Context(), sessionId)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		var result = make([]roulette.SessionItemDTO, len(items))
		for indx, i := range items {
			result[indx] = i.DTO()
		}
		api.ReturnOK(result, w)
	}
}

func (h *Handler) HandleAddTops() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			api.Return400("Поддерживается только json формат тела запроса", w)
			return
		}

		var body AddTopItemsRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		sessionId := r.Context().Value("uuid").(uuid.UUID)
		err := h.r.AddTopItemsToSession(r.Context(), sessionId, body.ItemIds)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		api.ReturnOK(nil, w)
	}
}

func (h *Handler) HandleGetSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Context().Value("uuid").(uuid.UUID)
		data, err := h.r.GetSessionById(r.Context(), sessionId)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		api.ReturnOK(data.DTO(), w)
	}
}

func (h *Handler) HandleFormSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
