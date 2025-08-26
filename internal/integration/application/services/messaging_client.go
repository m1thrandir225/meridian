package services

import (
	"context"
	"fmt"
	"log"

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
	resp, err := mc.client.SendMessage(ctx, req)
	if err != nil {
		log.Printf("gRPC call to SendMessage failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (mc *MessagingClient) RegisterBot(ctx context.Context, req *messagingpb.RegisterBotRequest) (*messagingpb.RegisterBotResponse, error) {
	resp, err := mc.client.RegisterBot(ctx, req)
	if err != nil {
		log.Printf("gRPC call to RegisterBot failed: %v", err)
		return nil, err
	}

	if resp.Error != "" || !resp.Success {
		return nil, fmt.Errorf("failed to register bot: %s", resp.Error)
	}

	return resp, nil
}

func (mc *MessagingClient) RemoveBot(ctx context.Context, req *messagingpb.RemoveBotRequest) (*messagingpb.RemoveBotResponse, error) {
	resp, err := mc.client.RemoveBot(ctx, req)
	if err != nil {
		log.Printf("gRPC call to RemoveBot failed: %v", err)
		return nil, err
	}

	if resp.Error != "" || !resp.Success {
		return nil, fmt.Errorf("failed to remove bot: %s", resp.Error)
	}

	return resp, nil
}

func (mc *MessagingClient) Close() error {
	return mc.conn.Close()
}
