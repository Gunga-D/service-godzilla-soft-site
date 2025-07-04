package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/code"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/order"
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

func (r *repo) CreateItemOrder(ctx context.Context, email string, amount int64, itemID int64, itemSlip string, itemName string, utm *string) (string, error) {
	var orderID string
	err := r.db.WithTx(ctx, func(ctx context.Context) error {
		codeValue, err := r.blockCode(ctx, itemID)
		if err != nil {
			return fmt.Errorf("block code: %v", err)
		}

		newOrderID, err := r.insertOrder(ctx, email, amount, codeValue, itemSlip, itemName, utm)
		if err != nil {
			return fmt.Errorf("insert order: %v", err)
		}
		orderID = newOrderID
		return nil
	})
	if err != nil {
		return "", err
	}
	return orderID, nil
}

func (r *repo) CreateItemGiftOrder(ctx context.Context, steamProfile string, amount int64, itemID int64, utm *string) (string, error) {
	newOrderID, err := r.insertOrder(ctx, steamProfile, amount, fmt.Sprintf("STEAM_GIFT_%d", itemID), "Нет инструкции", "Нет названия", utm)
	if err != nil {
		return "", fmt.Errorf("insert order: %v", err)
	}
	return newOrderID, nil
}

func (r *repo) CreateSteamOrder(ctx context.Context, steamLogin string, amount int64) (string, error) {
	newOrderID, err := r.insertOrder(ctx, steamLogin, amount, "STEAM_INVOICE_BY_LOGIN", "Нет инструкции", "Нет названия", nil)
	if err != nil {
		return "", fmt.Errorf("insert order: %v", err)
	}
	return newOrderID, nil
}

func (r *repo) PaidOrder(ctx context.Context, orderID string) error {
	query, args, err := sq.Update("public.order").
		Where(sq.Eq{"id": orderID}).
		Set("status", order.PaidStatus).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.Do(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *repo) FinishOrder(ctx context.Context, orderID string) error {
	return r.db.WithTx(ctx, func(ctx context.Context) error {
		codeValue, err := r.setOrderStatus(ctx, orderID, order.FinishedStatus)
		if err != nil {
			return err
		}

		return r.setCodeStatus(ctx, codeValue, code.DeliveredStatus)
	})
}

func (r *repo) FailedOrder(ctx context.Context, orderID string) error {
	return r.db.WithTx(ctx, func(ctx context.Context) error {
		codeValue, err := r.setOrderStatus(ctx, orderID, order.FailedStatus)
		if err != nil {
			return err
		}

		return r.setCodeStatus(ctx, codeValue, code.FreeStatus)
	})
}

func (r *repo) FetchPaidOrders(ctx context.Context) ([]order.PaidOrder, error) {
	query, args, err := sq.Select("id, email, code_value, item_slip", "amount").From(`public.order`).
		Where(sq.Eq{"status": order.PaidStatus}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []order.PaidOrder
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) FetchUserOrdersByEmail(ctx context.Context, email string) ([]order.UserOrder, error) {
	query, args, err := sq.Select("id, code_value, item_slip", "item_name", "amount", "status").From(`public.order`).
		Where(sq.And{
			sq.Or{
				sq.Eq{"status": order.FinishedStatus},
				sq.Eq{"status": order.PaidStatus},
			},
			sq.Eq{"email": email},
		}).
		OrderBy("created_at desc").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []order.UserOrder
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) setCodeStatus(ctx context.Context, codeValue string, status string) error {
	query, args, err := sq.Update("public.code").
		Where(sq.Eq{"value": codeValue}).
		Set("status", status).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.Do(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *repo) setOrderStatus(ctx context.Context, orderID string, status string) (string, error) {
	q := sq.Update("public.order").
		Where(sq.Eq{"id": orderID}).
		Set("status", status).
		Set("updated_at", time.Now())

	query, args, err := q.
		Suffix(`
RETURNING code_value
`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var codeValue string
	if err := r.db.Do(ctx).GetContext(ctx, &codeValue, query, args...); err != nil {
		return "", err
	}
	return codeValue, nil
}

func (r *repo) insertOrder(ctx context.Context, email string, amount int64, codeValue string, itemSlip string, itemName string, utm *string) (string, error) {
	q := sq.Insert("public.order").
		Columns(
			"email",
			"code_value",
			"item_slip",
			"item_name",
			"utm",
			"amount",
			"status",
			"created_at",
			"updated_at",
		).Values(
		email,
		codeValue,
		itemSlip,
		itemName,
		utm,
		amount,
		order.PendingStatus,
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
		return "", err
	}

	var id string
	if err := r.db.Do(ctx).GetContext(ctx, &id, query, args...); err != nil {
		return "", err
	}
	return id, nil
}

func (r *repo) blockCode(ctx context.Context, itemID int64) (string, error) {
	var codeValue string
	if err := r.db.Do(ctx).GetContext(ctx, &codeValue, `
		update public.code set status=$1, updated_at=$2 where
		value=(select value from public.code where status=$3 and item_id=$4 limit 1)
		returning value
	`, code.BlockedStatus, time.Now(), code.FreeStatus, itemID); err != nil {
		return "", err
	}
	return codeValue, nil
}
