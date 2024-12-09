package service

import (
	"encoding/json"
	"net/http"
)

const (
	_statusOK    = "ok"
	_statusError = "error"
)

type healthInfo struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		info := healthInfo{
			Status: _statusOK,
		}
		body, _ := json.MarshalIndent(&info, "", " ")
		w.Write(body)
	})
}
