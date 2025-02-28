package steam_calc_price

import (
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/steam"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SteamCalcPriceRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		if body.Amount < 30 {
			api.Return400("Минимальная сумма от 30 рублей", w)
			return
		}
		if body.Amount > 30000 {
			api.Return400("Максимальная сумма до 30 000 рублей", w)
			return
		}

		api.ReturnOK(SteamInvoiceResponse{
			Price: steam.SteamCalcPrice(body.Amount),
		}, w)
	}
}
