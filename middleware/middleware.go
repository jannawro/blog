package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type CacheKeyFunc func(*http.Request) string

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
