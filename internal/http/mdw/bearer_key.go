package mdw

import (
	"net/http"
	"strings"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type bearerMDW struct {
	secret string
}

func NewBearerMDW(secret string) *bearerMDW {
	return &bearerMDW{
		secret: secret,
	}
}

func (m *bearerMDW) VerifyUser(next http.Handler) http.Handler {
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

		if headerParts[1] != m.secret {
			api.Return401("Ключ невалидный", w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
