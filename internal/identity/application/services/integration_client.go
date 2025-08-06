package services

import (
	"context"
	"log"
	"time"

	integrationpb "github.com/m1thrandir225/meridian/internal/integration/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (ic *IntegrationClient) ValidateAPIToken(token string) (*integrationpb.ValidateAPITokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &integrationpb.ValidateAPITokenRequest{
		Token: token,
	}

	resp, err := ic.client.ValidateAPIToken(ctx, req)
	if err != nil {
		log.Printf("gRPC call to ValidateAPIToken failed: %v", err)
		return nil, err
	}

	return resp, nil
}

func (ic *IntegrationClient) Close() error {
	return ic.conn.Close()
}
