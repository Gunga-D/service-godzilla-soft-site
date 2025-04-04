package admin_recalc_price

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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

			kzSteamAppIds := make([]string, 0, len(allItems))
			ruSteamAppIds := make([]string, 0, len(allItems))
			for _, item := range allItems {
				kzSteamAppIds = append(kzSteamAppIds, fmt.Sprint(*item.SteamAppID))

				// Тут для ru исключаем Battlefield™ V и Metro 2033 Redux
				if *item.SteamAppID != 1238810 && *item.SteamAppID != 286690 {
					ruSteamAppIds = append(ruSteamAppIds, fmt.Sprint(*item.SteamAppID))
				}
			}
			kztPricesResp, err := h.steamClient.FetchPrices(ctx, kzSteamAppIds, pointer.ToString("KZ"))
			if err != nil {
				log.Printf("cannot get prices from kz steam for %v: %v\n", strings.Join(kzSteamAppIds, ","), err)
				return
			}
			kztMapPrices := *kztPricesResp

			ruPricesResp, err := h.steamClient.FetchPrices(ctx, ruSteamAppIds, pointer.ToString("RU"))
			if err != nil {
				log.Printf("cannot get prices from ru steam for %v: %v\n", strings.Join(ruSteamAppIds, ","), err)
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
