package new_user_email

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"text/template"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
)

type htmlRegistrationTemplateData struct {
	Login    string
	Password string
}

type handler struct {
	consumer             databus.Consumer
	userRepo             userRepo
	registrationTemplate *template.Template
	sendToEmailDatabus   sendToEmailDatabus
}

func NewHandler(consumer databus.Consumer, userRepo userRepo,
	registrationTemplate *template.Template, sendToEmailDatabus sendToEmailDatabus) *handler {

	return &handler{
		consumer:             consumer,
		userRepo:             userRepo,
		registrationTemplate: registrationTemplate,
		sendToEmailDatabus:   sendToEmailDatabus,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusNewUserEmail(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus new user email: %v", err)
	}
	for msg := range msgs {
		var data databus.NewUserEmailDTO
		json.Unmarshal(msg.Body, &data)

		if data.UserID != nil {
			log.Printf("[info] new user email %s for %d\n", data.Email, *data.UserID)
		} else {
			log.Printf("[info] new user email %s without userid\n", data.Email)
		}

		usr, err := h.userRepo.GetUserByEmail(ctx, data.Email)
		if err != nil {
			log.Printf("[error] cannot get user by email - %s: %v\n", data.Email, err)
			msg.Nack(false, true)
			continue
		}
		if usr == nil && data.UserID != nil {
			err = h.userRepo.AssignEmailToUser(ctx, *data.UserID, data.Email)
			if err != nil {
				log.Printf("[error] cannot assign email to user - %d: %v\n", *data.UserID, err)
				msg.Nack(false, true)
				continue
			}
			log.Printf("[info] assign email %s to user %d\n", data.Email, *data.UserID)
			msg.Ack(false)
			continue
		}
		if usr != nil {
			msg.Ack(false)
			continue
		}

		newPwd := newPassword(15)
		var body bytes.Buffer
		err = h.registrationTemplate.Execute(&body, htmlRegistrationTemplateData{
			Login:    data.Email,
			Password: newPwd,
		})
		if err != nil {
			log.Printf("[error] cannot execute registration template: %v\n", err)
			msg.Nack(false, true)
			continue
		}

		_, err = h.userRepo.CreateUser(ctx, user.User{
			Email:    pointer.ToString(data.Email),
			Password: pointer.ToString(auth.GeneratePassword(ctx, newPwd)),
		})
		if err != nil {
			log.Printf("[error] cannot create user: %v\n", err)
			msg.Nack(false, true)
			continue
		}
		log.Printf("[info] create user with email %s\n", data.Email)

		err = h.sendToEmailDatabus.PublishDatabusSendToEmail(ctx, databus.SendToEmailDTO{
			Email:   data.Email,
			Subject: "Регистрация на сайте Godzilla Soft",
			Body:    body.String(),
		})
		if err != nil {
			log.Printf("[error] cannot send to user: %v\n", err)
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
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
