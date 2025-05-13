package mdw

import (
	"context"
	"net/http"
	"strings"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type jwtMDW struct {
	jwtService jwtService
}

func NewJWT(jwtService jwtService) *jwtMDW {
	return &jwtMDW{
		jwtService: jwtService,
	}
}

func (m *jwtMDW) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}
		if len(headerParts[1]) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		userID, userEmail, err := m.jwtService.ParseToken(headerParts[1])
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), user.MetaUserIDKey{}, userID)
		if userEmail != nil {
			ctx = context.WithValue(ctx, user.MetaUserEmailKey{}, userEmail)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
