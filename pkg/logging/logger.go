package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(config Config) *Logger {
	var zapConfig zap.Config

	switch config.Environment {
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		zapConfig = zap.NewProductionConfig()
	}

	if level, err := zapcore.ParseLevel(config.LogLevel); err == nil {
		zapConfig.Level = zap.NewAtomicLevelAt(level)
	}
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.StacktraceKey = "stacktrace"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.LevelKey = "level"
	zapConfig.EncoderConfig.CallerKey = "caller"

	zapConfig.InitialFields = map[string]interface{}{
		"service": config.ServiceName,
	}

	logger, err := zapConfig.Build()
	if err != nil {
		basicLogger, _ := zap.NewDevelopment()
		return &Logger{Logger: basicLogger}
	}

	return &Logger{Logger: logger}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	fields := make([]zap.Field, 0, len(ctx))
	for k, v := range ctx {
		fields = append(fields, zap.Any(k, v))
	}
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithRequestID adds a request ID to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("request_id", requestID))}
}

// WithMethod adds a gRPC method name to the logger
func (l *Logger) WithMethod(method string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("grpc_method", method))}
}

// WithHTTPMethod adds an HTTP method name to the logger
func (l *Logger) WithHTTPMethod(method string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("http_method", method))}
}

// WithPath adds an HTTP path to the logger
func (l *Logger) WithPath(path string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("http_path", path))}
}

// WithStatusCode adds an HTTP status code to the logger
func (l *Logger) WithStatusCode(statusCode int) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Int("http_status", statusCode))}
}

// WithUserID adds a user ID to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("user_id", userID))}
}

// WithChannelID adds a channel ID to the logger
func (l *Logger) WithChannelID(channelID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("channel_id", channelID))}
}

// WithIntegrationID adds an integration ID to the logger
func (l *Logger) WithIntegrationID(integrationID string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("integration_id", integrationID))}
}

// WithError adds an error to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Error(err))}
}

// WithDuration adds a duration field to the logger
func (l *Logger) WithDuration(duration interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Any("duration", duration))}
}

// WithUserAgent adds a user agent to the logger
func (l *Logger) WithUserAgent(userAgent string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("user_agent", userAgent))}
}

// WithRemoteAddr adds a remote address to the logger
func (l *Logger) WithRemoteAddr(remoteAddr string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.String("remote_addr", remoteAddr))}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// GetZapLogger returns the underlying zap.Logger for advanced usage
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.Logger
}
