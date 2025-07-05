package healthcheck

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type Status string

const (
	StatusOK          Status = "OK"
	StatusUnavailable Status = "Unavailable"
)

type Health struct {
	Status     Status            `json:"status"`
	Failures   map[string]string `json:"failures,omitempty"`
	Components []string          `json:"components"`
}

func writeResponse(w http.ResponseWriter, body interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
