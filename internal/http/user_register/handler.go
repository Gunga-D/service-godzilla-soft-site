package user_register

import (
	"errors"
	"net/http"
	"time"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	errDuplicateKey = errors.New("duplicate key value violates unique constraint")
)

type handler struct {
	jwtSecretKey string
	userRepo     user.Repository
	pwdGenerator pwdGenerator
}

func NewHandler(jwtSecretKey string, userRepo user.Repository, pwdGenerator pwdGenerator) *handler {
	return &handler{
		jwtSecretKey: jwtSecretKey,
		userRepo:     userRepo,
		pwdGenerator: pwdGenerator,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRegisterRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Ошибка запроса, отправляемые данные некорректные", w)
			return
		}

		userID, err := h.userRepo.CreateUser(r.Context(), user.User{
			Email:    req.Email,
			Password: h.pwdGenerator.GeneratePassword(r.Context(), req.Password),
		})
		if err != nil {
			if errors.Is(err, errDuplicateKey) {
				api.Return400("Пользователь с такой почтой уже зарегистрирован", w)
				return
			}
			api.Return500("Невозможно создать пользователя, попробуйте чуть позже", w)
			return
		}

		payload := jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		accessToken, err := token.SignedString(h.jwtSecretKey)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(UserRegisterResponsePayload{
			UserID:      userID,
			AccessToken: accessToken,
		}, w)
	}
}
