package http

import (
	"encoding/json"
	"log"
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

	err := json.NewEncoder(w).Encode(response{
		Status: "ok",
		Data:   payload,
	})
	if err != nil {
		log.Printf("encode error: %v\n", err)
	}
}

func Return400(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
	})
}

func Return401(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
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

func Return409(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
	})
}

func Return500(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(response{
		Status:     "error",
		ErrMessage: pointer.ToString(msg),
	})
}
