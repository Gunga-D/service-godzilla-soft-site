package neuro

import (
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type Task struct {
	ID        string    `db:"id"`
	Query     string    `db:"query"`
	Result    string    `db:"result"`
	CreatedAt time.Time `db:"created_at"`
}

type TaskResult struct {
	Success bool            `json:"success"`
	Message *string         `json:"message"`
	Data    *TaskResultData `json:"data"`
}

type TaskResultData struct {
	Raw        string           `json:"raw"`
	Reflection string           `json:"reflection"`
	Items      []item.ItemCache `json:"items"`
}
