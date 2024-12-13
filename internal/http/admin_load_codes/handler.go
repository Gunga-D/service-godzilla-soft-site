package admin_load_codes

import (
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type handler struct {
	codeRepo               codeRepo
	itemRepo               itemRepo
	itemChangeStateDatabus itemChangeStateDatabus
}

func NewHandler(codeRepo codeRepo, itemRepo itemRepo, itemChangeStateDatabus itemChangeStateDatabus) *handler {
	return &handler{
		codeRepo:               codeRepo,
		itemRepo:               itemRepo,
		itemChangeStateDatabus: itemChangeStateDatabus,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminLoadCodeRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		i, err := h.itemRepo.GetItemByID(r.Context(), req.ItemID)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		if i == nil {
			api.Return400("Такого товара не существует", w)
			return
		}

		if err := h.codeRepo.CreateCodes(r.Context(), req.ItemID, req.Codes); err != nil {
			api.Return500(err.Error(), w)
			return
		}
		if i.Status == item.PausedStatus {
			h.itemChangeStateDatabus.PublishDatabusChangeItemState(r.Context(), databus.ChangeItemStateDTO{
				ItemID: req.ItemID,
				Status: item.ActiveStatus,
			})
		}
		api.ReturnOK(nil, w)
	}
}
