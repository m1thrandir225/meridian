package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
	integrationpb "github.com/m1thrandir225/meridian/internal/integration/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	integrationService *services.IntegrationService
	cache              *cache.RedisCache
	logger             *logging.Logger
	integrationpb.UnimplementedIntegrationServiceServer
}

func NewGRPCHandler(
	service *services.IntegrationService,
	cache *cache.RedisCache,
	logger *logging.Logger,
) *GRPCServer {
	return &GRPCServer{
		integrationService: service,
		cache:              cache,
		logger:             logger,
	}
}

func (h *GRPCServer) ValidateAPIToken(ctx context.Context, req *integrationpb.ValidateAPITokenRequest) (*integrationpb.ValidateAPITokenResponse, error) {
	logger := h.logger.WithMethod("ValidateAPIToken")
	logger.Info("Validating API Token")

	isValid, id, _, err := h.integrationService.ValidateApiToken(ctx, req.Token)
	if err != nil {
		logger.Error("Error validating API Token", zap.Error(err))
		return &integrationpb.ValidateAPITokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	if !isValid {
		logger.Warn("Invalid API Token provided")
		return &integrationpb.ValidateAPITokenResponse{
			Valid: false,
		}, nil
	}

	cmd := domain.GetIntegrationCommand{
		IntegrationID: id,
	}

	integration, err := h.integrationService.GetIntegration(ctx, cmd)
	if err != nil {
		logger.Error("Error getting integration", zap.Error(err))
		return nil, err
	}

	response := &integrationpb.ValidateAPITokenResponse{
		Valid:            isValid,
		IntegrationId:    id,
		IntegrationName:  integration.ServiceName,
		TargetChannelIds: integration.TargetChannelIDsAsStringSlice(),
	}

	logger.Info("API Token validation successful")

	return response, nil
}

func (h *GRPCServer) GetIntegration(ctx context.Context, req *integrationpb.GetIntegrationRequest) (*integrationpb.GetIntegrationResponse, error) {
	logger := h.logger.WithMethod("GetIntegration")
	logger.Info("Getting integration")

	cacheKey := fmt.Sprintf("grpc_integration:%s", req.IntegrationId)
	var cachedResponse integrationpb.GetIntegrationResponse
	if hit, _ := h.cache.GetWithMetrics(ctx, cacheKey, &cachedResponse); hit {
		logger.Info("Integration hit cache")
		return &cachedResponse, nil
	}

	cmd := domain.GetIntegrationCommand{
		IntegrationID: req.IntegrationId,
	}

	integration, err := h.integrationService.GetIntegration(ctx, cmd)
	if err != nil {
		logger.Error("Error getting integration", zap.Error(err))
		return nil, err
	}

	targetChannelIds := make([]string, len(integration.TargetChannelIDs))
	for i, id := range integration.TargetChannelIDs {
		targetChannelIds[i] = string(id)
	}

	response := &integrationpb.GetIntegrationResponse{
		Integration: &integrationpb.Integration{
			Id:               integration.ID.String(),
			ServiceName:      integration.ServiceName,
			TargetChannelIds: targetChannelIds,
			CreatedAt:        integration.CreatedAt.Format(time.RFC3339),
			CreatorUserId:    integration.CreatorUserID.String(),
			HashedApiToken:   integration.HashedAPIToken.Hash(),
			TokenLookupHash:  integration.TokenLookupHash,
			IsRevoked:        integration.IsRevoked,
		},
	}
	h.cache.Set(ctx, cacheKey, response, 15*time.Minute)
	logger.Info("Integration retrieved", zap.String("integration_id", integration.ID.String()))

	return response, nil
}

func StartGRPCServer(
	integrationService *services.IntegrationService,
	cache *cache.RedisCache,
	logger *logging.Logger,
	port string,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(integrationService, cache, logger)
	integrationpb.RegisterIntegrationServiceServer(s, grpcHandler)
	log.Printf("gRPC server listening on port %s", port)
	return s.Serve(lis)
}
