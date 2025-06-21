package user_keys

import (
	"log"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type handler struct {
	orderRepo orderRepo
}

func NewHandler(orderRepo orderRepo) *handler {
	return &handler{
		orderRepo: orderRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userEmail *string
		if v, ok := r.Context().Value(user.MetaUserEmailKey{}).(*string); ok {
			userEmail = v
		}
		if userEmail == nil {
			api.ReturnOK([]UserKeyDTO{}, w)
			return
		}

		orders, err := h.orderRepo.FetchUserOrdersByEmail(r.Context(), *userEmail)
		if err != nil {
			log.Printf("[error] fetch user keys by email: %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		res := make([]UserKeyDTO, 0, len(orders))
		for _, ord := range orders {
			itemName := "Неизвестно"
			if ord.ItemName != nil {
				itemName = *ord.ItemName
			}
			if ord.Amount == 0 {
				itemName = "Подарок " + itemName
			}
			itemSlip := "Нет инструкции"
			if ord.ItemSlip != nil {
				itemSlip = *ord.ItemSlip
			}
			res = append(res, UserKeyDTO{
				ID:       ord.ID,
				ItemName: itemName,
				ItemSlip: itemSlip,
				Code:     ord.CodeValue,
				Status:   ord.Status,
			})
		}
		api.ReturnOK(res, w)
	}
}
