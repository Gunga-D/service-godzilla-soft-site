package telegram_sign_in

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type jwtService interface {
	GenerateToken(userID int64, email *string) (string, error)
}

type telegramRegistrationDatabus interface {
	PublishDatabusTelegramRegistration(ctx context.Context, msg databus.TelegramRegistrationDTO) error
}
