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
	topicCacheIdsKey = "topics"
	topicCacheKey    = "topic:%d"
)

type Repo struct {
	redis redis.Redis
}

func NewRedisRepo(redis redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}

func (r *Repo) CreateTopic(ctx context.Context, topic topics.Topic) error {
	return r.addTopic(ctx, topic)
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

func (r *Repo) fetchIds(ctx context.Context) ([]int64, error) {
	members, err := r.redis.Members(ctx, topicCacheIdsKey)
	if err != nil {
		return nil, err
	}

	res := make([]int64, len(members))
	for i, val := range members {
		res[i] = val.(int64)
	}

	return res, nil
}

func (r *Repo) addTopic(ctx context.Context, topic topics.Topic) error {
	var arr []interface{}
	arr = append(arr, topic.Id)
	err := r.redis.Add(ctx, topicCacheIdsKey, arr)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, makeTopicKey(topic.Id), topic, nil)
}
