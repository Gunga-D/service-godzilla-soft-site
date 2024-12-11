package password

import (
	"context"
	"crypto/sha256"
	"fmt"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (u *Generator) GeneratePassword(_ context.Context, pwd string) string {
	return hash(pwd)
}

func (u *Generator) ValidatePassword(ctx context.Context, pwd string, checkedPwd string) bool {
	return pwd == hash(checkedPwd)
}

func hash(in string) string {
	h := sha256.New()
	h.Write([]byte(in))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
