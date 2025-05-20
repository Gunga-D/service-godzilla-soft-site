package steam_gift_resolve_profile

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	steamClient steam.Client
}

func NewHandler(steamClient steam.Client) *handler {
	return &handler{
		steamClient: steamClient,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SteamGiftResolveProfileRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		if body.ProfileURL == "" {
			api.Return400("URL на профиль - пустой", w)
			return
		}
		u, err := url.Parse(body.ProfileURL)
		if err != nil {
			api.Return400("Ссылка на профиль некорректная", w)
			return
		}
		var steamBase64ID int64
		// Если ссылка на профиль следующего вида:
		// https://steamcommunity.com/id/rombewfwf/
		if strings.Contains(u.Path, "/id/") {
			steamID := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/id/"), "/")
			if steamID == "" {
				api.Return400("Ссылка на профиль некорректная", w)
				return
			}

			steamBase64ID, err = h.steamClient.ResolveProfileID(r.Context(), steamID)
			if err != nil {
				logger.Get().Log(fmt.Sprintf("❗️ Не получилось получить профиль по ссылке: %s", body.ProfileURL))

				api.Return500("Что-то пошло не так, попробуйте чуть позже", w)
				return
			}
		}
		// Если ссылка на профиль следующего вида:
		// https://steamcommunity.com/profiles/76561198284384725/
		if strings.Contains(u.Path, "/profiles/") {
			steamID := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/profiles/"), "/")
			if steamID == "" {
				api.Return400("Ссылка на профиль некорректная", w)
				return
			}
			steamBase64ID, err = strconv.ParseInt(steamID, 10, 64)
			if err != nil {
				logger.Get().Log(fmt.Sprintf("❗️ Не получилось получить профиль по ссылке: %s", body.ProfileURL))

				api.Return500("Что-то пошло не так, попробуйте чуть позже", w)
				return
			}
		}

		profileInfo, err := h.steamClient.GetProfileInfo(r.Context(), steamBase64ID)
		if err != nil {
			log.Printf("error to get steam profile info: %v", err)
			api.Return500("Что-то пошло не так, попробуйте чуть позже", w)
			return
		}
		var avatarURL *string
		if profileInfo.AvatarUrl != "0000000000000000000000000000000000000000" {
			avatarURL = pointer.ToString(fmt.Sprintf("https://avatars.steamstatic.com/%s.jpg", profileInfo.AvatarUrl))
		}

		api.ReturnOK(SteamGiftResolveProfileResponse{
			AvatarURL:   avatarURL,
			ProfileName: profileInfo.PersonaName,
		}, w)
	}
}
