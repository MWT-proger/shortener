package grpc

import (
	"context"
	"fmt"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	pb "github.com/MWT-proger/shortener/internal/shortener/grpc/proto"
)

// GenerateShortKey.
func (s *ShortenerGRPCServer) GenerateShortKey(ctx context.Context, in *pb.GenerateShortKeyRequest) (*pb.GenerateShortKeyResponse, error) {
	var response pb.GenerateShortKeyResponse

	// var (
	// 	finalStatusCode = http.StatusCreated
	// )

	userID, ok := auth.UserIDFrom(ctx)
	fmt.Println(userID)
	if !ok {
		// http.Error(w, "", http.StatusInternalServerError)
		// return
	}

	// shortURL, err := h.shortService.GenerateShortURL(ctx, userID, in.Url, r.Host)

	// if err != nil {
	// 	finalStatusCode = h.setOrGetHTTPCode(w, err)

	// 	if finalStatusCode == 0 {
	// 		return
	// 	}
	// }

	// w.Header().Set("content-type", "text/plain")
	// w.WriteHeader(finalStatusCode)
	// w.Write([]byte(shortURL))

	return &response, nil
}

// GetFullURLByShortKey.
func (s *ShortenerGRPCServer) GetFullURLByShortKey(ctx context.Context, in *pb.GetFullURLByShortKeyRequest) (*pb.GetFullURLByShortKeyResponse, error) {
	var response pb.GetFullURLByShortKeyResponse

	return &response, nil
}

// GetListUserURLs.
func (s *ShortenerGRPCServer) GetListUserURLs(ctx context.Context, in *pb.GetListUserURLsRequest) (*pb.GetListUserURLsResponse, error) {
	var response pb.GetListUserURLsResponse

	return &response, nil
}

// DeleteListUserURLs.
func (s *ShortenerGRPCServer) DeleteListUserURLs(ctx context.Context, in *pb.DeleteListUserURLsRequest) (*pb.DeleteListUserURLsResponse, error) {
	var response pb.DeleteListUserURLsResponse

	return &response, nil
}
