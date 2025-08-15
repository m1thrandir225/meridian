package handlers

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	messagingpb "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/api"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	channelService *services.ChannelService
	messageService *services.MessageService
	messagingpb.UnimplementedMessagingServiceServer
	wsHandler *WebSocketHandler
}

func NewGRPCHandler(channelService *services.ChannelService, messageService *services.MessageService, wsHandler *WebSocketHandler) *GRPCServer {
	return &GRPCServer{
		channelService: channelService,
		messageService: messageService,
		wsHandler:      wsHandler,
	}
}

func (h *GRPCServer) SendMessage(ctx context.Context, req *messagingpb.SendMessageRequest) (*messagingpb.SendMessageResponse, error) {
	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	if len(req.TargetChannelIds) == 0 {
		return nil, fmt.Errorf("at least one target channel is required")
	}

	var messages []*domain.Message
	var responses []*messagingpb.MessageResponse

	for _, channelIDStr := range req.TargetChannelIds {
		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid channel ID %s: %w", channelIDStr, err)
		}

		messageContent := domain.NewMessageContent(req.Content)

		var cmd interface{}
		if req.SenderType == "integration" {
			integrationID, err := uuid.Parse(req.SenderId)
			if err != nil {
				return nil, fmt.Errorf("invalid integration ID %s: %w", req.SenderId, err)
			}

			cmd = domain.SendNotificationCommand{
				ChannelID:     channelID,
				IntegrationID: integrationID,
				Content:       messageContent,
			}
		} else {
			userID, err := uuid.Parse(req.SenderId)
			if err != nil {
				return nil, fmt.Errorf("invalid user ID %s: %w", req.SenderId, err)
			}

			cmd = domain.SendMessageCommand{
				ChannelID:       channelID,
				SenderUserID:    userID,
				Content:         messageContent,
				ParentMessageID: nil,
			}
		}

		var message *domain.Message
		if notifCmd, ok := cmd.(domain.SendNotificationCommand); ok {
			message, err = h.messageService.HandleNotificationSent(ctx, notifCmd)
		} else if msgCmd, ok := cmd.(domain.SendMessageCommand); ok {
			message, err = h.messageService.HandleMessageSent(ctx, msgCmd)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to send message to channel %s: %w", channelIDStr, err)
		}

		if message == nil {
			return nil, fmt.Errorf("failed to send message to channel %s: message is nil", channelIDStr)
		}

		messages = append(messages, message)

		response := &messagingpb.MessageResponse{
			MessageId:      message.GetId().String(),
			Success:        true,
			MessageContent: message.GetContent().GetText(),
			TargetChannelIds: []string{
				channelID.String(),
			},
		}

		if message.GetSenderUserId() != nil {
			response.SenderId = message.GetSenderUserId().String()
			response.SenderType = "user"
		}

		if message.GetIntegrationId() != nil {
			response.SenderId = message.GetIntegrationId().String()
			response.SenderType = "integration"
		}

		responses = append(responses, response)

		if h.wsHandler != nil {
			h.wsHandler.BroadcastMessage(message)
		}
	}

	return &messagingpb.SendMessageResponse{
		Success:   true,
		Responses: responses,
	}, nil

}
func StartGRPCServer(port string, channelService *services.ChannelService, messageService *services.MessageService, wsHandler *WebSocketHandler) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(channelService, messageService, wsHandler)
	messagingpb.RegisterMessagingServiceServer(s, grpcHandler)

	log.Printf("Messaging gRPC server listening on port %s", port)
	return s.Serve(lis)
}
