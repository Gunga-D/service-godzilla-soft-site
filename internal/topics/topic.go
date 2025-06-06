package topics

import (
	"database/sql"
	"time"
)

type Preview struct {
	ImageURL  string    `db:"image_url" json:"image"`
	Title     string    `db:"topic_title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Topic struct {
	Id         int64          `db:"id" json:"id"`
	PreviewURL sql.NullString `db:"preview_url" json:"preview_url"`
	Title      string         `db:"topic_title" json:"title"`
	Content    string         `db:"topic_content" json:"topic_content"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}
