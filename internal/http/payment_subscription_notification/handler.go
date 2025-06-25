package payment_subscription_notification

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	terminalKey string
	subRepo     subRepo
}

func NewHandler(terminalKey string, subRepo subRepo) *handler {
	return &handler{
		terminalKey: terminalKey,
		subRepo:     subRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PaymentNotificationRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Ошибка запроса, отправляемые данные некорректные", w)
			return
		}

		if req.TerminalKey != h.terminalKey {
			api.Return401("Ключ некорректный", w)
			return
		}

		if req.Status == "CONFIRMED" {
			logger.Get().Log("💸 Подписка оплачена")

			err := h.subRepo.PaidSubscriptionBill(r.Context(), req.OrderID, fmt.Sprint(req.RebillId))
			if err != nil {
				api.Return500("Неизвестная ошибка, попробуйте чуть позже", w)
				return
			}
		} else {
			log.Printf("[INFO][status - %s][orderId - %s][errorCode - %s] Unsuccessful status of payment\n", req.Status, req.OrderID, req.ErrorCode)
			if req.Status == "AUTH_FAIL" || req.Status == "REJECTED" || req.Status == "DEADLINE_EXPIRED" {
				err := h.subRepo.FailedSubscriptionBill(r.Context(), req.OrderID)
				if err != nil {
					api.Return500("Неизвестная ошибка, попробуйте чуть позже", w)
					return
				}
			}
		}

		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	}
}
