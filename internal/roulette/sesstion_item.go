package roulette

import (
	"github.com/google/uuid"
	"time"
)

type SessionItem struct {
	SessionId uuid.UUID  `json:"session_id" db:"session_id"`
	ItemId    int        `json:"item_id" db:"item_id"`
	IsTop     bool       `json:"is_top" db:"is_top"`
	WonAt     *time.Time `json:"won_at" db:"won_at"`
}

type SessionItemDTO struct {
	ItemId int        `json:"item_id" db:"item_id"`
	IsTop  bool       `json:"is_top" db:"is_top"`
	WonAt  *time.Time `json:"won_at" db:"won_at"`
}

func (i SessionItem) DTO() SessionItemDTO {
	return SessionItemDTO{
		ItemId: i.ItemId,
		IsTop:  i.IsTop,
		WonAt:  i.WonAt,
	}
}
