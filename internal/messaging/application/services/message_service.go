package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/kafka"
)

type MessageService struct {
	repo              persistence.ChannelRepository
	eventPub          kafka.EventPublisher
	identityClient    *IdentityClient
	integrationClient *IntegrationClient
}

func NewMessageService(repo persistence.ChannelRepository, eventPub kafka.EventPublisher, identityClient *IdentityClient, integrationClient *IntegrationClient) *MessageService {
	return &MessageService{
		repo:              repo,
		eventPub:          eventPub,
		identityClient:    identityClient,
		integrationClient: integrationClient,
	}
}

func (s *MessageService) HandleListMessages(ctx context.Context, cmd domain.ListMessagesForChannelCommand) ([]domain.Message, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	messages, err := s.repo.FindMessages(context.Background(), cmd.ChannelID, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, err
	}

	channel.Messages = messages

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()

	return messages, nil
}

func (s *MessageService) HandleMessageSent(ctx context.Context, cmd domain.SendMessageCommand) (*domain.Message, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	//If the message is a reply, we need to load the messages in the domain
	if cmd.ParentMessageID != nil {
		messages, err := s.repo.FindMessages(ctx, cmd.ChannelID, 1000, 0)
		if err != nil {
			return nil, fmt.Errorf("error finding messages: %w", err)
		}
		channel.Messages = messages
	}

	message, err := channel.PostMessage(cmd.SenderUserID, cmd.Content, cmd.ParentMessageID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveMessage(ctx, message); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return message, err
}

// HandleNotificationSent sends a notification to a channel
// Might be redundant, but keeping it for now
// TODO: Remove this if it's redundant
func (s *MessageService) HandleNotificationSent(ctx context.Context, cmd domain.SendNotificationCommand) (*domain.Message, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}
	message, err := channel.PostNotification(cmd.IntegrationID, cmd.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveMessage(ctx, message); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()

	return message, err
}

// HandleAddReaction adds a reaction to a message
func (s *MessageService) HandleAddReaction(ctx context.Context, cmd domain.AddReactionCommand) (*domain.Reaction, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	// FIXME:  should the channel return all the messages??
	messages, err := s.repo.FindMessages(ctx, cmd.ChannelID, 100, 0)
	if err != nil {
		return nil, err
	}
	channel.Messages = messages

	newReaction, err := channel.AddReaction(cmd.MessageID, cmd.UserID, cmd.ReactionType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveReaction(ctx, newReaction); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return newReaction, nil
}

// HandleRemoveReaction removes a reaction from a message
func (s *MessageService) HandleRemoveReaction(ctx context.Context, cmd domain.RemoveReactionCommand) (*domain.Reaction, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	messages, err := s.repo.FindMessages(ctx, cmd.ChannelID, 100, 0)
	if err != nil {
		return nil, err
	}

	channel.Messages = messages

	reaction, err := channel.RemoveReaction(cmd.MessageID, cmd.UserID, cmd.ReactionType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteReaction(ctx, cmd.MessageID, cmd.UserID, cmd.ReactionType); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return reaction, nil
}

func (s *MessageService) ToMessageDTOs(ctx context.Context, messages []domain.Message) ([]domain.MessageDTO, error) {
	dtos := make([]domain.MessageDTO, len(messages))
	for i, message := range messages {
		dto, err := s.ToMessageDTO(ctx, &message)
		if err != nil {
			return nil, err
		}
		dtos[i] = *dto
	}
	return dtos, nil
}

// ToMessageDTO returns the message as a DTO
func (s *MessageService) ToMessageDTO(ctx context.Context, message *domain.Message) (*domain.MessageDTO, error) {
	senderUserID := message.GetSenderUserId()
	integrationID := message.GetIntegrationId()

	if senderUserID != nil {
		user, err := s.getSenderUser(ctx, senderUserID.String())
		if err != nil {
			return nil, err
		}
		dto := domain.ToMessageDTO(message, user, nil)
		return &dto, nil
	}
	if integrationID != nil {
		integrationBot, err := s.getSenderIntegrationBot(ctx, integrationID.String())
		if err != nil {
			return nil, err
		}
		dto := domain.ToMessageDTO(message, nil, integrationBot)
		return &dto, nil
	}

	return nil, fmt.Errorf("message has no sender user or integration id")
}

// getSenderUser returns the user with information from the identity service
func (s *MessageService) getSenderUser(ctx context.Context, userID string) (*domain.User, error) {
	pbUser, err := s.identityClient.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(pbUser.User.Id)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(userId, pbUser.User.Username, pbUser.User.FirstName, pbUser.User.LastName, pbUser.User.Email)

	return user, nil
}

// getSenderIntegrationBot returns the integration bot with information from the integration service
func (s *MessageService) getSenderIntegrationBot(ctx context.Context, id string) (*domain.IntegrationBot, error) {
	pbIntegration, err := s.integrationClient.GetIntegration(ctx, id)
	if err != nil {
		return nil, err
	}

	if pbIntegration == nil {
		return nil, fmt.Errorf("integration response is nil")
	}

	if pbIntegration.Integration == nil {
		return nil, fmt.Errorf("integration data is nil")
	}

	if pbIntegration.Integration.Id == "" {
		return nil, fmt.Errorf("integration ID is empty")
	}

	integrationID, err := uuid.Parse(pbIntegration.Integration.Id)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse(time.RFC3339, pbIntegration.Integration.CreatedAt)
	if err != nil {
		return nil, err
	}

	integrationBot := domain.NewIntegrationBot(integrationID, pbIntegration.Integration.ServiceName, createdAt, pbIntegration.Integration.IsRevoked)

	return integrationBot, nil
}
