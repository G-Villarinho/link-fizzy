package responses

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		jsoniter.NewEncoder(w).Encode(data)
	}
}

func NoContent(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
