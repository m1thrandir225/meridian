package services

import (
	"context"
	"log"
	"strings"

	integrationpb "github.com/m1thrandir225/meridian/internal/integration/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type IntegrationClient struct {
	client integrationpb.IntegrationServiceClient
	conn   *grpc.ClientConn
}

func NewIntegrationClient(address string) (*IntegrationClient, error) {
	conn, err := grpc.NewClient(address,
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
	req := &integrationpb.GetIntegrationRequest{IntegrationId: integrationID}

	resp, err := ic.client.GetIntegration(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.NotFound ||
			(st.Code() == codes.Unknown && strings.Contains(st.Message(), "integration not found"))) {
			return nil, nil
		}
		log.Printf("gRPC call to GetIntegration failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (ic *IntegrationClient) GetIntegrations(ctx context.Context, integrationIDs []string) ([]*integrationpb.Integration, error) {
	var integrations []*integrationpb.Integration
	for _, integrationID := range integrationIDs {
		resp, err := ic.GetIntegration(ctx, integrationID)
		if err != nil || resp == nil || resp.Integration == nil {
			continue
		}
		integrations = append(integrations, resp.Integration)
	}
	return integrations, nil
}

func (ic *IntegrationClient) Close() error {
	return ic.conn.Close()
}
