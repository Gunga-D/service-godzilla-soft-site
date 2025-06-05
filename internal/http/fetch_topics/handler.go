package fetch_topics

import (
	"net/http"
)

type handler struct {
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
