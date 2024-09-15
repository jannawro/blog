package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
