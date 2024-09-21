package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type contextKey struct{}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = SetReqID(r)

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)
			slog.Info(
				"Request handled",
				"requestID", ReqIDFromCtx(r.Context()),
				"statusCode", wrapped.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
			)
		})
	}
}

func ReqIDFromCtx(ctx context.Context) uuid.UUID {
	v := ctx.Value(contextKey{})
	if v == nil {
		panic("uuid for request not found")
	}

	switch id := v.(type) {
	case uuid.UUID:
		return id
	default:
		panic("uuid for request not found")
	}
}

func SetReqID(r *http.Request) *http.Request {
	requestID := uuid.New()
	ctx := context.WithValue(r.Context(), contextKey{}, requestID)
	return r.WithContext(ctx)
}
