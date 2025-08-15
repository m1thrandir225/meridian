package handlers

import (
	"context"
	"fmt"
	"log"
	"net"

	identitypb "github.com/m1thrandir225/meridian/internal/identity/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	tokenVerifier auth.TokenVerifier
	identitypb.UnimplementedIdentityServiceServer
}

func NewGRPCServer(tokenVerifier auth.TokenVerifier) *GRPCServer {
	return &GRPCServer{
		tokenVerifier: tokenVerifier,
	}
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *identitypb.ValidateTokenRequest) (*identitypb.ValidateTokenResponse, error) {
	claims, err := s.tokenVerifier.Verify(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return &identitypb.ValidateTokenResponse{
		UserId: claims.Custom.UserID,
	}, nil
}

func StartGRPCServer(port string, tokenVerifier auth.TokenVerifier) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCServer(tokenVerifier)
	identitypb.RegisterIdentityServiceServer(s, grpcHandler)

	log.Printf("Identity gRPC server listening on port %s", port)
	return s.Serve(lis)
}
