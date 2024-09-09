package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var store *sessions.CookieStore

func TestAPIKeyAuth(t *testing.T) {
	config := APIKeyConfig{
		KeyName: "X-API-Key",
		Keys: map[string]bool{
			"valid-key": true,
		},
	}

	tests := []struct {
		name           string
		key            string
		expectedStatus int
	}{
		{"Valid API Key", "valid-key", http.StatusOK},
		{"Invalid API Key", "invalid-key", http.StatusUnauthorized},
		{"Missing API Key", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := APIKeyAuth(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.key != "" {
				req.Header.Set("X-API-Key", tt.key)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
		})
	}
}

func TestSessionAuth(t *testing.T) {
	store = sessions.NewCookieStore([]byte("test-secret"))
	config := SessionConfig{
		SessionName: "test-session",
		CookieStore: store,
	}

	tests := []struct {
		name           string
		authenticated  bool
		expectedStatus int
	}{
		{"Authenticated Session", true, http.StatusOK},
		{"Unauthenticated Session", false, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := SessionAuth(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			// Create a session and set the authenticated value
			session, err := store.New(req, config.SessionName)
			require.NoError(t, err, "Failed to create new session")

			session.Values["authenticated"] = tt.authenticated
			err = session.Save(req, rr)
			require.NoError(t, err, "Failed to save session")

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
		})
	}
}

func TestCombinedAuth(t *testing.T) {
	store = sessions.NewCookieStore([]byte("test-secret"))
	apiConfig := APIKeyConfig{
		KeyName: "X-API-Key",
		Keys: map[string]bool{
			"valid-key": true,
		},
	}

	sessionConfig := SessionConfig{
		SessionName: "test-session",
		CookieStore: store,
	}

	tests := []struct {
		name           string
		apiKey         string
		authenticated  bool
		expectedStatus int
	}{
		{"Valid API Key", "valid-key", false, http.StatusOK},
		{"Invalid API Key, Authenticated Session", "invalid-key", true, http.StatusOK},
		{"Invalid API Key, Unauthenticated Session", "invalid-key", false, http.StatusUnauthorized},
		{"No API Key, Authenticated Session", "", true, http.StatusOK},
		{"No API Key, Unauthenticated Session", "", false, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CombinedAuth(apiConfig, sessionConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rr := httptest.NewRecorder()

			// Create a session and set the authenticated value
			session, err := store.New(req, sessionConfig.SessionName)
			require.NoError(t, err, "Failed to create new session")

			session.Values["authenticated"] = tt.authenticated
			err = session.Save(req, rr)
			require.NoError(t, err, "Failed to save session")

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
		})
	}
}
