package middlewares

import (
	"log"
	"net/http"
)

func CorsMiddleware(allowedOrigins ...string) func(http.Handler) http.Handler {
	originSet := make(map[string]struct{})
	for _, origin := range allowedOrigins {
		originSet[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if _, exists := originSet[origin]; !exists {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			log.Println("Origin header:", origin)
			log.Println("Allowed origins:", originSet)

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
