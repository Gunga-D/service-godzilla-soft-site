package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	topicCacheKey = "topic:%d"
)

type Repo struct {
	redis redis.Redis
}

func NewRedisRepo(redis redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}

func (r *Repo) CreateTopic(ctx context.Context, topic topics.Topic, id int64) (string, error) {
	key := fmt.Sprintf(topicCacheKey, id)
	return key, r.redis.Set(ctx, key, topic, nil)
}

func (r *Repo) FetchTopic(ctx context.Context, key string) (topics.Topic, error) {
	bytes, err := redigo.Bytes(r.redis.Get(ctx, key))
	if err != nil {
		return topics.Topic{}, err
	}

	var t topics.Topic
	if err = json.Unmarshal(bytes, &t); err != nil {
		return topics.Topic{}, err
	}

	return t, nil
}

func (r *Repo) FetchTopicPreview(ctx context.Context, key string) (topics.Preview, error) {
	topic, err := r.FetchTopic(ctx, key)
	if err != nil {
		return topics.Preview{}, err
	}
	return topics.Preview{
		ImageURL:  "",
		Title:     topic.Title,
		CreatedAt: topic.CreatedAt,
	}, nil
}
