package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	content     []byte
	contentType string
	statusCode  int
	timestamp   time.Time
}

type Cache struct {
	store map[string]cacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]cacheEntry),
		ttl:   ttl,
	}
}

func (c *Cache) Set(key string, entry cacheEntry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = entry
}

func (c *Cache) Get(key string) (cacheEntry, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, found := c.store[key]
	if !found {
		return cacheEntry{}, false
	}
	if time.Since(entry.timestamp) > c.ttl {
		delete(c.store, key)
		return cacheEntry{}, false
	}
	return entry, true
}

type CacheKeyFunc func(*http.Request) string

func CacheMiddleware(cache *Cache, keyFunc CacheKeyFunc) Middleware {
	if keyFunc == nil {
		keyFunc = func(r *http.Request) string {
			return r.URL.Path
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := keyFunc(r)

			if entry, found := cache.Get(key); found {
				w.Header().Set("Content-Type", entry.contentType)
				w.WriteHeader(entry.statusCode)
				w.Write(entry.content)
				return
			}

			buf := &bytes.Buffer{}
			newW := &responseWriter{
				ResponseWriter: w,
				buf:            buf,
			}

			next.ServeHTTP(newW, r)

			cache.Set(key, cacheEntry{
				content:     buf.Bytes(),
				contentType: newW.Header().Get("Content-Type"),
				statusCode:  newW.statusCode,
				timestamp:   time.Now(),
			})

			w.WriteHeader(newW.statusCode)
			w.Write(buf.Bytes())
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.buf.Write(b)
	return rw.ResponseWriter.Write(b)
}
