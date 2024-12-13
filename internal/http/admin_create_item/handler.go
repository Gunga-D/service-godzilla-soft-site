package admin_create_item

import (
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/category"
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
			api.Return400("Невалидный запрос", w)
			return
		}
		var oldPrice *int64
		if req.OldPrice != nil && req.IsForSale {
			oldPrice = pointer.ToInt64(int64(*req.OldPrice * 100))
		}

		if _, found := category.CategoriesByID[req.CategoryID]; !found {
			api.Return400("Нет такой категории", w)
			return
		}

		newItemID, err := h.itemsRepo.CreateItem(r.Context(), item.Item{
			Title:        req.Title,
			Description:  req.Description,
			CategoryID:   req.CategoryID,
			Platform:     req.Platform,
			Region:       req.Region,
			CurrentPrice: int64(req.CurrentPrice * 100),
			IsForSale:    req.IsForSale,
			OldPrice:     oldPrice,
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
