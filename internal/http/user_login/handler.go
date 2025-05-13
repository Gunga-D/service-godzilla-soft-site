package user_login

import (
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
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
		if usr.Password == nil {
			api.Return400("Пароль не задан", w)
			return
		}

		if ok := auth.ValidatePassword(r.Context(), *usr.Password, req.Password); !ok {
			api.Return400("Пароль или почта введены некорректно", w)
			return
		}

		accessToken, err := h.jwtService.GenerateToken(usr.ID, usr.Email)
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
