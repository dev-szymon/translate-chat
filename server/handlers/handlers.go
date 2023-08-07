package handlers

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	err    string
	status int
}

func (e *apiError) Error() string {
	return e.err
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
