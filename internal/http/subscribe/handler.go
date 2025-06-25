package subscribe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
)

const (
	_subscriptionNotifyURL = "https://webhook.site/b04ace61-fb35-45b0-9fbc-cac1eb58bf15"
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
			api.Return500("Неизвестная ошибка", w)
			return
		}

		resp, err := h.tinkoffClient.CreateRecurrent(r.Context(), id, cost, fmt.Sprintf("Подписка GODZILLA SOFT на %s", durationName), fmt.Sprint(*userID), _subscriptionNotifyURL)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		api.ReturnOK(SubscribeResponse{
			UserAccessToken: accessToken,
			SubscriptionID:  id,
			PaymentLink:     resp.PaymentURL,
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
		return "", 0, errors.New("Невозможно создать пользователя, попробуйте чуть позже")
	}

	var body bytes.Buffer
	err = h.registrationTemplate.Execute(&body, htmlRegistrationTemplateData{
		Login:    email,
		Password: newPwd,
	})
	if err != nil {
		return "", 0, errors.New("Неизвестная ошибка")
	}

	err = h.sendToEmailDatabus.PublishDatabusSendToEmail(ctx, databus.SendToEmailDTO{
		Email:   email,
		Subject: "Регистрация на сайте Godzilla Soft",
		Body:    body.String(),
	})
	if err != nil {
		return "", 0, errors.New("Неизвестная ошибка")
	}

	accessToken, err := h.jwtService.GenerateToken(userID, &email)
	if err != nil {
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
