package admin_recalc_price

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	item_info "github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	sq "github.com/Masterminds/squirrel"
)

const (
	rateKztToRub = 5.7
)

type handler struct {
	itemRepo    item_info.Repository
	steamClient steam.Client
}

func NewHandler(itemRepo item_info.Repository, steamClient steam.Client) *handler {
	return &handler{
		itemRepo:    itemRepo,
		steamClient: steamClient,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("got request to recalc prices")
		go func() {
			ctx := context.Background()

			allItems, err := h.itemRepo.FetchItemsByFilter(ctx, sq.And{
				sq.Eq{"is_steam_gift": true},
				sq.Eq{"status": "active"},
				sq.NotEq{"steam_app_id": nil},
			}, 1000, 0, false)
			if err != nil {
				log.Printf("cannot fetch all items to recalc: %v\n", err)
				return
			}

			steamAppIds := make([]string, 0, len(allItems))
			for _, item := range allItems {
				steamAppIds = append(steamAppIds, fmt.Sprint(*item.SteamAppID))
			}
			kztPricesResp, err := h.steamClient.FetchPrices(ctx, steamAppIds, pointer.ToString("KZ"))
			if err != nil {
				log.Printf("cannot get prices from kz steam: %v\n", err)
				return
			}
			kztMapPrices := *kztPricesResp

			ruPricesResp, err := h.steamClient.FetchPrices(ctx, steamAppIds, pointer.ToString("RU"))
			if err != nil {
				log.Printf("cannot get prices from ru steam: %v\n", err)
				return
			}
			ruMapPrices := *ruPricesResp

			for _, item := range allItems {
				itemLoc := "KZ"
				kzPrice := kztMapPrices[fmt.Sprint(*item.SteamAppID)]
				kzInRubPrice := float64(kzPrice.Data.PriceOverview.Final) / rateKztToRub

				selfPrice := kzInRubPrice / 100

				ruPrice, foundInRu := ruMapPrices[fmt.Sprint(*item.SteamAppID)]
				if foundInRu && ruPrice.Success {
					if float64(ruPrice.Data.PriceOverview.Final)/100 < selfPrice {
						selfPrice = float64(ruPrice.Data.PriceOverview.Final) / 100
						itemLoc = "RU"
					}
				}

				itemPrice := int64(selfPrice*1.175) * 100
				itemLimitPrice := int64(selfPrice*1.084) * 100
				err = h.itemRepo.UpdatePrice(ctx, item.ID, itemPrice, itemLimitPrice, itemLoc)
				if err != nil {
					log.Printf("cannot update price of item - %s: %v\n", item.Title, err)
				}
			}
		}()
		api.ReturnOK(nil, w)
	}
}
