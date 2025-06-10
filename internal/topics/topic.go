package topics

import (
	"time"
)

type Topic struct {
	Id         int64     `db:"id" json:"id" redis:"-"`
	PreviewURL string    `db:"preview_url" json:"preview_url" redis:"preview_url"`
	Title      string    `db:"topic_title" json:"title" redis:"topic_title"`
	Content    string    `db:"topic_content" json:"topic_content" redis:"topic_content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at" redis:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" redis:"updated_at"`
}
