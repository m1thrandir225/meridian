package services

import (
	"context"
	"log"

	integrationpb "github.com/m1thrandir225/meridian/internal/integration/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IntegrationClient struct {
	client integrationpb.IntegrationServiceClient
	conn   *grpc.ClientConn
}

func NewIntegrationClient(address string) (*IntegrationClient, error) {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := integrationpb.NewIntegrationServiceClient(conn)
	return &IntegrationClient{
		client: client,
		conn:   conn,
	}, nil
}

func (ic *IntegrationClient) GetIntegration(ctx context.Context, integrationID string) (*integrationpb.GetIntegrationResponse, error) {
	req := &integrationpb.GetIntegrationRequest{
		IntegrationId: integrationID,
	}

	resp, err := ic.client.GetIntegration(ctx, req)
	if err != nil {
		log.Printf("gRPC call to GetIntegration failed: %v", err)
		return nil, err
	}

	return resp, nil
}

func (ic *IntegrationClient) GetIntegrations(ctx context.Context, integrationIDs []string) ([]*integrationpb.Integration, error) {
	var integrations []*integrationpb.Integration

	for _, integrationID := range integrationIDs {
		resp, err := ic.GetIntegration(ctx, integrationID)
		if err != nil {
			log.Printf("Failed to get integration %s: %v", integrationID, err)
			continue
		}
		integrations = append(integrations, resp.Integration)
	}

	return integrations, nil
}

func (ic *IntegrationClient) Close() error {
	return ic.conn.Close()
}
