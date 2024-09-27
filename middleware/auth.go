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
			if isValidAPIKey(r, config) {
				next.ServeHTTP(w, r)
			} else {
				respondUnauthorized(w)
			}
		})
	}
}

func isValidAPIKey(r *http.Request, config APIKeyConfig) bool {
	key := r.Header.Get(config.KeyName)
	return key != "" && config.Keys[key]
}

func respondUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
