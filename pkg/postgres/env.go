package postgres

import (
	"fmt"
	"os"
	"strconv"
)

const (
	_envHost   = "POSTGRES_HOST"
	_envPort   = "POSTGRES_PORT"
	_envUser   = "POSTGRES_USER"
	_envPwd    = "POSTGRES_PASSWORD"
	_envDBname = "POSTGRES_DB"
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

func loadUser() (string, error) {
	if user := os.Getenv(_envUser); user != "" {
		return user, nil
	}
	return "", fmt.Errorf("no user")
}

func loadPwd() (string, error) {
	if pwd := os.Getenv(_envPwd); pwd != "" {
		return pwd, nil
	}
	return "", fmt.Errorf("no pwd")
}

func loadDBname() (string, error) {
	if dbname := os.Getenv(_envDBname); dbname != "" {
		return dbname, nil
	}
	return "", fmt.Errorf("no dbname")
}
