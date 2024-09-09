package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("your-secret-key"))
}

type SessionConfig struct {
	SessionName string
}

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

func CombinedAuth(apiConfig APIKeyConfig, sessionConfig SessionConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for API key first
			key := r.Header.Get(apiConfig.KeyName)
			if key != "" {
				if _, valid := apiConfig.Keys[key]; valid {
					next.ServeHTTP(w, r)
					return
				}
			}

			// If no valid API key, check for session
			session, err := store.Get(r, sessionConfig.SessionName)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if auth, ok := session.Values["authenticated"].(bool); ok && auth {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}
