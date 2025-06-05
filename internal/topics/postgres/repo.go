package postgres

import (
	"context"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type Repo struct {
	db postgres.TxDatabase
}

type Topic struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRepo(db postgres.TxDatabase) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateTopic(ctx context.Context, topic Topic) (int64, error) {
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
