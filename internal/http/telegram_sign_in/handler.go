package telegram_sign_in

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type handler struct {
	jwtService   jwtService
	userRepo     user.Repository
	authBotToken string
}

func NewHandler(jwtService jwtService, userRepo user.Repository, authBotToken string) *handler {
	return &handler{
		jwtService:   jwtService,
		userRepo:     userRepo,
		authBotToken: authBotToken,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TelegramSignInRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Ошибка запроса, отправляемые данные некорректные", w)
			return
		}

		if ok := h.validateRequest(req); !ok {
			api.Return400("Ошибка запроса, некорректный телеграмм хэш", w)
			return
		}

		usr, err := h.userRepo.GetUserByTelegramID(r.Context(), req.ID)
		if err != nil {
			api.Return500("Произошла непредвиденная ошибка, попробуйте чуть позже", w)
			return
		}

		var userID int64
		var userEmail *string
		if usr == nil {
			createdUserID, err := h.userRepo.CreateUser(r.Context(), user.User{
				PhotoURL:   req.PhotoURL,
				Username:   req.Username,
				FirstName:  req.FirstName,
				TelegramID: &req.ID,
			})
			if err != nil {
				api.Return500("Произошла непредвиденная ошибка, попробуйте чуть позже", w)
				return
			}
			userID = createdUserID
		} else {
			userID = usr.ID
			userEmail = usr.Email
		}

		accessToken, err := h.jwtService.GenerateToken(userID, userEmail)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		api.ReturnOK(TelegramSignInResponsePayload{
			UserID:      userID,
			AccessToken: accessToken,
		}, w)
	}
}

func (h *handler) validateRequest(r TelegramSignInRequest) bool {
	dataToCheck := []string{}
	// Определение полей происходит в алфавитном порядке
	if r.AuthDate != nil {
		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%d", "auth_date", *r.AuthDate))
	}
	if r.FirstName != nil {
		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%s", "first_name", *r.FirstName))
	}
	dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%d", "id", r.ID))
	if r.LastName != nil {
		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%s", "last_name", *r.LastName))
	}
	if r.PhotoURL != nil {
		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%s", "photo_url", *r.PhotoURL))
	}
	if r.Username != nil {
		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%s", "username", *r.Username))
	}

	sha256hash := sha256.New()
	io.WriteString(sha256hash, h.authBotToken)
	hmachash := hmac.New(sha256.New, sha256hash.Sum(nil))
	io.WriteString(hmachash, strings.Join(dataToCheck, "\n"))
	hash := hex.EncodeToString(hmachash.Sum(nil))

	return r.Hash == hash
}
