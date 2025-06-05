package postgres

// this file implements
// 1. topics fetchig from postgres
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
