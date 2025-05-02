package collection

import "time"

type Collection struct {
	ID              int64      `db:"id"`
	CategoryID      int64      `db:"category_id"`
	Name            string     `db:"name"`
	Description     string     `db:"description"`
	BackgroundImage string     `db:"background_image"`
	HeaderImage     *string    `db:"header_image"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

type CollectionItem struct {
	ID           int64 `db:"id"`
	CollectionID int64 `db:"collection_id"`
	ItemID       int64 `db:"item_id"`
}
