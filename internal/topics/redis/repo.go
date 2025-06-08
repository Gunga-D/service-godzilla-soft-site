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

type OutOfBounds string

func (e OutOfBounds) Error() string {
	return string(e)
}

const OutOfBoundsError OutOfBounds = "[limit:offset] is out of bounds"

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

func (r *Repo) FetchTopics(ctx context.Context, limit uint64, offset uint64) ([]topics.Topic, error) {
	ids, err := r.fetchIds(ctx)
	if err != nil || len(ids) == 0 {
		return nil, err
	}

	first := offset
	if first > uint64(len(ids)-1) {
		return nil, OutOfBoundsError
	}

	last := limit + offset
	if last > uint64(len(ids)-1) {
		return nil, OutOfBoundsError
	}

	var res []topics.Topic
	for _, id := range ids[first:last] {
		topic, err := r.GetTopic(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, *topic)
	}
	return res, nil
}

func (r *Repo) GetTopic(ctx context.Context, id int64) (*topics.Topic, error) {
	bytes, err := redigo.Bytes(r.redis.Get(ctx, makeTopicKey(id)))
	if err != nil || len(bytes) == 0 {
		return nil, err
	}

	var t topics.Topic
	if err = json.Unmarshal(bytes, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repo) fetchIds(ctx context.Context) ([]int64, error) {
	res, err := redigo.Int64s(r.redis.Members(ctx, topicCacheIdsKey))
	if err != nil {
		return nil, err
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

func makeTopicKey(id int64) string {
	return fmt.Sprintf(topicCacheKey, id)
}
