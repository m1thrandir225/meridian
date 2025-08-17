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
	"github.com/m1thrandir225/meridian/pkg/cache"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	channelService *services.ChannelService
	messageService *services.MessageService
	messagingpb.UnimplementedMessagingServiceServer
	wsHandler *WebSocketHandler
	cache     *cache.RedisCache
}

func NewGRPCHandler(
	channelService *services.ChannelService,
	messageService *services.MessageService,
	wsHandler *WebSocketHandler,
	cache *cache.RedisCache,
) *GRPCServer {
	return &GRPCServer{
		channelService: channelService,
		messageService: messageService,
		wsHandler:      wsHandler,
		cache:          cache,
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

func (h *GRPCServer) RegisterBot(ctx context.Context, req *messagingpb.RegisterBotRequest) (*messagingpb.RegisterBotResponse, error) {
	integrationID, err := uuid.Parse(req.IntegrationId)
	if err != nil {
		return nil, fmt.Errorf("invalid integration ID %s: %w", req.IntegrationId, err)
	}

	channels := make([]*domain.Channel, 0, len(req.ChannelIds))
	success := true
	for _, channelIDStr := range req.ChannelIds {
		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid channel ID %s: %w", channelIDStr, err)
		}

		cmd := domain.AddBotToChannelCommand{
			ChannelID:     channelID,
			IntegrationID: integrationID,
		}

		channel, err := h.channelService.HandleAddBotToChannel(ctx, cmd)
		if err != nil {
			log.Printf("failed to add bot to channel %s: %v", channelIDStr, err)
			success = false
			continue
		}
		channels = append(channels, channel)
	}

	if !success {
		return &messagingpb.RegisterBotResponse{
			Success:       false,
			IntegrationId: integrationID.String(),
			ChannelIds:    req.ChannelIds,
			Error:         "failed to add bot to some channels",
		}, nil
	}

	for _, channelIDStr := range req.ChannelIds {
		channelCacheKey := fmt.Sprintf("channel:%s", channelIDStr)
		h.cache.Delete(ctx, channelCacheKey)
		log.Printf("Invalidated cache for channel: %s", channelIDStr)
	}

	channelIDs := make([]string, 0, len(channels))
	for _, channel := range channels {
		channelIDs = append(channelIDs, channel.ID.String())
	}

	response := &messagingpb.RegisterBotResponse{
		Success:       true,
		IntegrationId: integrationID.String(),
		ChannelIds:    channelIDs,
	}

	return response, nil

}
func StartGRPCServer(
	port string,
	channelService *services.ChannelService,
	messageService *services.MessageService,
	wsHandler *WebSocketHandler,
	cache *cache.RedisCache,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(channelService, messageService, wsHandler, cache)
	messagingpb.RegisterMessagingServiceServer(s, grpcHandler)

	log.Printf("Messaging gRPC server listening on port %s", port)
	return s.Serve(lis)
}
