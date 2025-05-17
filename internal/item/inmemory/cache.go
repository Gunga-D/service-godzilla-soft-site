package inmemory

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/fillers"
)

const (
	_itemsPerRequestLimit = 50
)

type cache struct {
	itemsByID  map[int64]item.ItemCache
	itemByName map[string]item.ItemCache
	fillers    []fillers.Filler
	getter     getter
	rec        recomendation
	lock       *sync.RWMutex
}

func NewCache(getter getter, fillers []fillers.Filler, rec recomendation) *cache {
	return &cache{
		itemsByID:  make(map[int64]item.ItemCache),
		itemByName: make(map[string]item.ItemCache),
		getter:     getter,
		rec:        rec,
		lock:       &sync.RWMutex{},
		fillers:    fillers,
	}
}

func (c *cache) WarmUp(ctx context.Context) {
	if err := c.lazySync(ctx, _itemsPerRequestLimit); err != nil {
		log.Printf("error to warmup cache: %v\n", err)
	}
}

func (c *cache) lazySync(ctx context.Context, limit uint64) error {
	defer func(start time.Time) {
		log.Printf("[info] sync lazy items latency: %v\n", time.Since(start))
	}(time.Now())

	newItemsByID := make(map[int64]item.ItemCache)
	newItemsByName := make(map[string]item.ItemCache)
	itemsBySteamAppID := make(map[int64]item.ItemCache)

	cursor := int64(0)
	for {
		gotItems, err := c.getter.FetchItemsPaginatedCursorItemId(ctx, limit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		var notCachedItems []item.ItemCache
		for _, gotItem := range gotItems {
			currentItem, err := c.GetItemByID(ctx, gotItem.ID)
			if err != nil {
				continue
			}

			if currentItem == nil {
				notCachedItems = append(notCachedItems, item.ItemCache{
					Item: gotItem,
				})
			} else {
				newItemsByID[gotItem.ID] = item.ItemCache{
					Item:         gotItem,
					YandexMarket: currentItem.YandexMarket,
					SteamBlock:   currentItem.SteamBlock,
				}
				newItemsByName[gotItem.Title] = item.ItemCache{
					Item:         gotItem,
					YandexMarket: currentItem.YandexMarket,
					SteamBlock:   currentItem.SteamBlock,
				}
				if gotItem.SteamAppID != nil {
					itemsBySteamAppID[*gotItem.SteamAppID] = item.ItemCache{
						Item:         gotItem,
						YandexMarket: currentItem.YandexMarket,
						SteamBlock:   currentItem.SteamBlock,
					}
				}
			}
		}

		for _, f := range c.fillers {
			err = f.Fill(ctx, notCachedItems)
			if err != nil {
				return err
			}
		}

		for _, i := range notCachedItems {
			newItemsByID[i.ID] = i
			newItemsByName[i.Title] = i
			if i.SteamAppID != nil {
				itemsBySteamAppID[*i.SteamAppID] = i
			}
		}

		if len(gotItems) < int(limit) {
			break
		}
		cursor = gotItems[len(gotItems)-1].ID
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.itemsByID = newItemsByID
	c.itemByName = newItemsByName

	go c.rec.Sync(context.Background(), itemsBySteamAppID)
	return nil
}

func (c *cache) StartSync(ctx context.Context) {
	if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
		log.Fatalf("[error] failed to sync items: %v\n", err)
	}

	t := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			log.Println("[info] sync items stop")
			return
		case <-t.C:
			if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
				log.Printf("[error] failed to sync items: %v\n", err)
			}
		}
	}
}

func (c *cache) sync(ctx context.Context, limit uint64) error {
	defer func(start time.Time) {
		log.Printf("[info] sync items latency: %v\n", time.Since(start))
	}(time.Now())

	cursor := int64(0)
	newItemsByID := make(map[int64]item.ItemCache)
	newItemsByName := make(map[string]item.ItemCache)
	itemsBySteamAppID := make(map[int64]item.ItemCache)
	for {
		gotItems, err := c.getter.FetchItemsPaginatedCursorItemId(ctx, limit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		cacheItems := make([]item.ItemCache, 0, len(gotItems))
		for _, gotItem := range gotItems {
			cacheItems = append(cacheItems, item.ItemCache{
				Item: gotItem,
			})
		}

		for _, f := range c.fillers {
			err = f.Fill(ctx, cacheItems)
			if err != nil {
				return err
			}
		}

		for _, i := range cacheItems {
			newItemsByID[i.ID] = i
			newItemsByName[i.Title] = i
			if i.SteamAppID != nil {
				itemsBySteamAppID[*i.SteamAppID] = i
			}
		}

		if len(gotItems) < int(limit) {
			break
		}
		cursor = gotItems[len(gotItems)-1].ID
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.itemsByID = newItemsByID
	c.itemByName = newItemsByName

	go c.rec.Sync(context.Background(), itemsBySteamAppID)

	return nil
}

func (c *cache) GetItemByID(ctx context.Context, id int64) (*item.ItemCache, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	res, found := c.itemsByID[id]
	if !found {
		return nil, nil
	}
	return &res, nil
}

func (c *cache) GetItemByName(ctx context.Context, name string) (*item.ItemCache, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	res, found := c.itemByName[name]
	if !found {
		return nil, nil
	}
	return &res, nil
}

func (c *cache) FetchAllItems(ctx context.Context) ([]item.ItemCache, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var res []item.ItemCache
	for _, i := range c.itemsByID {
		res = append(res, i)
	}
	return res, nil
}
