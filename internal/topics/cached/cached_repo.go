package cached

import (
	"context"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics/redis"
)

type Repo struct {
	pg     *postgres.Repo
	redis  *redis.Repo
	keyMap map[int64]string
}

func NewRepo(pg *postgres.Repo, redis *redis.Repo) *Repo {
	return &Repo{
		pg,
		redis,
		make(map[int64]string)}
}

func (r *Repo) CreateTopic(ctx context.Context, topic topics.Topic) error {
	id, err := r.pg.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	key, err := r.redis.CreateTopic(ctx, topic, id)
	if err != nil {
		return err
	}
	r.keyMap[id] = key
	return nil
}

func (r *Repo) GetPreviews(ctx context.Context, limit uint64, offset uint64) ([]topics.Preview, error) {
	ids, err := r.pg.FetchIds(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var previews []topics.Preview
	for _, id := range ids {
		key, contains := r.keyMap[int64(id)]
		var preview topics.Preview
		if contains {
			preview, err = r.redis.FetchTopicPreview(ctx, key)
			if err != nil {
				return nil, err
			}
		} else {
			// fetch from postgres
			topic, err := r.pg.FetchTopic(ctx, int64(id))
			if err != nil {
				return nil, err
			}
			preview = topics.Preview{
				Title:     topic.Title,
				CreatedAt: topic.CreatedAt,
				ImageURL:  "",
			}
			// cache in redis
			key, err := r.redis.CreateTopic(ctx, topic, int64(id))
			if err != nil {
				return nil, err
			}
			r.keyMap[int64(id)] = key
		}
		previews = append(previews, preview)
	}
	return previews, nil
}
