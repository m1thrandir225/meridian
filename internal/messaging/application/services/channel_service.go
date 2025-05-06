package services

import (
	"context"

	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type ChannelService struct {
	repo     persistence.ChannelRepository
	eventPub EventPublisher
}

func NewChannelService(repo persistence.ChannelRepository, eventPub EventPublisher) *ChannelService {
	return &ChannelService{
		repo:     repo,
		eventPub: eventPub,
	}
}

func (s *ChannelService) HandleCreateChannel(ctx context.Context, cmd domain.CreateChannelCommand) (*domain.Channel, error) {
	channel, err := domain.NewChannel(cmd.Name, cmd.Topic, cmd.CreatorUserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}
	return channel, err
}

func (s *ChannelService) HandleGetChannel(ctx context.Context, cmd domain.GetChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

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

	channel.Messages = messages

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
	return message, err
}

func (s *ChannelService) HandleAddReaction(ctx context.Context, cmd domain.AddReactionCommand) (*domain.Reaction, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	newReaction, err := channel.AddReaction(cmd.MessageID, cmd.UserID, cmd.ReactionType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveReaction(ctx, newReaction); err != nil {
		return nil, err
	}
	return newReaction, nil
}

func (s *ChannelService) HandleRemoveReaction(ctx context.Context, cmd domain.RemoveReactionCommand) (*domain.Reaction, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}
	reaction, err := channel.RemoveReaction(cmd.MessageID, cmd.UserID, cmd.ReactionType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, reaction.GetId()); err != nil {
		return nil, err
	}

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
	return channel, err
}

func (s *ChannelService) HandleArchiveChannel(ctx context.Context, cmd domain.ArchiveChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	channel.ArchiveChannel(cmd.UserID)

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}
	return channel, err
}

func (s *ChannelService) HandleUnarchiveChannel(ctx context.Context, cmd domain.UnarchiveChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	channel.UnarchiveChannel(cmd.UserID)

	if err := s.repo.Save(ctx, channel); err != nil {
		return nil, err
	}
	return channel, err
}
