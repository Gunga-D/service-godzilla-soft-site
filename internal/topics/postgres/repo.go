package postgres

// this file implements
// 1. topics fetching from postgres
// 2. topics loading to postgres

import (
	"context"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type Repo struct {
	db postgres.TxDatabase
}

func NewRepo(db postgres.TxDatabase) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateTopic(ctx context.Context, topic topics.Topic) (int64, error) {
	q := sq.Insert("public.topics").
		Columns(
			"topic_title",
			"topic_content",
			"created_at",
			"updated_at",
		).
		Values(
			topic.Title,
			topic.Content,
			topic.CreatedAt,
			time.Now())

	query, args, err := q.Suffix(`RETURNING id`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, err
	}

	var id int64
	if err = r.db.Do(ctx).GetContext(ctx, &id, query, args...); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) FetchIds(ctx context.Context, limit uint64, offset uint64) ([]int64, error) {
	query, args, err := sq.Select("id").From(`public.topics`).
		Offset(offset).
		OrderBy("created_at DESC").
		Limit(limit).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var res []int64
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repo) FetchTopics(ctx context.Context, limit uint64, offset uint64) ([]topics.Topic, error) {
	query, args, err := sq.Select("*").From(`public.topics`).
		Offset(offset).
		OrderBy("created_at DESC").
		Limit(limit).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var res []topics.Topic
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repo) FetchTopic(ctx context.Context, id int64) (topics.Topic, error) {
	query, args, err := sq.Select("*").From(`public.topics`).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return topics.Topic{}, err
	}

	var res topics.Topic
	err = r.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		return topics.Topic{}, err
	}

	return res, nil
}

func (r *Repo) FetchTopicPreview(ctx context.Context, id int64) (topics.Preview, error) {
	topic, err := r.FetchTopic(ctx, id)
	if err != nil {
		return topics.Preview{}, err
	}
	return topics.Preview{
		ImageURL:  "",
		Title:     topic.Title,
		CreatedAt: topic.CreatedAt,
	}, err
}
