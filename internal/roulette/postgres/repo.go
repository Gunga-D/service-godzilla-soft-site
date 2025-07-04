package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"log"
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

func (r *Repo) CreateSession(ctx context.Context) (uuid.UUID, error) {
	s, err := r.newSession(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	return s.Id, nil
}

func (r *Repo) FetchItems(ctx context.Context, sessionId uuid.UUID) ([]roulette.SessionItem, error) {
	query, args, err := sq.Select("*").
		From("public.roulette_session_items").
		Where(sq.Eq{"session_id": sessionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []roulette.SessionItem
	err = r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

func (r *Repo) AddItemsToRoulette(ctx context.Context, ids []int64) error {
	query, args, err := sq.Select("*").From("public.roulette").
		Where(sq.Eq{"item.id": ids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	var existingItems []item.Item
	err = r.db.SelectContext(ctx, &existingItems, query, args...)
	if err != nil {
		return err
	}
	if len(ids)-len(existingItems) > 0 {
		log.Printf("unable to add %d existingItems to roulette", len(ids)-len(existingItems))
	}
	if len(existingItems) == 0 {
		return errors.New("no items to add")
	}

	q := sq.Insert("public.roulette_random_item").
		Columns("id", "total_cost", "item_category")

	for _, it := range existingItems {
		q = q.Values(
			it.ID,
			it.CurrentPrice,
			roulette.CategoryFromPrice(it.CurrentPrice))
	}

	query, args, err = q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) FetchAvailableItems(ctx context.Context) ([]int64, error) {
	query, args, err := sq.Select("id").
		From("public.roulette_random_item").
		ToSql()

	if err != nil {
		return nil, err
	}

	var ids []int64
	err = r.db.SelectContext(ctx, &ids, query, args...)
	return ids, err
}

func (r *Repo) AddTopItemsToSession(ctx context.Context, sId uuid.UUID, itemsIds []int64) error {
	q := sq.Insert("public.roulette_session_items").
		Columns("session_id", "item_id", "is_top")

	for _, id := range itemsIds {
		q = q.Values(sId, id, true)
	}

	query, args, err := q.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) GetSessionById(ctx context.Context, sessionId uuid.UUID) (*roulette.Session, error) {
	sql, args, err := sq.Select("*").
		From("public.roulette_session").
		Where(sq.Eq{"id": sessionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var s roulette.Session
	err = r.db.GetContext(ctx, &s, sql, args...)
	return &s, err
}

func (r *Repo) newSession(ctx context.Context) (roulette.Session, error) {
	q := sq.Insert("public.roulette_session").
		Columns(
			"created_at",
			"updated_at").
		Values(
			time.Now(),
			time.Now(),
		)
	query, args, err := q.Suffix(`RETURNING *`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		_ = fmt.Errorf(`failed to build sql query: %v`, err)
	}

	var s roulette.Session
	err = r.db.GetContext(ctx, &s, query, args...)
	if err != nil {
		return roulette.Session{}, err
	}

	return s, nil
}

type AverageInfo struct {
	avg   int64 `db:"avg"`
	count int64 `db:"count"`
}

func (r *Repo) getAverageInfo(ctx context.Context, sId uuid.UUID) (AverageInfo, error) {
	sql, args, err := sq.Select("avg(total_cost), count(total_cost)").
		From("public.roulette_session_items si").
		Join("public.roulette_random_item ri ON si.item_id = ri.id").
		Where(sq.Eq{"session_id": sId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return AverageInfo{}, err
	}
	var avgCost AverageInfo
	err = r.db.GetContext(ctx, &avgCost, sql, args...)
	return avgCost, err
}
