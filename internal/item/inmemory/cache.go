package inmemory

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

const (
	_itemsPerRequestLimit = 50
)

type cache struct {
	itemsByID  map[int64]item.Item
	itemByName map[string]item.Item
	getter     getter
	lock       *sync.RWMutex
}

func NewCache(getter getter) *cache {
	return &cache{
		itemsByID:  make(map[int64]item.Item),
		itemByName: make(map[string]item.Item),
		getter:     getter,
		lock:       &sync.RWMutex{},
	}
}

func (c *cache) WarmUp(ctx context.Context) error {
	if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
		return err
	}
	return nil
}

func (c *cache) StartSync(ctx context.Context) {
	if err := c.sync(ctx, _itemsPerRequestLimit); err != nil {
		log.Fatalf("[error] failed to sync items: %v\n", err)
	}

	t := time.NewTicker(5 * time.Minute)
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
	newItemsByID := make(map[int64]item.Item)
	newItemsByName := make(map[string]item.Item)
	for {
		gotItems, err := c.getter.FetchItemsPaginatedCursorItemId(ctx, limit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		for _, gotItem := range gotItems {
			newItemsByID[gotItem.ID] = gotItem
			newItemsByName[gotItem.Title] = gotItem
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

func (c *cache) GetItemByID(ctx context.Context, id int64) (*item.Item, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	res, found := c.itemsByID[id]
	if !found {
		return nil, nil
	}
	return &res, nil
}

func (c *cache) GetItemByName(ctx context.Context, name string) (*item.Item, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	res, found := c.itemByName[name]
	if !found {
		return nil, nil
	}
	return &res, nil
}
