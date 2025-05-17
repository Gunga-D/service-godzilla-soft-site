package think

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
	"github.com/google/uuid"
)

type handler struct {
	neuroDatabus neuroDatabus
}

func NewHandler(neuroDatabus neuroDatabus) *handler {
	return &handler{
		neuroDatabus: neuroDatabus,
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
		if err := h.neuroDatabus.PublishDatabusNeuroTask(r.Context(), databus.NeuroTaskDTO{
			ID:    taskID,
			Query: body.Query,
		}); err != nil {
			log.Printf("cannot publish to neuro task: %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		err := h.neuroDatabus.PublishDatabusNeuroNewItems(r.Context(), databus.NeuroNewItemsDTO{
			Query: body.Query,
		})
		if err != nil {
			log.Printf("cannot publish to neuro new items queue: %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		api.ReturnOK(ThinkResponse{
			ID: taskID,
		}, w)
	}
}
