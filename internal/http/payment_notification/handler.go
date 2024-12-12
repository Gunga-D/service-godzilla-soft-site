package payment_notification

import "net/http"

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
