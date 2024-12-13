package payment_notification

import "net/http"

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Тут должна быть реакция на успешную и неуспешную оплату:
		// 1) При успешной оплате переводим order в статус paid
		// 2) При неуспешной оплате переводим order в статус failed
	}
}
