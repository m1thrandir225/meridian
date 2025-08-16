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
	"google.golang.org/grpc"
)

type GRPCServer struct {
	integrationService *services.IntegrationService
	cache              *cache.RedisCache
	integrationpb.UnimplementedIntegrationServiceServer
}

func NewGRPCHandler(
	service *services.IntegrationService,
	cache *cache.RedisCache,
) *GRPCServer {
	return &GRPCServer{
		integrationService: service,
		cache:              cache,
	}
}

func (h *GRPCServer) ValidateAPIToken(ctx context.Context, req *integrationpb.ValidateAPITokenRequest) (*integrationpb.ValidateAPITokenResponse, error) {
	log.Printf("gRPC: Validating API Token")

	cacheKey := fmt.Sprintf("grpc_api_token_validation:%s", req.Token)
	var cachedResponse integrationpb.ValidateAPITokenResponse
	if hit, _ := h.cache.GetWithMetrics(ctx, cacheKey, &cachedResponse); hit {
		return &cachedResponse, nil
	}

	isValid, id, _, err := h.integrationService.ValidateApiToken(ctx, req.Token)
	if err != nil {
		return &integrationpb.ValidateAPITokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	if !isValid {
		return &integrationpb.ValidateAPITokenResponse{
			Valid: false,
		}, nil
	}

	cmd := domain.GetIntegrationCommand{
		IntegrationID: id,
	}

	integration, err := h.integrationService.GetIntegration(ctx, cmd)
	if err != nil {
		return nil, err
	}

	response := &integrationpb.ValidateAPITokenResponse{
		Valid:            isValid,
		IntegrationId:    id,
		IntegrationName:  integration.ServiceName,
		TargetChannelIds: integration.TargetChannelIDsAsStringSlice(),
	}
	h.cache.Set(ctx, cacheKey, response, 15*time.Minute)

	return response, nil
}

func (h *GRPCServer) GetIntegration(ctx context.Context, req *integrationpb.GetIntegrationRequest) (*integrationpb.GetIntegrationResponse, error) {
	log.Printf("gRPC: Getting integration")

	cacheKey := fmt.Sprintf("grpc_integration:%s", req.IntegrationId)
	var cachedResponse integrationpb.GetIntegrationResponse
	if hit, _ := h.cache.GetWithMetrics(ctx, cacheKey, &cachedResponse); hit {
		return &cachedResponse, nil
	}

	cmd := domain.GetIntegrationCommand{
		IntegrationID: req.IntegrationId,
	}

	integration, err := h.integrationService.GetIntegration(ctx, cmd)
	if err != nil {
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

	return response, nil
}

func StartGRPCServer(
	integrationService *services.IntegrationService,
	cache *cache.RedisCache,
	port string,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(integrationService, cache)
	integrationpb.RegisterIntegrationServiceServer(s, grpcHandler)
	log.Printf("gRPC server listening on port %s", port)
	return s.Serve(lis)
}
