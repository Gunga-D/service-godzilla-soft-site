package mdw

import (
	"context"
	"net/http"
	"strings"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/platform"
)

type useragentMDW struct{}

func NewUseragent() *useragentMDW {
	return &useragentMDW{}
}

func (m *useragentMDW) IdentifyPlatform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("user-agent")
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(header, "Android") ||
			strings.Contains(header, "webOS") ||
			strings.Contains(header, "iPhone") ||
			strings.Contains(header, "iPad") ||
			strings.Contains(header, "iPod") ||
			strings.Contains(header, "BlackBerry") ||
			strings.Contains(header, "IEMobile") ||
			strings.Contains(header, "Opera Mini") {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), platform.MetaPlatform{}, platform.MAV)))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), platform.MetaPlatform{}, platform.Desktop)))
	})
}
