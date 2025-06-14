package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
)

type repo struct {
	db postgres.TxDatabase
}

func NewRepo(db postgres.TxDatabase) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateItem(ctx context.Context, i item.Item) (int64, error) {
	q := sq.Insert("public.item").
		Columns(
			"title",
			"description",
			"category_id",
			"platform",
			"region",
			"current_price",
			"is_for_sale",
			"old_price",
			"limit_price",
			"thumbnail_url",
			"background_url",
			"status",
			"slip",
			"is_steam_gift",
			"yandex_id",
			"steam_app_id",
			"price_loc",
			"popular",
			"new",
			"unavailable",
			"horizontal_image",
			"vertical_image",
			"steam_raw_data",
			"exists_in_random",
			"created_at",
			"updated_at",
		).Values(
		i.Title,
		i.Description,
		i.CategoryID,
		i.Platform,
		i.Region,
		i.CurrentPrice,
		i.IsForSale,
		i.OldPrice,
		i.LimitPrice,
		i.ThumbnailURL,
		i.BackgroundURL,
		i.Status,
		i.Slip,
		i.IsSteamGift,
		i.YandexID,
		i.SteamAppID,
		i.PriceLoc,
		i.Popular,
		i.New,
		i.Unavailable,
		i.HorizontalImage,
		i.VerticalImage,
		i.SteamRawData,
		i.ExistsInRandom,
		time.Now(),
		time.Now(),
	)
	query, args, err := q.
		Suffix(`
			RETURNING id
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := r.db.GetContext(ctx, &id, query, args...); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) ChangeItemState(ctx context.Context, itemID int64, status string) error {
	query, args, err := sq.Update("public.item").
		Where(sq.Eq{"id": itemID}).
		Set("status", status).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetItemByID(ctx context.Context, id int64) (*item.Item, error) {
	query, args, err := sq.Select("*").From(`public.item`).
		Where(sq.Eq{"id": id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []item.Item
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) GetItemBySteamAppID(ctx context.Context, steamAppID int64) (*item.Item, error) {
	query, args, err := sq.Select("*").From(`public.item`).
		Where(sq.Eq{"steam_app_id": steamAppID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []item.Item
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) FetchItemsByFilter(ctx context.Context, criteria sq.And, limit uint64, offset uint64, orderBy []string) ([]item.Item, error) {
	query, args, err := sq.Select("*").From(`public.item`).
		Where(criteria).
		OrderBy(orderBy...).
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []item.Item
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) GetItemsCountByFilter(ctx context.Context, criteria sq.And) (int64, error) {
	query, args, err := sq.Select("count(*)").From(`public.item`).
		Where(criteria).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	var res int64
	err = r.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *repo) FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error) {
	query, args, err := sq.Select("*").From(`public.item`).
		Where(sq.And{
			sq.Gt{"id": cursor},
			sq.Eq{"status": item.ActiveStatus},
		}).
		OrderBy("id").
		Limit(limit).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var res []item.Item
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) UpdatePrice(ctx context.Context, itemID int64, currentPrice int64, limitPrice int64, priceLoc string) error {
	q := sq.Update("item").
		Where(sq.Eq{"id": itemID}).
		Set("current_price", currentPrice).
		Set("limit_price", limitPrice).
		Set("price_loc", priceLoc)

	query, args, err := q.
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *repo) UpdateSteamRawData(ctx context.Context, itemID int64, steamRawData string) error {
	q := sq.Update("item").
		Where(sq.Eq{"id": itemID}).
		Set("steam_raw_data", steamRawData)

	query, args, err := q.
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
