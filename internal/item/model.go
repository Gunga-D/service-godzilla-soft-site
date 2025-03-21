package item

import (
	"time"

	"github.com/lib/pq"
)

const (
	ActiveStatus   = "active"
	PausedStatus   = "paused"
	ArchivedStatus = "archived"
)

type Item struct {
	ID          int64   `db:"id"`
	Title       string  `db:"title"`
	Description *string `db:"description"`
	CategoryID  int64   `db:"category_id"`
	Platform    string  `db:"platform"`
	Region      string  `db:"region"`
	Publisher   *string `db:"publisher"`
	Creator     *string `db:"creator"`
	ReleaseDate *string `db:"release_date"`
	// Указывается с копейками, таким образом:
	// 100 рублей = 10000
	CurrentPrice  int64          `db:"current_price"`
	IsForSale     bool           `db:"is_for_sale"`
	OldPrice      *int64         `db:"old_price"`
	LimitPrice    *int64         `db:"limit_price"`
	ThumbnailURL  string         `db:"thumbnail_url"`
	BackgroundURL *string        `db:"background_url"`
	BxImageURL    *string        `db:"bx_image_url"`
	BxGalleryUrls pq.StringArray `db:"bx_gallery_urls"`
	Status        string         `db:"status"`
	Slip          string         `db:"slip"`
	YandexID      *string        `db:"yandex_id"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}
