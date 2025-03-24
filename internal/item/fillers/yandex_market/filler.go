package yandex_market

import (
	"context"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_market"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type filler struct {
	yandexMarket yandex_market.Client
}

func NewFiller(yandexMarket yandex_market.Client) *filler {
	return &filler{
		yandexMarket: yandexMarket,
	}
}

func (f *filler) Fill(ctx context.Context, items []item.ItemCache) error {
	var yandexSkus []string
	for _, i := range items {
		if i.YandexID != nil {
			yandexSkus = append(yandexSkus, *i.YandexID)
		}
	}
	if len(yandexSkus) == 0 {
		return nil
	}

	var modelIds []int64
	prices := make(map[string]float64)
	mapSkuModelID := make(map[string]int64)
	offerMappings, err := f.yandexMarket.OfferMappings(ctx, yandex_market.OfferMappingsRequest{
		OfferIds: yandexSkus,
	})
	if err != nil {
		return err
	}

	for _, offerMap := range offerMappings.Result.OfferMappings {
		prices[offerMap.Offer.OfferId] = offerMap.Offer.BasicPrice.Value
		modelIds = append(modelIds, offerMap.Mapping.MarketModelId)
		mapSkuModelID[offerMap.Offer.OfferId] = offerMap.Mapping.MarketModelId
	}

	sumRating := make(map[int64]float64)
	countRating := make(map[int64]int)
	var nextPageToken *string
	for {
		goodsFeedback, err := f.yandexMarket.GoodsFeedback(ctx, yandex_market.GoodsFeedbackRequest{
			ModelIds: modelIds,
		}, nextPageToken)
		if err != nil {
			log.Printf("[error] cannot get goods feedback for %v: %v\n", modelIds, err)
			return nil
		}
		for _, f := range goodsFeedback.Result.Feedbacks {
			sumRating[f.Identifiers.ModelID] += float64(f.Statistics.Rating)
			countRating[f.Identifiers.ModelID]++
		}
		nextPageToken = goodsFeedback.Result.Paging.NextPageToken
		if nextPageToken == nil {
			break
		}
	}

	for idx := 0; idx < len(items); idx++ {
		v := items[idx]
		if v.YandexID != nil {
			modelID := mapSkuModelID[*v.YandexID]
			rating := sumRating[modelID] / float64(countRating[modelID])

			items[idx].YandexMarket = &item.ItemYandexMarketBlock{
				Price:  prices[*v.YandexID],
				Rating: rating,
			}
		}
	}
	return nil
}
