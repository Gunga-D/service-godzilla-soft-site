package user_login

import (
	"net/http"
	"time"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	jwt "github.com/golang-jwt/jwt/v5"
)

type handler struct {
	jwtSecretKey string
	userRepo     user.Repository
	pwdValidator pwdValidator
}

func NewHandler(jwtSecretKey string, userRepo user.Repository, pwdValidator pwdValidator) *handler {
	return &handler{
		jwtSecretKey: jwtSecretKey,
		userRepo:     userRepo,
		pwdValidator: pwdValidator,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserLoginRequest
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
			api.Return404("Пользователь не найден", w)
			return
		}

		if ok := h.pwdValidator.ValidatePassword(r.Context(), usr.Password, req.Password); !ok {
			api.Return400("Пароль или почта введены некорректно", w)
			return
		}

		payload := jwt.MapClaims{
			"sub": usr.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		accessToken, err := token.SignedString(h.jwtSecretKey)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(UserLoginResponsePayload{
			UserID:      usr.ID,
			AccessToken: accessToken,
		}, w)
	}
}
