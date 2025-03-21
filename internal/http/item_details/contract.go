package item_details

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/yandex_market_filler"
)

type itemGetter interface {
	GetItemByID(ctx context.Context, id int64) (*item.Item, error)
}

type yandexGetter interface {
	GetYandexItem(yandexSku string) *yandex_market_filler.YandexMarketOffer
}
