package user_change_password

import (
	"bytes"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
)

type htmlChangePwdTemplateData struct {
	Login    string
	Password string
}

type handler struct {
	userRepo           user.Repository
	changePwdTemplate  *template.Template
	sendToEmailDatabus sendToEmailDatabus
}

func NewHandler(userRepo user.Repository, changePwdTemplate *template.Template, sendToEmailDatabus sendToEmailDatabus) *handler {
	return &handler{
		userRepo:           userRepo,
		changePwdTemplate:  changePwdTemplate,
		sendToEmailDatabus: sendToEmailDatabus,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserChangePasswordRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Ошибка запроса, отправляемые данные некорректные", w)
			return
		}

		usr, err := h.userRepo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if usr == nil {
			api.Return404("Такого пользователя не существует", w)
			return
		}

		newPwd := newPassword(15)

		if err := h.userRepo.ChangePassword(r.Context(), usr.ID, auth.GeneratePassword(r.Context(), newPwd)); err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		var body bytes.Buffer
		err = h.changePwdTemplate.Execute(&body, htmlChangePwdTemplateData{
			Login:    req.Email,
			Password: newPwd,
		})
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		err = h.sendToEmailDatabus.PublishDatabusSendToEmail(r.Context(), databus.SendToEmailDTO{
			Email:   req.Email,
			Subject: "Регистрация на сайте Godzilla Soft",
			Body:    body.String(),
		})
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(nil, w)
	}
}

func newPassword(l int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz" + "0123456789")
	s := make([]rune, l)
	for j := 0; j < l; j++ {
		s[j] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}
