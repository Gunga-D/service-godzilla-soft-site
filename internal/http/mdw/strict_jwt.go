package mdw

import (
	"context"
	"net/http"
	"strings"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type strictJwtMDW struct {
	jwtService jwtService
}

func NewStrictJWT(jwtService jwtService) *strictJwtMDW {
	return &strictJwtMDW{
		jwtService: jwtService,
	}
}

func (m *strictJwtMDW) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			api.Return401("Авторизационный токен пустой", w)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			api.Return401("Ошибка авторизации", w)
			return
		}
		if len(headerParts[1]) == 0 {
			api.Return401("Ошибка авторизации", w)
			return
		}

		userID, userEmail, err := m.jwtService.ParseToken(headerParts[1])
		if err != nil {
			api.Return401("Ошибка авторизации", w)
			return
		}
		ctxWithUserID := context.WithValue(r.Context(), user.MetaUserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctxWithUserID, user.MetaUserEmailKey{}, userEmail)))
	})
}
