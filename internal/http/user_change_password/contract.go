package user_change_password

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
)

type sendToEmailDatabus interface {
	PublishDatabusSendToEmail(ctx context.Context, msg databus.SendToEmailDTO) error
}
