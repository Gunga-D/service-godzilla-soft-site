package roulette_session

import (
	"github.com/google/uuid"
)

type CreateSessionResponse struct {
	SessionId  uuid.UUID `json:"session_id"`
	PaymentUrl string    `json:"payment_url"`
}

type AddTopItemsRequest struct {
	ItemIds []int64 `json:"items"`
}
