package mdw

import (
	"context"
	"fmt"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

type VerifyUUID struct {
	name string
}

func NewVerifyUUID(paramName string) *VerifyUUID {
	return &VerifyUUID{
		name: paramName,
	}
}

func (mdw *VerifyUUID) VerifyUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, mdw.name)
		uuidValue, err := uuid.Parse(uuidStr)
		if err != nil {
			api.Return400(fmt.Sprintf("В запросе задан неверный uuid: %s", err.Error()), w)
			return
		}
		ctx := context.WithValue(r.Context(), "uuid", uuidValue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
