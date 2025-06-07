package cached

import (
	"context"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/redis"
)

type Repo struct {
	pg    *postgres.Repo
	redis *redis.Repo
}

func NewRepo(pg *postgres.Repo, redis *redis.Repo) *Repo {
	return &Repo{
		pg,
		redis}
}

func (r *Repo) CreateTopic(ctx context.Context, topic topics.Topic) (int64, error) {
	// assign id provided by postgres
	id, err := r.pg.CreateTopic(ctx, topic)
	topic.Id = id
	if err != nil {
		return -1, err
	}
	return topic.Id, r.redis.CreateTopic(ctx, topic)
}

func (r *Repo) FetchTopics(ctx context.Context, limit uint64, offset uint64) ([]topics.Topic, error) {
	fetched, err := r.redis.FetchTopics(ctx, limit, offset)
	if err == nil {
		return fetched, nil
	}
	return r.pg.FetchTopics(ctx, limit, offset)
}

func (r *Repo) GetTopic(ctx context.Context, id int64) (*topics.Topic, error) {
	topic, err := r.redis.GetTopic(ctx, id)
	if err == nil {
		return topic, nil
	}
	return r.pg.GetTopic(ctx, id)
}
