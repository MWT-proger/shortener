package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MWT-proger/shortener/configs"
)

// RequestLoggerInterceptor — interceptor-логер для входящих HTTP-запросов.
func RequestLoggerInterceptor(conf configs.Config) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		var UserID uuid.UUID
		token, ok := getTokenAuth(ctx)
		if ok {
			UserID = getUserID(conf, token)

			if UserID == uuid.Nil {
				return nil, status.Error(codes.Unauthenticated, "")
			}

		} else {
			UserID = uuid.New()

			token, err = buildJWTString(conf, UserID)

			if err != nil {
				return nil, status.Error(codes.Internal, "")
			}

		}
		grpc.SetTrailer(ctx, metadata.New(map[string]string{gRPCnameKey: token}))

		ctx = WithUserID(ctx, UserID)

		resp, err = handler(ctx, req)
		return resp, err

	}

}

func getTokenAuth(ctx context.Context) (token string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get(gRPCnameKey)
		fmt.Println(values)
		if len(values) > 0 {
			token = values[0]
			return token, ok
		}
	}
	return "", false
}
