package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type SessionConfig struct {
	CookieStore *sessions.CookieStore
}

func NewSessionConfig(cookieStore *sessions.CookieStore) *SessionConfig {
	return &SessionConfig{
		CookieStore: cookieStore,
	}
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

func SessionAuth(config *SessionConfig) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            sessionName := os.Getenv("SESSION_NAME")
            if sessionName == "" {
                sessionName = "default_session_name" // fallback
            }

            session, err := config.CookieStore.Get(r, sessionName)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Check if the session is authenticated
            if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // If we reach here, the session is authenticated
            next.ServeHTTP(w, r)
        })
    }
}
