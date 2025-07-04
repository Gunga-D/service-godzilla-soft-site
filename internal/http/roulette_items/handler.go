package roulette_items

import (
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette/postgres"
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

func (h *Handler) HandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			api.Return400("Поддерживается только json формат тела запроса", w)
			return
		}

		var req CreateRouletteItemsRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		err := h.r.AddItemsToRoulette(r.Context(), req.Items)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		api.ReturnOK(nil, w)
	}
}

func (h *Handler) HandleFetchItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := h.r.FetchAvailableItems(r.Context())
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		api.ReturnOK(items, w)
	}
}

func (h *Handler) HandleRoll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handler) PaymentCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
