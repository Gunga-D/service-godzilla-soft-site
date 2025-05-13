package new_user_steam_link

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type handler struct {
	consumer databus.Consumer
	userRepo userRepo
}

func NewHandler(consumer databus.Consumer, userRepo userRepo) *handler {
	return &handler{
		consumer: consumer,
		userRepo: userRepo,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusNewUserSteamLink(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus new user steam link: %v", err)
	}
	for msg := range msgs {
		var data databus.NewUserSteamLinkDTO
		json.Unmarshal(msg.Body, &data)

		log.Printf("[info] new steam link %s for user %d\n", data.SteamLink, data.UserID)

		usr, err := h.userRepo.GetUserByID(ctx, data.UserID)
		if err != nil {
			log.Printf("[error] cannot get user by id - %d: %v\n", data.UserID, err)
			msg.Nack(false, true)
			continue
		}

		if usr.SteamLink == nil || (usr.SteamLink != nil && *usr.SteamLink != data.SteamLink) {
			err = h.userRepo.AssignSteamLinkToUser(ctx, data.UserID, data.SteamLink)
			if err != nil {
				log.Printf("[error] cannot assign steam link to user - %d: %v\n", data.UserID, err)
				msg.Nack(false, true)
				continue
			}
			log.Printf("[info] assign new steam link %s to user %d\n", data.SteamLink, data.UserID)
		}

		msg.Ack(false)
	}
}
