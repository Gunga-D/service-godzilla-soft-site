package check_voucher

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/voucher"
)

type handler struct {
	itemGetter        itemGetter
	voucherActivation voucherActivation
}

func NewHandler(itemGetter itemGetter, voucherActivation voucherActivation) *handler {
	return &handler{
		itemGetter:        itemGetter,
		voucherActivation: voucherActivation,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CheckVoucherReq
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		i, err := h.itemGetter.GetItemByID(r.Context(), body.ItemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if i == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		newPrice, hasLimitation, err := h.voucherActivation.PeekVoucher(r.Context(), body.Voucher, i.Item)
		if err != nil {
			if errors.Is(err, voucher.ErrNotFoundVoucher) {
				api.Return404("Купон не найден или уже был активирован", w)
				return
			}
			log.Printf("[error] peek voucher error: %v\n", err)
			api.Return500("Неизвестная ошибка", w)
			return
		}

		normNewPrice := float64(newPrice) / 100

		var warning *string
		if hasLimitation {
			warning = pointer.ToString(fmt.Sprintf("К сожалению, мы можем снизить цену только до %.2f ₽", normNewPrice))
		}

		api.ReturnOK(CheckVoucherResp{
			OldPrice: float64(i.CurrentPrice) / 100,
			NewPrice: normNewPrice,
			Currency: "RUB",
			Warning:  warning,
		}, w)
	}
}
