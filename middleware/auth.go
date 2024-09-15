package middleware

import (
	"net/http"
)

type APIKeyConfig struct {
	KeyName string
	Keys    map[string]bool
}

func APIKeyAuth(config APIKeyConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get(config.KeyName)
			if key == "" {
				http.Error(w, "Missing API key", http.StatusUnauthorized)
				return
			}

			if _, valid := config.Keys[key]; !valid {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
