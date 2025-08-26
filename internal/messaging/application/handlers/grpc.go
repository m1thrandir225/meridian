package handlers

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	messagingpb "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	channelService *services.ChannelService
	messageService *services.MessageService
	messagingpb.UnimplementedMessagingServiceServer
	wsHandler *WebSocketHandler
	cache     *cache.RedisCache
	logger    *logging.Logger
}

func NewGRPCHandler(
	channelService *services.ChannelService,
	messageService *services.MessageService,
	wsHandler *WebSocketHandler,
	cache *cache.RedisCache,
	logger *logging.Logger,
) *GRPCServer {
	return &GRPCServer{
		channelService: channelService,
		messageService: messageService,
		wsHandler:      wsHandler,
		cache:          cache,
		logger:         logger,
	}
}

func (h *GRPCServer) SendMessage(ctx context.Context, req *messagingpb.SendMessageRequest) (*messagingpb.SendMessageResponse, error) {
	logger := h.logger.WithMethod("SendMessage")
	logger.Info("Sending message")

	if req.Content == "" {
		logger.Error("Content is required")
		return nil, fmt.Errorf("content is required")
	}

	if len(req.TargetChannelIds) == 0 {
		logger.Error("At least one target channel is required")
		return nil, fmt.Errorf("at least one target channel is required")
	}

	var messages []*domain.Message
	var responses []*messagingpb.MessageResponse

	for _, channelIDStr := range req.TargetChannelIds {
		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			logger.Error("Invalid channel ID", zap.String("channel_id", channelIDStr), zap.Error(err))
			return nil, fmt.Errorf("invalid channel ID %s: %w", channelIDStr, err)
		}

		messageContent := domain.NewMessageContent(req.Content)

		var cmd interface{}
		if req.SenderType == "integration" {
			integrationID, err := uuid.Parse(req.SenderId)
			if err != nil {
				logger.Error("Invalid integration ID", zap.String("integration_id", req.SenderId), zap.Error(err))
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
				logger.Error("Invalid user ID", zap.String("user_id", req.SenderId), zap.Error(err))
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
			logger.Error("Failed to send message to channel", zap.String("channel_id", channelIDStr), zap.Error(err))
			return nil, fmt.Errorf("failed to send message to channel %s: %w", channelIDStr, err)
		}

		if message == nil {
			logger.Error("Failed to send message to channel", zap.String("channel_id", channelIDStr), zap.Error(fmt.Errorf("message is nil")))
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
		logger.Info("Message sent to channel", zap.String("channel_id", channelIDStr), zap.String("message_id", message.GetId().String()))
	}

	return &messagingpb.SendMessageResponse{
		Success:   true,
		Responses: responses,
	}, nil

}

func (h *GRPCServer) RegisterBot(ctx context.Context, req *messagingpb.RegisterBotRequest) (*messagingpb.RegisterBotResponse, error) {
	logger := h.logger.WithMethod("RegisterBot")
	logger.Info("Registering bot")

	integrationID, err := uuid.Parse(req.IntegrationId)
	if err != nil {
		logger.Error("Invalid integration ID", zap.String("integration_id", req.IntegrationId), zap.Error(err))
		return nil, fmt.Errorf("invalid integration ID %s: %w", req.IntegrationId, err)
	}

	requestorID, err := uuid.Parse(req.RequestorId)
	if err != nil {
		logger.Error("Invalid requestor ID", zap.String("requestor_id", req.RequestorId), zap.Error(err))
		return nil, fmt.Errorf("invalid requestor ID %s: %w", req.RequestorId, err)
	}

	channels := make([]*domain.Channel, 0, len(req.ChannelIds))
	success := true
	for _, channelIDStr := range req.ChannelIds {
		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			logger.Error("Invalid channel ID", zap.String("channel_id", channelIDStr), zap.Error(err))
			return nil, fmt.Errorf("invalid channel ID %s: %w", channelIDStr, err)
		}

		cmd := domain.AddBotToChannelCommand{
			ChannelID:     channelID,
			IntegrationID: integrationID,
			RequestorID:   requestorID,
		}

		channel, err := h.channelService.HandleAddBotToChannel(ctx, cmd)
		if err != nil {
			logger.Error("Failed to add bot to channel", zap.String("channel_id", channelIDStr), zap.Error(err))
			success = false
			continue
		}
		channels = append(channels, channel)
	}

	if !success {
		logger.Error("Failed to add bot to some channels")
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
		logger.Info("Invalidated cache for channel", zap.String("channel_id", channelIDStr))
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

	logger.Info("Bot registered successfully")
	return response, nil

}

func (h *GRPCServer) RemoveBot(ctx context.Context, req *messagingpb.RemoveBotRequest) (*messagingpb.RemoveBotResponse, error) {
	logger := h.logger.WithMethod("RemoveBot")
	logger.Info("Removing bot")

	integrationID, err := uuid.Parse(req.IntegrationId)
	if err != nil {
		logger.Error("Invalid integration ID", zap.String("integration_id", req.IntegrationId), zap.Error(err))
		return nil, fmt.Errorf("invalid integration ID %s: %w", req.IntegrationId, err)
	}
	requestorID, err := uuid.Parse(req.RequestorId)
	if err != nil {
		logger.Error("Invalid requestor ID", zap.String("requestor_id", req.RequestorId), zap.Error(err))
		return nil, fmt.Errorf("invalid requestor ID %s: %w", req.RequestorId, err)
	}

	success := true
	for _, channelIDStr := range req.ChannelIds {
		channelID, err := uuid.Parse(channelIDStr)
		if err != nil {
			logger.Error("Invalid channel ID", zap.String("channel_id", channelIDStr), zap.Error(err))
			return nil, fmt.Errorf("invalid channel ID %s: %w", channelIDStr, err)
		}

		cmd := domain.RemoveBotFromChannelCommand{
			ChannelID:     channelID,
			IntegrationID: integrationID,
			RequestorID:   requestorID,
		}

		_, err = h.channelService.HandleRemoveBotFromChannel(ctx, cmd)
		if err != nil {
			logger.Error("Failed to remove bot from channel", zap.String("channel_id", channelIDStr), zap.Error(err))
			success = false
		} else {
			channelCacheKey := fmt.Sprintf("channel:%s", channelIDStr)
			h.cache.Delete(ctx, channelCacheKey)
			logger.Info("Invalidated cache for channel", zap.String("channel_id", channelIDStr))
		}
	}

	if !success {
		return &messagingpb.RemoveBotResponse{
			Success:       false,
			IntegrationId: req.IntegrationId,
			ChannelIds:    req.ChannelIds,
			Error:         "failed to remove bot from some channels",
		}, nil
	}

	return &messagingpb.RemoveBotResponse{
		Success:       true,
		IntegrationId: req.IntegrationId,
		ChannelIds:    req.ChannelIds,
	}, nil
}

func StartGRPCServer(
	port string,
	channelService *services.ChannelService,
	messageService *services.MessageService,
	wsHandler *WebSocketHandler,
	cache *cache.RedisCache,
	logger *logging.Logger,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcHandler := NewGRPCHandler(channelService, messageService, wsHandler, cache, logger)
	messagingpb.RegisterMessagingServiceServer(s, grpcHandler)

	logger.Info("Messaging gRPC server listening", zap.String("port", port))
	return s.Serve(lis)
}
