package telegram_registration

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	tele "gopkg.in/telebot.v4"
)

type handler struct {
	consumer databus.Consumer
	bot      *tele.Bot
}

func NewHandler(consumer databus.Consumer, bot *tele.Bot) *handler {
	return &handler{
		consumer: consumer,
		bot:      bot,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusTelegramRegistration(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus telegram registration: %v", err)
	}
	for msg := range msgs {
		var data databus.TelegramRegistrationDTO
		json.Unmarshal(msg.Body, &data)

		menu := &tele.ReplyMarkup{ResizeKeyboard: true}
		menu.Inline(
			tele.Row{menu.Data("Проверить подписку", "checkSubscription")},
		)
		_, err := h.bot.Send(&tele.User{
			ID: data.TelegramID,
		}, "Привет\\! Это бот GODZILLA SOFT, спасибо за регистрацию\\! Для получения *БЕСПЛАТНОЙ СЛУЧАЙНОЙ STEAM ИГРЫ* осталось только подписаться на наш [канал](https://t.me/godzillasoftmedia) и нажать кнопку \"Проверить подписку\"\\.", menu, tele.ModeMarkdownV2)
		if err != nil {
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}
}
