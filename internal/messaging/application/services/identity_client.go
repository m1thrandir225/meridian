package services

import (
	"context"
	"log"
	"time"

	identitypb "github.com/m1thrandir225/meridian/internal/identity/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdentityClient struct {
	client identitypb.IdentityServiceClient
	conn   *grpc.ClientConn
}

func NewIdentityClient(address string) (*IdentityClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := identitypb.NewIdentityServiceClient(conn)
	return &IdentityClient{
		client: client,
		conn:   conn,
	}, nil
}

func (ic *IdentityClient) ValidateToken(token string) (*identitypb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &identitypb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := ic.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("gRPC call to ValidateToken failed: %v", err)
		return nil, err
	}

	return resp, nil
}

func (ic *IdentityClient) GetUserByID(userID string) (*identitypb.GetUserByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &identitypb.GetUserByIDRequest{
		UserId: userID,
	}

	resp, err := ic.client.GetUserByID(ctx, req)
	if err != nil {
		log.Printf("gRPC call to GetUserByID failed: %v", err)
		return nil, err
	}

	return resp, nil
}

func (ic *IdentityClient) GetUsers(context context.Context, userIDs []string) (*identitypb.GetUsersResponse, error) {

	req := &identitypb.GetUsersRequest{
		UserIds: userIDs,
	}
	resp, err := ic.client.GetUsers(context, req)
	if err != nil {
		log.Printf("gRPC call to GetUsers failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (ic *IdentityClient) Close() error {
	return ic.conn.Close()
}
