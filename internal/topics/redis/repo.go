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
	topicCacheIdsKey = "topics:created_at"
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

func (r *Repo) EmplaceTopics(ctx context.Context, topics []topics.Topic) error {
	keys, values := toMultiSetArgs(topics)
	err := r.redis.MultiSet(ctx, keys, values, nil)
	if err != nil {
		return err
	}

	// set ids to ordered set
	return r.redis.Execute(ctx, func(conn redigo.Conn) error {
		_, err := conn.Do("ZADD", toZAddArgs(topics)...)
		return err
	})
}

func (r *Repo) FetchTopics(ctx context.Context, limit uint64, offset uint64) ([]topics.Topic, error) {
	var res []topics.Topic
	er := r.redis.Execute(ctx, func(conn redigo.Conn) error {
		var err error
		res, err = r.fetchTopicsImpl(conn, limit, offset)
		return err
	})
	return res, er
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

func (r *Repo) addTopic(ctx context.Context, topic topics.Topic) error {
	return r.redis.Execute(ctx, func(conn redigo.Conn) error {
		return r.addTopicImpl(conn, topic)
	})
}

func makeTopicKey(id int64) string {
	return fmt.Sprintf(topicCacheKey, id)
}

func (r *Repo) Clear(ctx context.Context) error {
	return r.redis.Execute(ctx, func(conn redigo.Conn) error {
		return r.removeAllTopicsImpl(conn)
	})
}

func (r *Repo) FetchAllTopics(ctx context.Context) ([]topics.Topic, error) {
	var ids []int64
	err := r.redis.Execute(ctx, func(conn redigo.Conn) error {
		var err error
		ids, err = redigo.Int64s(conn.Do("ZRANGE", "topics:created_at", 0, -1))
		return err
	})
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = makeTopicKey(id)
	}

	result, err := r.redis.MultiGet(ctx, toMultiGetStrings(ids))
	if err != nil {
		return nil, err
	}

	var res []topics.Topic
	for _, data := range result {
		var t topics.Topic
		err := json.Unmarshal(data.([]byte), &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (r *Repo) addTopicImpl(c redigo.Conn, topic topics.Topic) error {
	bytes, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	if _, err = c.Do("SET", makeTopicKey(topic.Id), bytes); err != nil {
		return err
	}

	_, err = c.Do("ZADD", toZAddArgs([]topics.Topic{topic})...)
	return err
}

func (r *Repo) removeTopicImpl(c redigo.Conn, id int64) error {
	_, err := c.Do("DEL", makeTopicKey(id))
	if err != nil {
		return err
	}

	args := redigo.Args{"topics:created_at"}.Add(id)
	_, err = c.Do("ZREM", args...)
	return err
}

func (r *Repo) removeAllTopicsImpl(c redigo.Conn) error {
	ids, err := redigo.Int64s(c.Do("ZRANGE", "topics:created_at", 0, -1))
	if err != nil {
		return err
	}

	// clear set
	_, err = c.Do("DEL", "topics:created_at")
	if err != nil {
		return err
	}

	// if keys doesn't present, return
	if len(ids) == 0 {
		return nil
	}

	// 2. Remove topics:%d cache
	_, err = c.Do("DEL", toMultiGetArgs(ids)...)
	return err
}

func (r *Repo) fetchTopicsImpl(conn redigo.Conn, limit uint64, offset uint64) ([]topics.Topic, error) {
	// 1. Get IDs from sorted set
	ids, err := redigo.Int64s(conn.Do("ZRANGE", "topics:created_at", offset, offset+limit-1))
	if err != nil {
		return nil, fmt.Errorf("failed to get topic IDs: %w", err)
	}

	if len(ids) == 0 {
		return []topics.Topic{}, nil
	}

	values, err := redigo.Values(conn.Do("MGET", toMultiGetArgs(ids)...))
	if err != nil {
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	// 4. Parse the results
	var result []topics.Topic
	for _, val := range values {
		var t topics.Topic
		buf, ok := val.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected value type: %T", val)
		}

		if err := json.Unmarshal(buf, &t); err != nil {
			return nil, fmt.Errorf("failed to unmarshal t: %w", err)
		}

		result = append(result, t)
	}

	return result, nil
}

func toZAddArgs(topics []topics.Topic) []interface{} {
	args := redigo.Args{topicCacheIdsKey}
	for _, topic := range topics {
		args = args.Add(topic.CreatedAt.Unix(), topic.Id)
	}
	return args
}

func toMultiSetArgs(topics []topics.Topic) ([]string, []interface{}) {
	keys := make([]string, len(topics))
	values := make([]interface{}, len(topics))
	for i, t := range topics {
		bytes, err := json.Marshal(t)
		if err != nil {
			return nil, nil
		}
		keys[i] = makeTopicKey(t.Id)
		values[i] = bytes
	}
	return keys, values
}
func toMultiGetArgs(ids []int64) []interface{} {
	var res = make([]interface{}, len(ids))
	strings := toMultiGetStrings(ids)
	for i, str := range strings {
		res[i] = str
	}
	return res
}

func toMultiGetStrings(ids []int64) []string {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = makeTopicKey(id)
	}
	return keys
}
