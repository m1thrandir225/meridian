package handlers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	identitypb "github.com/m1thrandir225/meridian/internal/identity/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	tokenVerifier   auth.TokenVerifier
	identityService *services.IdentityService
	cache           *cache.RedisCache
	logger          *logging.Logger
	identitypb.UnimplementedIdentityServiceServer
}

func NewGRPCServer(
	tokenVerifier auth.TokenVerifier,
	identityService *services.IdentityService,
	cache *cache.RedisCache,
	logger *logging.Logger,
) *GRPCServer {
	return &GRPCServer{
		tokenVerifier:   tokenVerifier,
		identityService: identityService,
		cache:           cache,
		logger:          logger,
	}
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *identitypb.ValidateTokenRequest) (*identitypb.ValidateTokenResponse, error) {
	logger := s.logger.WithMethod("ValidateToken")
	logger.Info("Validating token")

	cacheKey := fmt.Sprintf("grpc_token_validation:%s", req.Token)
	var cachedResponse identitypb.ValidateTokenResponse
	if hit, _ := s.cache.GetWithMetrics(ctx, cacheKey, &cachedResponse); hit {
		logger.Info("Token validation hit cache")
		return &cachedResponse, nil
	}

	claims, err := s.tokenVerifier.Verify(req.Token)
	if err != nil {
		logger.Error("Error validating token", zap.Error(err))
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	response := &identitypb.ValidateTokenResponse{
		UserId: claims.Custom.UserID,
	}

	s.cache.Set(ctx, cacheKey, response, 15*time.Minute)
	logger.Info("Token validation successful")

	return response, nil
}

func (s *GRPCServer) GetUserByID(ctx context.Context, req *identitypb.GetUserByIDRequest) (*identitypb.GetUserByIDResponse, error) {
	logger := s.logger.WithMethod("GetUserByID")
	logger.Info("Getting user by ID")

	cacheKey := fmt.Sprintf("grpc_user:%s", req.UserId)
	var cachedUser identitypb.GetUserByIDResponse
	if hit, _ := s.cache.GetWithMetrics(ctx, cacheKey, &cachedUser); hit {
		logger.Info("User hit cache")
		return &cachedUser, nil
	}

	cmd := domain.GetUserCommand{
		UserID: req.UserId,
	}

	user, err := s.identityService.GetUser(ctx, cmd)
	if err != nil {
		logger.Error("Error getting user by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}

	response := &identitypb.GetUserByIDResponse{
		User: &identitypb.User{
			Id:        user.ID.String(),
			Username:  user.Username.String(),
			Email:     user.Email.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	s.cache.Set(ctx, cacheKey, response, 15*time.Minute)
	logger.Info("User retrieved", zap.String("user_id", user.ID.String()))

	return response, nil
}

func (s *GRPCServer) GetUsers(ctx context.Context, req *identitypb.GetUsersRequest) (*identitypb.GetUsersResponse, error) {
	logger := s.logger.WithMethod("GetUsers")
	logger.Info("Getting users")

	cacheKey := fmt.Sprintf("grpc_users:%s", req.UserIds)
	var cachedUsers identitypb.GetUsersResponse
	if hit, _ := s.cache.GetWithMetrics(ctx, cacheKey, &cachedUsers); hit {
		logger.Info("Users hit cache")
		return &cachedUsers, nil
	}

	cmd := domain.GetUsersCommand{
		UserIds: req.UserIds,
	}
	users, err := s.identityService.GetUsers(ctx, cmd)
	if err != nil {
		logger.Error("Error getting users", zap.Error(err))
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

	response := &identitypb.GetUsersResponse{
		Users: pbUsers,
	}

	s.cache.Set(ctx, cacheKey, response, 15*time.Minute)
	logger.Info("Users retrieved", zap.Int("count", len(users)))

	return response, nil
}

func StartGRPCServer(
	port string,
	tokenVerifier auth.TokenVerifier,
	identityService *services.IdentityService,
	cache *cache.RedisCache,
	logger *logging.Logger,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCServer(tokenVerifier, identityService, cache, logger)
	identitypb.RegisterIdentityServiceServer(s, grpcHandler)

	logger.Info("Identity gRPC server listening on port", zap.String("port", port))
	return s.Serve(lis)
}
