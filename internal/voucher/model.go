package voucher

import "time"

type Voucher struct {
	ID           int64      `db:"id"`
	Type         string     `db:"type"`
	Value        string     `db:"value"`
	Impact       int64      `db:"impact"`
	HasActivated bool       `db:"has_activated"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
