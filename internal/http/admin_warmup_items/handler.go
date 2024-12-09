package admin_warmup_items

import (
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
	cache itemsCache
}

func NewHandler(cache itemsCache) *handler {
	return &handler{
		cache: cache,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.cache.WarmUp(r.Context())
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}

		api.ReturnOK(nil, w)
	}
}
