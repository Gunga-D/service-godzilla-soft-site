package neuro_task

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/neuro"
)

type handler struct {
	neuroSearch neuroSearch
	neuroCache  neuroCache
	neuroRepo   neuro.Repository
	consumer    databus.Consumer
}

func NewHandler(consumer databus.Consumer, neuroSearch neuroSearch, neuroCache neuroCache, neuroRepo neuro.Repository) *handler {
	return &handler{
		consumer:    consumer,
		neuroSearch: neuroSearch,
		neuroCache:  neuroCache,
		neuroRepo:   neuroRepo,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusNeuroTask(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus neuro task: %v", err)
	}
	for msg := range msgs {
		msg := msg
		go func() {
			var data databus.NeuroTaskDTO
			json.Unmarshal(msg.Body, &data)

			log.Printf("[info] neurotask with id - %s\n", data.ID)

			res := h.neuroSearch.Search(ctx, data.ID, data.Query)
			if err := h.neuroCache.SetTaskResult(ctx, data.ID, res); err != nil {
				msg.Nack(false, true)
				return
			}

			if res.Data != nil {
				h.neuroRepo.CreateFinishedNeuroTask(ctx, neuro.Task{
					ID:     data.ID,
					Query:  data.Query,
					Result: res.Data.Raw,
				})
			}

			msg.Ack(false)
		}()
	}
}
