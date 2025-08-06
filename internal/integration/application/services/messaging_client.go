package services

import (
	"context"
	"log"
	"time"

	messagingpb "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MessagingClient struct {
	client messagingpb.MessagingServiceClient
	conn   *grpc.ClientConn
}

func NewMessagingClient(address string) (*MessagingClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	client := messagingpb.NewMessagingServiceClient(conn)

	return &MessagingClient{
		client: client,
		conn:   conn,
	}, nil
}

func (mc *MessagingClient) SendMessage(ctx context.Context, req *messagingpb.SendMessageRequest) (*messagingpb.SendMessageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := mc.client.SendMessage(ctx, req)
	if err != nil {
		log.Printf("gRPC call to SendMessage failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (mc *MessagingClient) Close() error {
	return mc.conn.Close()
}
