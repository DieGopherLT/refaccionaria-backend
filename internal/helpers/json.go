package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func WriteJsonMessage(w http.ResponseWriter, code int, response Response) {
	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func WriteJsonResponse(w http.ResponseWriter, code int, value interface{}) {
	js, _ := json.Marshal(value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}
