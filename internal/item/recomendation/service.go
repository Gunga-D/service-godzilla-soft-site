package recomendation

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type service struct {
	steamClient steam.Client
	genreItems  map[string][]item.ItemCache
	lock        *sync.RWMutex
}

func NewService(steamClient steam.Client) *service {
	return &service{
		steamClient: steamClient,
		genreItems:  make(map[string][]item.ItemCache),
		lock:        &sync.RWMutex{},
	}
}

func (s *service) Sync(ctx context.Context, itemsBySteamAppID map[int64]item.ItemCache) error {
	defer func(start time.Time) {
		log.Printf("[info] actualize recommendation cache latency: %v\n", time.Since(start))
	}(time.Now())

	for genreName, genreID := range mapGenreNameToID {
		if _, skip := skipGenres[genreName]; skip {
			continue
		}

		time.Sleep(50 * time.Millisecond)
		resp, err := s.steamClient.GetGenreApps(ctx, genreID)
		if err != nil {
			return err
		}
		var itemsForGenre []item.ItemCache
		for _, steamItem := range resp.Tabs.TopSellers.Items {
			i, ok := itemsBySteamAppID[steamItem.ID]
			if !ok {
				continue
			}

			itemsForGenre = append(itemsForGenre, i)
			s.lock.Lock()
			s.genreItems[genreName] = itemsForGenre
			s.lock.Unlock()
		}
	}
	return nil
}

func (s *service) Recommend(ctx context.Context, itemID int64, genres []string) ([]item.ItemCache, error) {
	var res []item.ItemCache
	for _, genre := range genres {
		if _, skip := skipGenres[genre]; skip {
			continue
		}

		s.lock.RLock()
		cachedItems, exists := s.genreItems[genre]
		if !exists {
			s.lock.RUnlock()
			continue
		}
		for _, cachedItem := range cachedItems {
			if cachedItem.ID == itemID {
				continue
			}
			res = append(res, cachedItem)
		}
		s.lock.RUnlock()
	}
	return res, nil
}
