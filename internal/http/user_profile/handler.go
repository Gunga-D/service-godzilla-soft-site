package user_profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
	"go.octolab.org/pointer"
)

type handler struct {
	redis    redis.Redis
	userRepo user.Repository
}

func NewHandler(redis redis.Redis, userRepo user.Repository) *handler {
	return &handler{
		redis:    redis,
		userRepo: userRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID *int64
		if v, ok := r.Context().Value(user.MetaUserIDKey{}).(int64); ok {
			userID = pointer.ToInt64(v)
		}
		if userID == nil {
			api.Return400("Пользователь не авторизован", w)
			return
		}

		raw, err := redigo.Bytes(h.redis.Get(r.Context(), fmt.Sprintf(user.UserCacheKey, *userID)))
		if err != nil {
			if err == redigo.ErrNil {
				repoUser, err := h.userRepo.GetUserByID(r.Context(), *userID)
				if err != nil {
					api.Return500("Неизвестная ошибка", w)
					return
				}
				if repoUser == nil {
					api.Return404("Пользователь не найден", w)
					return
				}

				userRaw, err := json.Marshal(repoUser)
				if err != nil {
					api.Return500("Неизвестная ошибка", w)
					return
				}
				err = h.redis.Set(r.Context(), fmt.Sprintf(user.UserCacheKey, *userID), userRaw, nil)
				if err != nil {
					api.Return500("Неизвестная ошибка", w)
					return
				}

				api.ReturnOK(UserProfileResponsePayload{
					UserID:    repoUser.ID,
					Email:     repoUser.Email,
					SteamLink: repoUser.SteamLink,
					PhotoURL:  repoUser.PhotoURL,
					Username:  repoUser.Username,
					FirstName: repoUser.FirstName,
				}, w)
				return
			}
			api.Return500("Неизвестная ошибка", w)
			return
		}
		var cacheUser user.User
		err = json.Unmarshal(raw, &cacheUser)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		api.ReturnOK(UserProfileResponsePayload{
			UserID:    cacheUser.ID,
			Email:     cacheUser.Email,
			SteamLink: cacheUser.SteamLink,
			PhotoURL:  cacheUser.PhotoURL,
			Username:  cacheUser.Username,
			FirstName: cacheUser.FirstName,
		}, w)
	}
}
