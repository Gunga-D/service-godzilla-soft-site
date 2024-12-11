package admin_create_item

import (
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type handler struct {
	itemsRepo item.Repository
}

func NewHandler(itemsRepo item.Repository) *handler {
	return &handler{
		itemsRepo: itemsRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminCreateItemRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("cannot parse request body", w)
			return
		}

		newItemID, err := h.itemsRepo.CreateItem(r.Context(), item.Item{
			Title:        req.Title,
			Description:  req.Description,
			CategoryID:   req.CategoryID,
			Platform:     req.Platform,
			Region:       req.Region,
			CurrentPrice: req.CurrentPrice,
			IsForSale:    req.IsForSale,
			OldPrice:     req.OldPrice,
			ThumbnailURL: req.ThumbnailURL,
			Status:       item.ActiveStatus,
			Slip:         req.Slip,
		})
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}

		api.ReturnOK(AdminCreateItemResponsePayload{
			ID: newItemID,
		}, w)
	}
}
