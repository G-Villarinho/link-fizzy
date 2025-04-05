package middlewares

import (
	"net/http"
	"strings"
)

var allowedOrigins = map[string]bool{
	"http://localhost:3001": true,
	"http://localhost:5175": true}

var allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
var allowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
var exposedHeaders = []string{"Content-Length", "X-Kuma-Revision"}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
