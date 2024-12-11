package mdw

import (
	"context"
	"net/http"
	"strings"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtMDDW struct {
	jwtSecretKey string
}

func NewJWT(jwtSecretKey string) *jwtMDDW {
	return &jwtMDDW{
		jwtSecretKey: jwtSecretKey,
	}
}

func (m *jwtMDDW) VerifyUser(next http.Handler) http.Handler {
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

		token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecretKey, nil
		})
		if err != nil {
			api.Return401("Ошибка авторизации", w)
			return
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			api.Return401("Ошибка авторизации", w)
			return
		}
		userID, err := claims.GetSubject()
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), user.MetaUserIDKey{}, userID)))
	})
}
