package service

import (
	"os"
	"strconv"
)

func envInt(env string, def int) int {
	v, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return def
	}
	return v
}
