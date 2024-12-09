package http

import (
	"encoding/json"
	"net/http"

	"github.com/AlekSi/pointer"
)

type response struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	ErrMessage *string     `json:"message,omitempty"`
}

func ReturnOK(payload interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response{
		Status: "ok",
		Data:   payload,
	})
}

func Return404(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
	})
}

func Return500(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
	})
}
