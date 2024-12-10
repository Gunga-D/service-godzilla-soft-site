package admin_create_item

import "net/http"

type handler struct {
	itemsRepo itemsRepo
}

func NewHandler(itemsRepo itemsRepo) *handler {
	return &handler{
		itemsRepo: itemsRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
