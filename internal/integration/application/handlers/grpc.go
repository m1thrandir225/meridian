package handlers

import (
	"context"
	"log"
	"net"

	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	integrationpb "github.com/m1thrandir225/meridian/internal/integration/infrastructure/api"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	integrationService *services.IntegrationService
	integrationpb.UnimplementedIntegrationServiceServer
}

func NewGRPCHandler(service *services.IntegrationService) *GRPCHandler {
	return &GRPCHandler{
		integrationService: service,
	}
}

func (h *GRPCHandler) ValidateAPIToken(ctx context.Context, req *integrationpb.ValidateAPITokenRequest) (*integrationpb.ValidateAPITokenResponse, error) {
	log.Printf("gRPC: Validating API Token")

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

	return &integrationpb.ValidateAPITokenResponse{
		Valid:           isValid,
		IntegrationId:   id,
		IntegrationName: "",
	}, nil
}

func (h *GRPCHandler) GetIntegration(ctx context.Context, req *integrationpb.GetIntegrationRequest) (*integrationpb.GetIntegrationResponse, error) {
	log.Printf("gRPC: Getting integration")
	//TODO: implement
	return nil, nil
}

func StartGRPCServer(integrationService *services.IntegrationService, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(integrationService)
	integrationpb.RegisterIntegrationServiceServer(s, grpcHandler)
	log.Printf("gRPC server listening on port %s", port)
	return s.Serve(lis)
}
