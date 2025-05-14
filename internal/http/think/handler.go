package think

import (
	"fmt"
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
	"github.com/google/uuid"
)

type handler struct {
	neuroTaskDatabus neuroTaskDatabus
}

func NewHandler(neuroTaskDatabus neuroTaskDatabus) *handler {
	return &handler{
		neuroTaskDatabus: neuroTaskDatabus,
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

		taskID := uuid.NewString()
		if err := h.neuroTaskDatabus.PublishDatabusNeuroTask(r.Context(), databus.NeuroTaskDTO{
			ID:    taskID,
			Query: body.Query,
		}); err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		api.ReturnOK(ThinkResponse{
			ID: taskID,
		}, w)
	}
}
