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
		go h.cache.WarmUp(r.Context())

		api.ReturnOK(nil, w)
	}
}
