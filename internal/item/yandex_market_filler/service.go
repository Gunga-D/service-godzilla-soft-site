package yandex_market_filler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_market"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

const fetchLimit = 10

type service struct {
	yandexItems  map[string]YandexMarketOffer
	itemGetter   itemGetter
	yandexMarket yandex_market.Client
	lock         *sync.RWMutex
}

func NewService(itemGetter itemGetter, yandexMarket yandex_market.Client) *service {
	return &service{
		yandexItems:  make(map[string]YandexMarketOffer),
		itemGetter:   itemGetter,
		yandexMarket: yandexMarket,
		lock:         &sync.RWMutex{},
	}
}

func (s *service) StartFetching(ctx context.Context) {
	if err := s.fetch(ctx); err != nil {
		log.Printf("[error] failed to sync items for suggest: %v\n", err)
	}

	t := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("[info] sync items for suggest stop")
			return
		case <-t.C:
			if err := s.fetch(ctx); err != nil {
				log.Printf("[error] failed to sync items for suggest: %v\n", err)
			}
		}
	}
}

func (s *service) fetch(ctx context.Context) error {
	cursor := int64(0)
	yandexMarketOffers := make(map[string]YandexMarketOffer)
	for {
		gotItems, err := s.itemGetter.FetchItemsPaginatedCursorItemId(ctx, fetchLimit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		var yandexSkus []string
		for _, gotItem := range gotItems {
			if _, ok := item.NotShowedItems[gotItem.ID]; ok {
				continue
			}
			if gotItem.YandexID != nil {
				yandexSkus = append(yandexSkus, *gotItem.YandexID)
			}
		}

		var modelIds []int64
		prices := make(map[string]float64)
		mapSkuModelID := make(map[string]int64)
		offerMappings, err := s.yandexMarket.OfferMappings(ctx, yandex_market.OfferMappingsRequest{
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
		goodsFeedback, err := s.yandexMarket.GoodsFeedback(ctx, yandex_market.GoodsFeedbackRequest{
			ModelIds: modelIds,
		})
		if err != nil {
			return err
		}
		for _, f := range goodsFeedback.Result.Feedbacks {
			sumRating[f.Identifiers.ModelID] += float64(f.Statistics.Rating)
			countRating[f.Identifiers.ModelID]++
		}

		for _, sku := range yandexSkus {
			modelID := mapSkuModelID[sku]
			rating := sumRating[modelID] / float64(countRating[modelID])
			yandexMarketOffers[sku] = YandexMarketOffer{
				Price:  prices[sku],
				Rating: rating,
			}
		}

		if len(gotItems) < fetchLimit {
			break
		}
		cursor = gotItems[len(gotItems)-1].ID
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.yandexItems = yandexMarketOffers
	return nil
}

func (s *service) GetYandexItem(yandexSku string) *YandexMarketOffer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	i, ok := s.yandexItems[yandexSku]
	if !ok {
		return nil
	}
	return &i
}
