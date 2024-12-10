package categories_tree

import (
	"net/http"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/category"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res []CategoryDTO
		for _, v := range category.CategoriesByID {
			if v.ParentID == nil {
				res = append(res, CategoryDTO{
					ID:       v.ID,
					Name:     v.Name,
					Children: categoryChildren(v.ID),
				})
			}
		}

		api.ReturnOK(res, w)
	}
}

func categoryChildren(categoryID int64) []CategoryDTO {
	var res []CategoryDTO
	for _, v := range category.CategoriesByID {
		if v.ParentID != nil && *v.ParentID == categoryID {
			res = append(res, CategoryDTO{
				ID:       v.ID,
				Name:     v.Name,
				Children: categoryChildren(v.ID),
			})
		}
	}
	return res
}
