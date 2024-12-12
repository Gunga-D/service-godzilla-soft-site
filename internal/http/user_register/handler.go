package user_register

import (
	"errors"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
)

var (
	errDuplicateKey = errors.New("duplicate key value violates unique constraint")
)

type handler struct {
	jwtService jwtService
	userRepo   user.Repository
}

func NewHandler(jwtService jwtService, userRepo user.Repository) *handler {
	return &handler{
		jwtService: jwtService,
		userRepo:   userRepo,
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
			Password: auth.GeneratePassword(r.Context(), req.Password),
		})
		if err != nil {
			if errors.Is(err, errDuplicateKey) {
				api.Return400("Пользователь с такой почтой уже зарегистрирован", w)
				return
			}
			api.Return500("Невозможно создать пользователя, попробуйте чуть позже", w)
			return
		}

		accessToken, err := h.jwtService.GenerateToken(userID, req.Email)
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
