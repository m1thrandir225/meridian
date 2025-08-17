package logging

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SafeString returns a safe version of a string for logging (truncates sensitive data)
func SafeString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// GetTokenPrefix returns a safe prefix of the token for logging
func GetTokenPrefix(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:8] + "***"
}

// GetContentPreview returns a safe preview of the message content for logging
func GetContentPreview(content string) string {
	return SafeString(content, 50)
}

// LogDatabaseOperation logs database operations with timing
func (l *Logger) LogDatabaseOperation(operation string, table string, duration time.Duration, err error) {
	if err != nil {
		l.Error("Database operation failed",
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
	} else {
		l.Debug("Database operation completed",
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Duration("duration", duration),
		)
	}
}

// LogCacheOperation logs cache operations
func (l *Logger) LogCacheOperation(operation string, key string, hit bool, err error) {
	if err != nil {
		l.Error("Cache operation failed",
			zap.String("operation", operation),
			zap.String("key", key),
			zap.Bool("cache_hit", hit),
			zap.Error(err),
		)
	} else {
		l.Debug("Cache operation completed",
			zap.String("operation", operation),
			zap.String("key", key),
			zap.Bool("cache_hit", hit),
		)
	}
}

// LogHTTPRequest logs HTTP requests with timing
func (l *Logger) LogHTTPRequest(method, path, status string, duration time.Duration, err error) {
	if err != nil {
		l.Error("HTTP request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("status", status),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
	} else {
		l.Info("HTTP request completed",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("status", status),
			zap.Duration("duration", duration),
		)
	}
}

// LogKafkaOperation logs Kafka operations
func (l *Logger) LogKafkaOperation(operation, topic string, partition int32, offset int64, err error) {
	if err != nil {
		l.Error("Kafka operation failed",
			zap.String("operation", operation),
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
			zap.Error(err),
		)
	} else {
		l.Info("Kafka operation completed",
			zap.String("operation", operation),
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
	}
}

// LogServiceStart logs service startup information
func (l *Logger) LogServiceStart(serviceName, version string, config map[string]interface{}) {
	l.Info("Service starting",
		zap.String("service", serviceName),
		zap.String("version", version),
		zap.Any("config", config),
	)
}

// LogServiceStop logs service shutdown information
func (l *Logger) LogServiceStop(serviceName string, reason string) {
	l.Info("Service stopping",
		zap.String("service", serviceName),
		zap.String("reason", reason),
	)
}

// LogHealthCheck logs health check results
func (l *Logger) LogHealthCheck(component string, healthy bool, details map[string]interface{}) {
	if healthy {
		l.Debug("Health check passed",
			zap.String("component", component),
			zap.Any("details", details),
		)
	} else {
		l.Warn("Health check failed",
			zap.String("component", component),
			zap.Any("details", details),
		)
	}
}

// SanitizeLogData removes sensitive information from log data
func SanitizeLogData(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	sensitiveKeys := []string{"password", "token", "secret", "key", "authorization"}

	for k, v := range data {
		keyLower := strings.ToLower(k)
		isSensitive := false

		for _, sensitiveKey := range sensitiveKeys {
			if strings.Contains(keyLower, sensitiveKey) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			if str, ok := v.(string); ok {
				sanitized[k] = GetTokenPrefix(str)
			} else {
				sanitized[k] = "***"
			}
		} else {
			sanitized[k] = v
		}
	}

	return sanitized
}

// LogGRPCOperation logs gRPC operations
func (l *Logger) LogGRPCOperation(method string, duration time.Duration, err error) {
	if err != nil {
		l.Error("gRPC operation failed",
			zap.String("method", method),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
	} else {
		l.Debug("gRPC operation completed",
			zap.String("method", method),
			zap.Duration("duration", duration),
		)
	}
}

func GenerateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + RandomString(8)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func GetRemoteAddr(r *http.Request) string {
	// Check for forwarded headers
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

func GetStackTrace() string {
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	return string(stack[:length])
}
