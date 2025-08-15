package handlers

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	identitypb "github.com/m1thrandir225/meridian/internal/identity/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	tokenVerifier   auth.TokenVerifier
	identityService *services.IdentityService
	identitypb.UnimplementedIdentityServiceServer
}

func NewGRPCServer(tokenVerifier auth.TokenVerifier, identityService *services.IdentityService) *GRPCServer {
	return &GRPCServer{
		tokenVerifier:   tokenVerifier,
		identityService: identityService,
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

func (s *GRPCServer) GetUserByID(ctx context.Context, req *identitypb.GetUserByIDRequest) (*identitypb.GetUserByIDResponse, error) {

	cmd := domain.GetUserCommand{
		UserID: req.UserId,
	}

	user, err := s.identityService.GetUser(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}
	return &identitypb.GetUserByIDResponse{
		User: &identitypb.User{
			Id:        user.ID.String(),
			Username:  user.Username.String(),
			Email:     user.Email.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *GRPCServer) GetUsers(ctx context.Context, req *identitypb.GetUsersRequest) (*identitypb.GetUsersResponse, error) {
	cmd := domain.GetUsersCommand{
		UserIds: req.UserIds,
	}
	users, err := s.identityService.GetUsers(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	pbUsers := make([]*identitypb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &identitypb.User{
			Id:        user.ID.String(),
			Username:  user.Username.String(),
			Email:     user.Email.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}
	return &identitypb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}

func StartGRPCServer(port string, tokenVerifier auth.TokenVerifier, identityService *services.IdentityService) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCServer(tokenVerifier, identityService)
	identitypb.RegisterIdentityServiceServer(s, grpcHandler)

	log.Printf("Identity gRPC server listening on port %s", port)
	return s.Serve(lis)
}
