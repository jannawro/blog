package middleware_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jannawro/blog/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetReqID(t *testing.T) {
	originalReq, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err, "Failed to create request")

	// Call the SetReqID function
	modifiedReq := middleware.SetReqID(originalReq)

	// Check if the returned request is not nil
	assert.NotNil(t, modifiedReq, "SetReqID returned nil request")

	// Get the request ID from the context
	requestID, ok := modifiedReq.Context().Value(middleware.ContextKey{}).(uuid.UUID)
	assert.True(t, ok, "Request ID not found in context or not of type uuid.UUID")

	// Check if the request ID is a valid UUID
	assert.NotEqual(t, uuid.Nil, requestID, "Request ID is nil UUID")

	// Verify that the original request context is different from the modified request context
	assert.NotEqual(t, originalReq.Context(), modifiedReq.Context(), "Request context was not modified")

	// Verify that calling SetReqID again produces a different UUID
	secondModifiedReq := middleware.SetReqID(modifiedReq)
	secondRequestID, _ := secondModifiedReq.Context().Value(middleware.ContextKey{}).(uuid.UUID)
	assert.NotEqual(t, requestID, secondRequestID, "SetReqID did not generate a new UUID on second call")
}

func TestReqIDFromCtx(t *testing.T) {
	t.Run("Successful UUID retrieval", func(t *testing.T) {
		expectedUUID := uuid.New()
		ctx := context.WithValue(context.Background(), middleware.ContextKey{}, expectedUUID)

		result := middleware.ReqIDFromCtx(ctx)
		assert.Equal(t, expectedUUID, result, "Retrieved UUID should match the one set in context")
	})

	t.Run("Missing UUID in context", func(t *testing.T) {
		ctx := context.Background()

		assert.Panics(t, func() {
			middleware.ReqIDFromCtx(ctx)
		}, "Function should panic when UUID is not found in context")
	})

	t.Run("Incorrect type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), middleware.ContextKey{}, "not a UUID")

		assert.Panics(t, func() {
			middleware.ReqIDFromCtx(ctx)
		}, "Function should panic when value in context is not a UUID")
	})

	t.Run("Nil context", func(t *testing.T) {
		assert.Panics(t, func() {
			middleware.ReqIDFromCtx(nil)
		}, "Function should panic when context is nil")
	})
}
