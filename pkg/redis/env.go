package redis

import (
	"fmt"
	"os"
	"strconv"
)

const (
	_envHost = "REDIS_HOST"
	_envPort = "REDIS_PORT"
	_envPwd  = "REDIS_PASSWORD"
)

func loadHost() (string, error) {
	if host := os.Getenv(_envHost); host != "" {
		return host, nil
	}
	return "", fmt.Errorf("no host")
}

func loadPort() (int, error) {
	port, err := strconv.Atoi(os.Getenv(_envPort))
	if err != nil {
		return 0, fmt.Errorf("no port")
	}
	return port, nil
}

func loadPwd() (string, error) {
	if pwd := os.Getenv(_envPwd); pwd != "" {
		return pwd, nil
	}
	return "", fmt.Errorf("no pwd")
}
