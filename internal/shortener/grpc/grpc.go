package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/auth"
	pb "github.com/MWT-proger/shortener/internal/shortener/grpc/proto"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
)

// ShortenerGRPCServer поддерживает все необходимые методы сервера.
type ShortenerGRPCServer struct {
	pb.UnimplementedShortenerServer
	shortService ShortenerServicer
}

// NewShortenerGRPCServer создает новую структуру ShortenerGRPCServer.
func NewShortenerGRPCServer(service ShortenerServicer) (s *ShortenerGRPCServer, err error) {
	s = &ShortenerGRPCServer{
		shortService: service,
	}

	return s, err
}

// ShortenerServicer интерфейс описывающий необходимые методы для сервисного слоя.
type ShortenerServicer interface {
	GetFullURLByShortKey(ctx context.Context, shortKey string) (string, error)
	GetListUserURLs(ctx context.Context, userID uuid.UUID, requestHost string) ([]*models.JSONShortURL, error)

	GenerateShortURL(ctx context.Context, userID uuid.UUID, fullURL string, requestHost string) (string, error)
	GenerateMultyShortURL(ctx context.Context, userID uuid.UUID, data []models.JSONShortURL, requestHost string) error

	DeleteListUserURLs(ctx context.Context, userID uuid.UUID, data []string)
	GetStats(ctx context.Context) (urls int, users int, err error)

	PingStorage() bool
}

// Handler интерфейс определяет необходимые методы для инициализации маршрутизатора.
// type Handler interface {
// 	JSONMultyGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request)
// 	PingDB(w http.ResponseWriter, r *http.Request)
// 	GetStats(w http.ResponseWriter, r *http.Request)
// }

// Run запускает gRPC server на установленном порту.
func (s *ShortenerGRPCServer) Run(ctx context.Context, conf configs.Config) error {

	listen, err := net.Listen("tcp", conf.HostServer)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	pb.RegisterShortenerServer(gRPCServer, s)

	logger.Log.Info("Running GRPC server on", logger.StringField("host", conf.HostServer))

	return gRPCServer.Serve(listen)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	UserID := uuid.New()
	ctx = auth.WithUserID(ctx, UserID)

	return handler(ctx, req)
}
