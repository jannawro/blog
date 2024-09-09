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
	sessionConfig := NewSessionConfig(store)

	tests := []struct {
		name           string
		authenticated  bool
		expectedStatus int
		envSessionName string
	}{
		{"Authenticated Session", true, http.StatusOK, "test-session"},
		{"Unauthenticated Session", false, http.StatusUnauthorized, "test-session"},
		{"Authenticated Session with Default Name", true, http.StatusOK, ""},
		{"Unauthenticated Session with Default Name", false, http.StatusUnauthorized, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset the SESSION_NAME environment variable
			if tt.envSessionName != "" {
				t.Setenv("SESSION_NAME", tt.envSessionName)
			} else {
				t.Setenv("SESSION_NAME", "")
			}

			handler := SessionAuth(sessionConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			// Determine the session name
			sessionName := tt.envSessionName
			if sessionName == "" {
				sessionName = "default_session_name"
			}

			// Create and save a session with the authentication status
			session, err := store.New(req, sessionName)
			require.NoError(t, err, "Failed to create new session")
			session.Values["authenticated"] = tt.authenticated
			err = session.Save(req, rr)
			require.NoError(t, err, "Failed to save session")

			// Copy the session cookie from the recorder to the request
			req.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))

			// Reset the recorder for the actual request
			rr = httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
		})
	}
}
