package logging

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLoggingMiddleware creates a Gin middleware for logging HTTP requests
func GinLoggingMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = GenerateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// Create logger with request context
		requestLogger := logger.WithRequestID(requestID).
			WithHTTPMethod(c.Request.Method).
			WithPath(c.Request.URL.Path).
			WithRemoteAddr(c.ClientIP()).
			WithUserAgent(c.Request.UserAgent())

		// Log request start
		requestLogger.Info("HTTP request started",
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("referer", c.Request.Referer()),
		)

		// Capture response
		responseWriter := &HTTPResponseWriter{
			ResponseWriter: c.Writer,
			statusCode:     200,
			body:           nil,
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		requestLogger = requestLogger.WithStatusCode(responseWriter.statusCode).
			WithDuration(duration)

		switch {
		case responseWriter.statusCode >= 500:
			requestLogger.Error("HTTP request failed (server error)",
				zap.String("error", c.Errors.String()),
			)
		case responseWriter.statusCode >= 400:
			requestLogger.Warn("HTTP request failed (client error)",
				zap.String("error", c.Errors.String()),
			)
		default:
			requestLogger.Info("HTTP request completed")
		}
	}
}

func StandardHTTPLoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = GenerateRequestID()
				r.Header.Set("X-Request-ID", requestID)
			}

			requestLogger := logger.WithRequestID(requestID).
				WithHTTPMethod(r.Method).
				WithPath(r.URL.Path).
				WithRemoteAddr(GetRemoteAddr(r)).
				WithUserAgent(r.UserAgent())

			// Log request start
			requestLogger.Info("HTTP request started",
				zap.String("query", r.URL.RawQuery),
				zap.String("referer", r.Referer()),
			)

			// Capture response
			responseWriter := &HTTPResponseWriter{
				ResponseWriter: w,
				statusCode:     200,
				body:           nil,
			}

			next.ServeHTTP(responseWriter, r)

			duration := time.Since(start)

			requestLogger = requestLogger.WithStatusCode(responseWriter.statusCode).
				WithDuration(duration)

			switch {
			case responseWriter.statusCode >= 500:
				requestLogger.Error("HTTP request failed (server error)")
			case responseWriter.statusCode >= 400:
				requestLogger.Warn("HTTP request failed (client error)")
			default:
				requestLogger.Info("HTTP request completed")
			}
		})
	}
}

func GinRecoveryMiddleware(logger *Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = GenerateRequestID()
		}

		requestLogger := logger.WithRequestID(requestID).
			WithHTTPMethod(c.Request.Method).
			WithPath(c.Request.URL.Path)

		requestLogger.Error("HTTP request panicked",
			zap.Any("panic", recovered),
			zap.String("stack", GetStackTrace()),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func BodyLoggingMiddleware(logger *Logger, maxBodySize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil && c.Request.ContentLength > 0 && c.Request.ContentLength <= maxBodySize {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				requestLogger := logger.WithRequestID(c.GetHeader("X-Request-ID"))
				requestLogger.Debug("HTTP request body",
					zap.String("content_type", c.GetHeader("Content-Type")),
					zap.String("body", string(bodyBytes)),
				)
			}
		}

		// Capture response body
		responseWriter := &HTTPResponseWriter{
			ResponseWriter: c.Writer,
			statusCode:     200,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = responseWriter

		c.Next()

		// Log response body if present and not too large
		if responseWriter.body.Len() > 0 && responseWriter.body.Len() <= int(maxBodySize) {
			requestLogger := logger.WithRequestID(c.GetHeader("X-Request-ID"))
			requestLogger.Debug("HTTP response body",
				zap.String("content_type", c.Writer.Header().Get("Content-Type")),
				zap.String("body", responseWriter.body.String()),
			)
		}
	}
}

func HealthCheckFilter(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		GinLoggingMiddleware(logger)(c)
	}
}

// RateLimitLoggingMiddleware creates a middleware that logs rate limit events
func RateLimitLoggingMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if rate limited (you would set this in your rate limiting middleware)
		if rateLimited := c.GetBool("rate_limited"); rateLimited {
			requestLogger := logger.WithRequestID(c.GetHeader("X-Request-ID")).
				WithHTTPMethod(c.Request.Method).
				WithPath(c.Request.URL.Path).
				WithRemoteAddr(c.ClientIP())

			requestLogger.Warn("HTTP request rate limited")
		}
	}
}
