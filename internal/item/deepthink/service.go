package deepthink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/deepseek"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

const (
	limit             = 50
	promptTemplate    = "У тебя есть следующий список игр:\n%s. Подбери самые лучшие игры по запросу пользователя из данного списка и верни ответ в следующем формате (строго соблюдай данный формат, ты можешь отвечать только в квадратных скобках): [твои мысли об игре без упоминания их id] | [id игр с разделительным символом точка с запятой]"
	thinkingResultKey = "thinking_result:%s"
)

type service struct {
	deepseekClient deepseek.Client
	currentPrompt  string
	itemCache      itemCache
	g              getter
	redis          redis.Redis
}

func NewService(deepseekClient deepseek.Client, itemCache itemCache, g getter, redis redis.Redis) *service {
	return &service{
		deepseekClient: deepseekClient,
		itemCache:      itemCache,
		g:              g,
		redis:          redis,
	}
}

func (s *service) StartSync(ctx context.Context) {
	if err := s.sync(ctx); err != nil {
		log.Printf("[error] failed to sync items for deepthink: %v\n", err)
	}

	t := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("[info] sync items for deepthink stop")
			return
		case <-t.C:
			if err := s.sync(ctx); err != nil {
				log.Printf("[error] failed to sync items for deepthink: %v\n", err)
			}
		}
	}
}

func (s *service) sync(ctx context.Context) error {
	defer func(start time.Time) {
		log.Printf("[info] generate AI prompt for thinking: %v\n", time.Since(start))
	}(time.Now())

	cursor := int64(0)
	var promptItems string
	for {
		gotItems, err := s.g.FetchItemsPaginatedCursorItemId(ctx, limit, cursor)
		if err != nil {
			return fmt.Errorf("failed to fetch items: %v", err)
		}

		for _, gotItem := range gotItems {
			if _, ok := item.NotShowedItems[gotItem.ID]; ok {
				continue
			}
			if gotItem.CategoryID == 10001 {
				promptItems += fmt.Sprintf("%s с id - %d\n", gotItem.Title, gotItem.ID)
			}
		}

		if len(gotItems) < int(limit) {
			break
		}
		cursor = gotItems[len(gotItems)-1].ID
	}

	s.currentPrompt = fmt.Sprintf(promptTemplate, promptItems)
	return nil
}

func (s *service) StartThinking(ctx context.Context, query string) string {
	id := uuid.NewString()
	go s.think(context.Background(), id, query)
	return id
}

func (s *service) GetThinkingResult(ctx context.Context, id string) (*ThinkResult, error) {
	raw, err := redigo.Bytes(s.redis.Get(ctx, fmt.Sprintf(thinkingResultKey, id)))
	if err != nil {
		if err == redigo.ErrNil {
			return nil, nil
		}
	}
	var res ThinkResult
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *service) think(ctx context.Context, id string, query string) error {
	log.Printf("[info] Промпт на генерацию: %s\n\nЗапрос: %s\n", s.currentPrompt, query)
	resp, err := s.deepseekClient.Completions(ctx, deepseek.CompletionsRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.MessageDTO{
			{
				Role:    "system",
				Content: s.currentPrompt,
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		Stream: false,
	})
	if err != nil {
		return err
	}
	if len(resp.Choices) == 0 {
		log.Printf("[error] no deepseek choices\n")
		return errors.New("no answer")
	}
	log.Printf("[info] Результат обработки: %s\n", resp.Choices[0].Message.Content)
	contentFields := strings.Split(resp.Choices[0].Message.Content, " | ")
	if len(contentFields) != 2 {
		log.Printf("[error] contentFields != 2\n")
		return errors.New("invalid answer")
	}

	var items []item.ItemCache
	for _, rawItemID := range strings.Split(strings.ReplaceAll(strings.ReplaceAll(contentFields[1], "[", ""), "]", ""), ";") {
		itemID, err := strconv.ParseInt(rawItemID, 10, 64)
		if err != nil {
			continue
		}
		cacheItem, err := s.itemCache.GetItemByID(ctx, itemID)
		if err != nil {
			continue
		}
		if cacheItem == nil {
			continue
		}
		items = append(items, *cacheItem)
	}

	thinkingResRaw, err := json.Marshal(ThinkResult{
		Reflection: strings.ReplaceAll(strings.ReplaceAll(contentFields[0], "[", ""), "]", ""),
		Items:      items,
	})
	if err != nil {
		log.Printf("[error] cannot marshal: %v\n", err)
		return err
	}
	err = s.redis.Set(ctx, fmt.Sprintf(thinkingResultKey, id), thinkingResRaw, nil)
	if err != nil {
		log.Printf("[error] cannot set thinking result to redis: %v\n", err)
		return err
	}
	return nil
}
