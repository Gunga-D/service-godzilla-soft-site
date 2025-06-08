package suggest

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
)

const (
	limit = 50
)

type service struct {
	g         getter
	itemCache itemCache
	indexDesc suggest.IndexDescription
	suggester *suggest.Service
}

func NewService(g getter, itemCache itemCache) *service {
	s := &service{
		g: g,
		indexDesc: suggest.IndexDescription{
			Name:      "items",
			NGramSize: 3,
			Wrap:      [2]string{"$", "$"},
			Pad:       "$",
			Alphabet:  []string{"english", "$"},
		},
		itemCache: itemCache,
	}
	if err := s.sync(context.Background()); err != nil {
		log.Fatalf("[error] failed to sync items for suggest: %v\n", err)
	}

	return s
}

func (s *service) StartSync(ctx context.Context) {
	t := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("[info] sync items for suggest stop")
			return
		case <-t.C:
			if err := s.sync(ctx); err != nil {
				log.Printf("[error] failed to sync items for suggest: %v\n", err)
			}
		}
	}
}

func (s *service) sync(ctx context.Context) error {
	cursor := int64(0)
	var items []string
	for {
		gotItems, err := s.g.FetchItemsPaginatedCursorItemId(ctx, limit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		for _, gotItem := range gotItems {
			if _, ok := item.NotShowedItems[gotItem.ID]; ok {
				continue
			}
			items = append(items, gotItem.Title)
		}

		if len(gotItems) < int(limit) {
			break
		}
		cursor = gotItems[len(gotItems)-1].ID
	}

	dict := dictionary.NewInMemoryDictionary(items)
	suggestBuilder, err := suggest.NewRAMBuilder(dict, s.indexDesc)
	if err != nil {
		return err
	}
	service := suggest.NewService()
	if err := service.AddIndex(s.indexDesc.Name, dict, suggestBuilder); err != nil {
		return err
	}
	s.suggester = service
	return nil
}

func (s *service) Suggest(ctx context.Context, text string) ([]Suggested, error) {
	if s.suggester == nil {
		return []Suggested{}, nil
	}
	normQuery := strings.ToLower(text)

	// Обработка аббревиатур
	fullText, found := abbr[normQuery]
	if found {
		text = fullText
	}

	searchConf, err := suggest.NewSearchConfig(text, 5, metric.DiceMetric(), 0.1)
	if err != nil {
		return nil, err
	}
	suggests, err := s.suggester.Suggest(s.indexDesc.Name, searchConf)
	if err != nil {
		return nil, err
	}

	res := make([]Suggested, 0, len(suggests))

	res = append(res, Suggested{
		Type: "banner",
		Banner: &SuggestedBanner{
			Image:       "https://disk.godzillasoft.ru/random_game_banner.png",
			Title:       "Случайная Steam игра",
			Description: "Испытай удачу и выиграй заветную игру всего лишь за 208₽",
			URL:         "https://godzillasoft.ru/random",
		},
		Probability: 1.0,
	})

	for _, suggest := range suggests {
		i, err := s.itemCache.GetItemByName(ctx, suggest.Value)
		if err != nil {
			continue
		}
		if i == nil {
			continue
		}

		res = append(res, Suggested{
			Type:        "item",
			Item:        i,
			Probability: suggest.Score,
		})
	}
	return res, nil
}
