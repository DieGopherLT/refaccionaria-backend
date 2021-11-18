package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func WriteJsonResponse(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&value)
}
