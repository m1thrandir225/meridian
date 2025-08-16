package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

type CacheMiddleware struct {
	cache *cache.RedisCache
	ttl   time.Duration
}

func NewCacheMiddleware(cache *cache.RedisCache, ttl time.Duration) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
		ttl:   ttl,
	}
}

func (m *CacheMiddleware) CacheResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		cacheKey := m.generateCacheKey(c)

		var cachedResponse map[string]interface{}
		if err := m.cache.Get(c.Request.Context(), cacheKey, &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CacheResponseWriter is a middleware that caches the response of a GET request
func (m *CacheMiddleware) CacheResponseWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		cacheKey := m.generateCacheKey(c)

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() == http.StatusOK && writer.body.Len() > 0 {
			var response map[string]interface{}
			if json.Unmarshal(writer.body.Bytes(), &response) == nil {
				m.cache.Set(c.Request.Context(), cacheKey, response, m.ttl)
			}
		}
	}
}

// generateCacheKey generates a cache key for a given request
func (m *CacheMiddleware) generateCacheKey(c *gin.Context) string {
	userID := c.GetHeader("X-User-ID")
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	keyData := path + "?" + query
	if userID != "" {
		keyData = userID + ":" + keyData
	}

	hash := md5.Sum([]byte(keyData))
	return "cache:" + hex.EncodeToString(hash[:])
}
