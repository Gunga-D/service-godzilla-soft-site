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
	lock       *sync.RWMutex
}

func NewCache(getter getter, fillers []fillers.Filler) *cache {
	return &cache{
		itemsByID:  make(map[int64]item.ItemCache),
		itemByName: make(map[string]item.ItemCache),
		getter:     getter,
		lock:       &sync.RWMutex{},
		fillers:    fillers,
	}
}

func (c *cache) WarmUp(ctx context.Context) {
	if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
		log.Printf("error to warmup cache: %v\n", err)
	}
}

func (c *cache) StartSync(ctx context.Context) {
	if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
		log.Fatalf("[error] failed to sync items: %v\n", err)
	}

	t := time.NewTicker(60 * time.Minute)
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
