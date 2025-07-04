package roulette

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Id        uuid.UUID     `json:"id" db:"id"`
	Status    SessionStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

type SesssionDTO struct {
	Status    SessionStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

func (s Session) DTO() SesssionDTO {
	return SesssionDTO{
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
