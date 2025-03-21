package item_details

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	itemGetter   itemGetter
	yandexGetter yandexGetter
}

func NewHandler(itemGetter itemGetter, yandexGetter yandexGetter) *handler {
	return &handler{
		itemGetter:   itemGetter,
		yandexGetter: yandexGetter,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Query().Get("item_id"), 10, 64)
		if err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		item, err := h.itemGetter.GetItemByID(r.Context(), id)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if item == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		var oldPrice *float64
		if item.OldPrice != nil {
			oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
		}

		itemDTO := ItemDTO{
			ID:            item.ID,
			Title:         item.Title,
			Description:   item.Description,
			CategoryID:    item.CategoryID,
			Platform:      item.Platform,
			Region:        item.Region,
			Publisher:     item.Publisher,
			Creator:       item.Creator,
			ReleaseDate:   item.ReleaseDate,
			CurrentPrice:  float64(item.CurrentPrice) / 100,
			IsForSale:     item.IsForSale,
			OldPrice:      oldPrice,
			ThumbnailURL:  item.ThumbnailURL,
			BackgroundURL: item.BackgroundURL,
			BxImageURL:    item.BxImageURL,
			BxGalleryUrls: item.BxGalleryUrls,
			Slip:          item.Slip,
		}

		if item.YandexID != nil {
			yaItem := h.yandexGetter.GetYandexItem(*item.YandexID)
			if yaItem != nil {
				yaBlock := YandexMarketDTO{
					Rating: yaItem.Rating,
					Price:  yaItem.Price,
				}
				itemDTO.YandexMarket = &yaBlock
			}
		}

		logger.Get().Log(fmt.Sprintf("👀 Товар\"%s\" просмотрели", item.Title))

		api.ReturnOK(itemDTO, w)
	}
}
