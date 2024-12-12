package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
)

func GeneratePassword(_ context.Context, pwd string) string {
	return hash(pwd)
}

func ValidatePassword(ctx context.Context, pwd string, checkedPwd string) bool {
	return pwd == hash(checkedPwd)
}

func hash(in string) string {
	h := sha256.New()
	h.Write([]byte(in))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
