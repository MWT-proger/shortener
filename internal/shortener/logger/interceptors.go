package logger

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// RequestLoggerInterceptor — interceptor-логер для входящих HTTP-запросов.
func RequestLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	timeStart := time.Now()

	defer func() {
		Log.Info("got incoming gRPC request",
			StringField("path", info.FullMethod),
			StringField("status", status.Code(err).String()),
			DurationField("time", time.Since(timeStart)),
		)
	}()

	resp, err = handler(ctx, req)
	return resp, err

}
