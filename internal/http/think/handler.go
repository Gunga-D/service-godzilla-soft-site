package think

import (
	"fmt"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	thinker thinker
}

func NewHandler(thinker thinker) *handler {
	return &handler{
		thinker: thinker,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ThinkRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		logger.Get().Log(fmt.Sprintf("🤔 Размышляем по запросу - %s", body.Query))

		res := h.thinker.StartThinking(r.Context(), body.Query)
		api.ReturnOK(ThinkResponse{
			ID: res,
		}, w)
	}
}
