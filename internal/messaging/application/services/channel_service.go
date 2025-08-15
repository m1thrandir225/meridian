package services

import (
	"context"
	"fmt"

	"github.com/m1thrandir225/meridian/pkg/kafka"

	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type ChannelService struct {
	repo           persistence.ChannelRepository
	eventPub       kafka.EventPublisher
	identityClient *IdentityClient
}

func NewChannelService(repo persistence.ChannelRepository, eventPub kafka.EventPublisher, identityClient *IdentityClient) *ChannelService {
	return &ChannelService{
		repo:           repo,
		eventPub:       eventPub,
		identityClient: identityClient,
	}
}

func (s *ChannelService) HandleGetUserChannels(ctx context.Context, cmd domain.GetUserChannelsCommand) ([]*domain.Channel, error) {
	channels, err := s.repo.FindUserChannels(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *ChannelService) HandleCreateChannel(ctx context.Context, cmd domain.CreateChannelCommand) (*domain.Channel, error) {
	channel, err := domain.NewChannel(cmd.Name, cmd.Topic, cmd.CreatorUserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, err
}

func (s *ChannelService) HandleGetChannel(ctx context.Context, cmd domain.GetChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()

	return channel, nil
}

func (s *ChannelService) HandleListMessages(ctx context.Context, cmd domain.ListMessagesForChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	messages, err := s.repo.FindMessages(context.Background(), cmd.ChannelID, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, err
	}

	messages, err = s.enrichMessagesWithUserInfo(ctx, messages)
	if err != nil {
		return nil, err
	}

	channel.Messages = messages

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, nil
}

func (s *ChannelService) HandleJoinChannel(ctx context.Context, cmd domain.JoinChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = channel.AddMember(cmd.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}
	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, err
}

func (s *ChannelService) HandleLeaveChannel(ctx context.Context, cmd domain.LeaveChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = channel.RemoveMember(cmd.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, nil
}

func (s *ChannelService) HandleMessageSent(ctx context.Context, cmd domain.SendMessageCommand) (*domain.Message, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
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

func (s *ChannelService) HandleNotificationSent(ctx context.Context, cmd domain.SendNotificationCommand) (*domain.Message, error) {
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

func (s *ChannelService) HandleAddReaction(ctx context.Context, cmd domain.AddReactionCommand) (*domain.Reaction, error) {
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

func (s *ChannelService) HandleRemoveReaction(ctx context.Context, cmd domain.RemoveReactionCommand) (*domain.Reaction, error) {
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

func (s *ChannelService) HandleSetChannelTopic(ctx context.Context, cmd domain.SetChannelTopicCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	channel.SetTopic(cmd.UserID, cmd.Topic)

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, err
}

func (s *ChannelService) HandleArchiveChannel(ctx context.Context, cmd domain.ArchiveChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = channel.ArchiveChannel(cmd.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, err
}

func (s *ChannelService) HandleUnarchiveChannel(ctx context.Context, cmd domain.UnarchiveChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = channel.UnarchiveChannel(cmd.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, err
}

func (s *ChannelService) enrichMessagesWithUserInfo(ctx context.Context, messages []domain.Message) ([]domain.Message, error) {
	if len(messages) == 0 {
		return messages, nil
	}

	userIDs := make(map[string]bool)
	for _, msg := range messages {
		if msg.GetSenderUserId() != nil {
			userIDs[msg.GetSenderUserId().String()] = true
		}
	}

	var userIDList []string
	for userID := range userIDs {
		userIDList = append(userIDList, userID)
	}

	users, err := s.identityClient.GetUsers(ctx, userIDList)
	if err != nil {
		return messages, fmt.Errorf("failed to fetch user information: %w", err)
	}

	enrichedMessages := make([]domain.Message, len(messages))
	for i, msg := range messages {
		enrichedMessages[i] = msg
		if msg.GetSenderUserId() != nil {
			userID := msg.GetSenderUserId().String()
			for _, user := range users.Users {
				if user.Id == userID {
					user := domain.NewUser(user.Id, user.Username, user.Email, user.FirstName, user.LastName)
					msg.SetSenderUser(user)
					enrichedMessages[i] = msg
					break
				}
			}
		}
	}

	return enrichedMessages, nil
}
