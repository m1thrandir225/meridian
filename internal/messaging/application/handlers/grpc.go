package handlers

import (
	"context"
	"log"
	"net"

	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	messagingpb "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/api"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	messagingService *services.ChannelService
	messagingpb.UnimplementedMessagingServiceServer
}

func NewGRPCHandler(service *services.ChannelService) *GRPCHandler {
	return &GRPCHandler{
		messagingService: service,
	}
}

func (h *GRPCHandler) SendMessage(ctx context.Context, req *messagingpb.SendMessageRequest) (*messagingpb.SendMessageResponse, error) {
	return nil, nil
}

func (h *GRPCHandler) GetMessages(ctx context.Context, req *messagingpb.GetMessagesRequest) (*messagingpb.GetMessagesResponse, error) {
	return nil, nil
}

func StartGRPCServer(port string, service *services.ChannelService) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(service)
	messagingpb.RegisterMessagingServiceServer(s, grpcHandler)

	log.Printf("Messaging gRPC server listening on port %s", port)
	return s.Serve(lis)
}
