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
			if key == "" || !config.Keys[key] {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
