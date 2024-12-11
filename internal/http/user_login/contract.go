package user_login

import "context"

type pwdValidator interface {
	ValidatePassword(ctx context.Context, pwd string, checkedPwd string) bool
}
