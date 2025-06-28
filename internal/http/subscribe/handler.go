package subscribe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

const (
	_subscriptionNotifyURL = "https://api.godzillasoft.ru/v1/payment_subscription_notification"
	_isSoon                = true
)

type htmlRegistrationTemplateData struct {
	Login    string
	Password string
}

type handler struct {
	jwtService           jwtService
	registrationTemplate *template.Template
	userRepo             user.Repository
	sendToEmailDatabus   sendToEmailDatabus
	subRepo              subRepo
	subChecker           subChecker
	tinkoffClient        tinkoff.Client
}

func NewHandler(jwtService jwtService, registrationTemplate *template.Template,
	userRepo user.Repository, sendToEmailDatabus sendToEmailDatabus,
	subRepo subRepo, subChecker subChecker, tinkoffClient tinkoff.Client) *handler {

	return &handler{
		jwtService:           jwtService,
		registrationTemplate: registrationTemplate,
		userRepo:             userRepo,
		sendToEmailDatabus:   sendToEmailDatabus,
		subRepo:              subRepo,
		subChecker:           subChecker,
		tinkoffClient:        tinkoffClient,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID *int64
		if v, ok := r.Context().Value(user.MetaUserIDKey{}).(int64); ok {
			userID = pointer.ToInt64(v)
		}

		var body SubscribeRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		if userID != nil {
			hasSub, err := h.subChecker.HasSubscription(r.Context(), *userID)
			if err != nil {
				api.Return500("Неизвестная ошибка", w)
				return
			}

			if hasSub {
				api.Return409("У вас уже есть подписка", w)
				return
			}
		}

		if userID == nil && body.Email == nil {
			api.Return401("Пользователь неавторизован", w)
			return
		}

		var accessToken *string
		if body.Email != nil && userID == nil {
			if ok := auth.ValidateEmail(*body.Email); !ok {
				api.Return400("Почта пользователя невалидная", w)
				return
			}

			newUserAccessToken, newUserID, err := h.syncRegisterUser(r.Context(), *body.Email)
			if err != nil {
				api.Return500(err.Error(), w)
				return
			}
			userID = &newUserID
			accessToken = &newUserAccessToken
		}

		if _isSoon {
			if body.Email != nil {
				logger.Get().Log(fmt.Sprintf("❤️ Баер захотел приобрести подписку, но мы его просто зарегали %s", *body.Email))
			} else {
				if userID != nil {
					logger.Get().Log(fmt.Sprintf("❤️ Зареганный пользователь %d захотел приобрести подписку", *userID))
				}
			}
			api.ReturnOK(SubscribeResponse{
				UserAccessToken: accessToken,
				RedirectLink:    "https://godzillasoft.ru/subscription_soon",
			}, w)
			return
		}

		cost, ok := prices[body.Period]
		if !ok {
			api.Return400("Период принимает значения month и year", w)
			return
		}
		duration, ok := durations[body.Period]
		if !ok {
			api.Return400("Период принимает значения month и year", w)
			return
		}
		durationName, ok := durationNames[body.Period]
		if !ok {
			api.Return400("Период принимает значения month и year", w)
			return
		}

		id, err := h.subRepo.CreateSubscriptionBill(r.Context(), *userID, cost, duration)
		if err != nil {
			log.Printf("[error] cannot create subscription bill %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		resp, err := h.tinkoffClient.CreateRecurrent(r.Context(), id, cost, fmt.Sprintf("Подписка GODZILLA SOFT на %s", durationName), fmt.Sprint(*userID), _subscriptionNotifyURL)
		if err != nil {
			log.Printf("[error] cannot create recurrent payment %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		api.ReturnOK(SubscribeResponse{
			UserAccessToken: accessToken,
			SubscriptionID:  id,
			RedirectLink:    resp.PaymentURL,
		}, w)
	}
}

func (h *handler) syncRegisterUser(ctx context.Context, email string) (string, int64, error) {
	newPwd := newPassword(15)
	userID, err := h.userRepo.CreateUser(ctx, user.User{
		Email:    pointer.ToString(email),
		Password: pointer.ToString(auth.GeneratePassword(ctx, newPwd)),
	})
	if err != nil {
		return "", 0, errors.New("Необходимо авторизоваться перед покупкой подписки")
	}

	var body bytes.Buffer
	err = h.registrationTemplate.Execute(&body, htmlRegistrationTemplateData{
		Login:    email,
		Password: newPwd,
	})
	if err != nil {
		log.Printf("[error] cannot execute registration template: %v\n", err)
		return "", 0, errors.New("Неизвестная ошибка")
	}

	err = h.sendToEmailDatabus.PublishDatabusSendToEmail(ctx, databus.SendToEmailDTO{
		Email:   email,
		Subject: "Регистрация на сайте Godzilla Soft",
		Body:    body.String(),
	})
	if err != nil {
		log.Printf("[error] cannot publish to send to email databus: %v\n", err)
		return "", 0, errors.New("Неизвестная ошибка")
	}

	accessToken, err := h.jwtService.GenerateToken(userID, &email)
	if err != nil {
		log.Printf("[error] cannot generate access token: %v\n", err)
		return "", 0, errors.New("Неизвестная ошибка")
	}
	return accessToken, userID, nil
}

func newPassword(l int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz" + "0123456789")
	s := make([]rune, l)
	for j := 0; j < l; j++ {
		s[j] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}
