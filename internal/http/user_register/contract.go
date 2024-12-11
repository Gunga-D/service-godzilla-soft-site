package user_register

import "context"

type pwdGenerator interface {
	GeneratePassword(_ context.Context, pwd string) string
}
