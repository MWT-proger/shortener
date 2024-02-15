package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	pb "github.com/MWT-proger/shortener/internal/shortener/grpc/proto"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GenerateShortKey.
func (s *ShortenerGRPCServer) GenerateShortKey(ctx context.Context, in *pb.GenerateShortKeyRequest) (*pb.GenerateShortKeyResponse, error) {
	var response pb.GenerateShortKeyResponse

	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		return nil, status.Error(codes.Internal, "")
	}

	shortURL, err := s.shortService.GenerateShortURL(ctx, userID, in.Url, "")
	if err != nil {
		if err := s.setOrGetGRPCCode(err); err != nil {
			return nil, err
		}
	}

	response.Result = shortURL

	return &response, nil
}

// GetFullURLByShortKey.
func (s *ShortenerGRPCServer) GetFullURLByShortKey(ctx context.Context, in *pb.GetFullURLByShortKeyRequest) (*pb.GetFullURLByShortKeyResponse, error) {
	var response pb.GetFullURLByShortKeyResponse

	fullURL, err := s.shortService.GetFullURLByShortKey(ctx, in.ShortKey)

	if err != nil {
		if err := s.setOrGetGRPCCode(err); err != nil {
			return nil, err
		}
	}
	response.FullUrl = fullURL

	return &response, nil
}

// GetListUserURLs.
func (s *ShortenerGRPCServer) GetListUserURLs(ctx context.Context, in *pb.GetListUserURLsRequest) (*pb.GetListUserURLsResponse, error) {
	var response pb.GetListUserURLsResponse

	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		return nil, status.Error(codes.Internal, "")
	}

	listURLs, err := s.shortService.GetListUserURLs(ctx, userID, "")

	if err != nil {
		if err := s.setOrGetGRPCCode(err); err != nil {
			return nil, err
		}
	}
	for _, v := range listURLs {
		response.Urls = append(response.Urls, &pb.URL{OriginalUrl: v.OriginalURL, ShortUrl: v.ShortURL})
	}

	return &response, nil
}

// DeleteListUserURLs.
func (s *ShortenerGRPCServer) DeleteListUserURLs(ctx context.Context, in *pb.DeleteListUserURLsRequest) (*pb.DeleteListUserURLsResponse, error) {
	var response pb.DeleteListUserURLsResponse

	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		return nil, status.Error(codes.Internal, "")
	}
	fmt.Println(userID, in.ShortKeys)
	s.shortService.DeleteListUserURLs(ctx, userID, in.ShortKeys)

	return &response, nil
}

// setOrGetHTTPCode присваивает response статус ответа при ошибке в сервисном слою.
// Возвращает HTTPCode, при необходимости досрочного выполнения кода отправляе HTTPCode = 0.
func (s *ShortenerGRPCServer) setOrGetGRPCCode(err error) error {
	var serviceError *lErrors.ServicesError

	if errors.As(err, &serviceError) {
		if serviceError.GRPCCode == codes.OK {
			return nil
		}
		return status.Error(serviceError.GRPCCode, serviceError.Error())

	} else {
		return status.Error(codes.Internal, "Ошибка сервера, попробуйте позже.")
	}
}
