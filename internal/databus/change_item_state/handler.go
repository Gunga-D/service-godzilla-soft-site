package change_item_state

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type handler struct {
	consumer databus.Consumer
	itemRepo itemRepo
}

func NewHandler(consumer databus.Consumer, itemRepo itemRepo) *handler {
	return &handler{
		consumer: consumer,
		itemRepo: itemRepo,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusChangeItemState(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus change item state: %v", err)
	}
	for msg := range msgs {
		var data databus.ChangeItemStateDTO
		json.Unmarshal(msg.Body, &data)

		log.Printf("[info] item %d change state to %s\n", data.ItemID, data.Status)

		if err := h.itemRepo.ChangeItemState(ctx, data.ItemID, data.Status); err != nil {
			log.Printf("[error] cannot change state: %v\n", err)
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
}
