package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/telegram_registration"
	order_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/order/postgres"
	user_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/user/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)
	redis := redis.Get(ctx)
	databusClient := databus.NewClient(ctx)

	userRepo := user_postgres.NewRepo(postgres, redis)
	orderRepo := order_postgres.NewRepo(postgres)

	telebot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_LOGIN_WIDGET_BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Printf("[error] cant init telegram bot: %v", err)
		return
	}

	states := sync.Map{}
	telebot.Handle(&tele.InlineButton{Unique: "checkSubscription"}, func(ctx tele.Context) error {
		usr, err := userRepo.GetUserByTelegramID(context.Background(), ctx.Sender().ID)
		if err != nil {
			return ctx.Send("Неизвестная ошибка. Обратись, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
		}
		if usr == nil {
			return ctx.Send("Данный пользователь не зарегистрирован на сайте https://godzillasoft.ru")
		}
		if !usr.HasRegistrationGift {
			return ctx.Send("Ты уже использовал данный подарок, но можешь приобрести случайную игру [здесь](https://godzillasoft.ru/random)\\.", tele.ModeMarkdownV2)
		}

		chatmember, err := telebot.ChatMemberOf(&tele.Chat{
			ID: -1002697382470,
		}, ctx.Sender())
		if err != nil {
			log.Printf("[error] cannot get info about chat members: %v\n", err)
			return ctx.Send("Произошла ошибка получения информации о подписке. Обратитесь, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
		}
		if chatmember.Role == tele.Administrator || chatmember.Role == tele.Creator || chatmember.Role == tele.Member {
			states.Store(ctx.Sender().ID, "email_state")
			return ctx.Send("Супер\\! Давай теперь перейдем к тому, как получить *БЕСПЛАТНУЮ СЛУЧАЙНУЮ STEAM ИГРУ*\\. Отправь мне свою электронную почту для получения своего подарка:", tele.ModeMarkdownV2)
		}

		menu := &tele.ReplyMarkup{ResizeKeyboard: true}
		menu.Inline(
			tele.Row{menu.Data("Проверить подписку", "checkSubscription")},
		)
		return ctx.Send("Ты еще не подписан\\! Для получения *БЕСПЛАТНОЙ СЛУЧАЙНОЙ STEAM ИГРЫ* необходимо подписаться на наш телеграмм [канал](https://t.me/godzillasoftmedia)\\. Подпишись и еще раз нажми кнопку \"Проверить подписку\"\\.", menu, tele.ModeMarkdownV2)
	})
	telebot.Handle(tele.OnText, func(ctx tele.Context) error {
		state, ok := states.Load(ctx.Sender().ID)
		if ok && state.(string) == "email_state" {
			email := ctx.Message().Text
			_, err := mail.ParseAddress(email)
			if err != nil {
				return ctx.Send(fmt.Sprintf("Почта *%s* введена некорректно, давай попробуем еще раз\\. Отправь мне свою электронную почту для получения своего подарка:", ctx.Message().Text), tele.ModeMarkdownV2)
			}

			usr, err := userRepo.GetUserByTelegramID(context.Background(), ctx.Sender().ID)
			if err != nil {
				return ctx.Send("Неизвестная ошибка. Обратись, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
			}
			if usr == nil {
				return ctx.Send("Данный пользователь не зарегистрирован на сайте https://godzillasoft.ru")
			}
			if !usr.HasRegistrationGift {
				return ctx.Send("Ты уже использовал данный подарок, но можешь приобрести случайную игру [здесь](https://godzillasoft.ru/random)\\.", tele.ModeMarkdownV2)
			}

			err = userRepo.RemoveFreeGift(context.Background(), ctx.Sender().ID)
			if err != nil {
				log.Printf("[error] remove free gift: %v\n", err)
				return ctx.Send("Неизвестная ошибка. Обратись, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
			}

			orderID, err := orderRepo.CreateItemOrder(context.Background(), email, 0, 42, `<ol class='BxItemInstruction'>
   <li>После оплаты вам на указанную почту придет ключ активации</li>
   <li>Данный код активации необходимо ввести в клиентское приложение</li>
   <li>После активации, игра окажется в вашей библиотеке, можно играть</li>
</ol>`)
			if err != nil {
				log.Printf("[error] create item order: %v\n", err)
				return ctx.Send("Неизвестная ошибка. Обратись, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
			}

			err = orderRepo.PaidOrder(context.Background(), orderID)
			if err != nil {
				log.Printf("[error] paid order: %v\n", err)
				return ctx.Send("Неизвестная ошибка. Обратись, пожалуйста, в наш чат поддержки @GODZILLASOFT_bot")
			}

			states.Delete(ctx.Sender().ID)

			teleEmail := normalizeForTeleMarkup(email)
			return ctx.Send(fmt.Sprintf("Еще раз спасибо за регистрацию\\) Подарок придет на почту *[%s](%s)* в течение 5\\-10 минут\\. В следующий раз случайную игру ты можешь забрать [здесь](https://godzillasoft.ru/random)\\.", teleEmail, teleEmail), tele.ModeMarkdownV2)
		}
		return nil
	})

	log.Println("start chatbot")
	go telebot.Start()

	log.Println("start consume telegram registration")
	go telegram_registration.NewHandler(databusClient, telebot).Consume(ctx)

	<-ctx.Done()
}

func normalizeForTeleMarkup(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "_", "\\_"), ")", "\\)")
}
