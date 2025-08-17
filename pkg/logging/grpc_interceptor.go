package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerLoggingInterceptor(logger *Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Create a logger with method context
		methodLogger := logger.WithMethod(info.FullMethod)

		// Extract request ID from metadata if available
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
				methodLogger = methodLogger.WithRequestID(requestIDs[0])
			}
		}

		// Log the incoming request
		methodLogger.Info("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.String("peer", getPeerFromContext(ctx)),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		if err != nil {
			st, _ := status.FromError(err)
			methodLogger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.String("error", err.Error()),
				zap.String("error_code", st.Code().String()),
				zap.Duration("duration", duration),
				zap.String("peer", getPeerFromContext(ctx)),
			)
		} else {
			methodLogger.Info("gRPC request completed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.String("peer", getPeerFromContext(ctx)),
			)
		}

		return resp, err
	}
}

// StreamServerLoggingInterceptor creates a gRPC stream interceptor that logs stream operations
func StreamServerLoggingInterceptor(logger *Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Create a logger with method context
		methodLogger := logger.WithMethod(info.FullMethod)

		// Extract request ID from metadata if available
		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
				methodLogger = methodLogger.WithRequestID(requestIDs[0])
			}
		}

		// Log the stream start
		methodLogger.Info("gRPC stream started",
			zap.String("method", info.FullMethod),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
			zap.String("peer", getPeerFromContext(ss.Context())),
		)

		// Call the handler
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Log the stream completion
		if err != nil {
			st, _ := status.FromError(err)
			methodLogger.Error("gRPC stream failed",
				zap.String("method", info.FullMethod),
				zap.String("error", err.Error()),
				zap.String("error_code", st.Code().String()),
				zap.Duration("duration", duration),
				zap.String("peer", getPeerFromContext(ss.Context())),
			)
		} else {
			methodLogger.Info("gRPC stream completed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.String("peer", getPeerFromContext(ss.Context())),
			)
		}

		return err
	}
}

// UnaryClientLoggingInterceptor creates a gRPC unary client interceptor for logging
func UnaryClientLoggingInterceptor(logger *Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		methodLogger := logger.WithMethod(method)

		methodLogger.Debug("gRPC client request started",
			zap.String("method", method),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			methodLogger.Error("gRPC client request failed",
				zap.String("method", method),
				zap.String("error", err.Error()),
				zap.String("error_code", st.Code().String()),
				zap.Duration("duration", duration),
			)
		} else {
			methodLogger.Debug("gRPC client request completed",
				zap.String("method", method),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}

// getPeerFromContext extracts peer information from the context
func getPeerFromContext(ctx context.Context) string {
	// This is a simplified version - in a real implementation you might want to
	// extract more detailed peer information from the context
	return "unknown"
}

// LoggingServerOption returns a gRPC server option with logging interceptors
func LoggingServerOption(logger *Logger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		UnaryServerLoggingInterceptor(logger),
	)
}

// LoggingStreamServerOption returns a gRPC server option with stream logging interceptors
func LoggingStreamServerOption(logger *Logger) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(
		StreamServerLoggingInterceptor(logger),
	)
}

// LoggingClientOption returns a gRPC client option with logging interceptors
func LoggingClientOption(logger *Logger) grpc.DialOption {
	return grpc.WithUnaryInterceptor(
		UnaryClientLoggingInterceptor(logger),
	)
}
