package subscribe

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type jwtService interface {
	GenerateToken(userID int64, email *string) (string, error)
}

type subRepo interface {
	CreateSubscriptionBill(ctx context.Context, userID int64, amount int64, forPeriod time.Duration) (string, error)
}

type subChecker interface {
	HasSubscription(ctx context.Context, userID int64) (bool, error)
}

type sendToEmailDatabus interface {
	PublishDatabusSendToEmail(ctx context.Context, msg databus.SendToEmailDTO) error
}
