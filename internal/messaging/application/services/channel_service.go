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

// TODO: implement
func (s *ChannelService) HandleCreateChannel(ctx context.Context, cmd domain.CreateChannelCommand) (*domain.Channel, error) {
	return nil, nil
}

// TODO: implement
func (s *ChannelService) HandleJoinChannel(ctx context.Context, cmd domain.JoinChannelCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleLeaveChannel(ctx context.Context, cmd domain.LeaveChannelCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleMessageSent(ctx context.Context, cmd domain.SendMessageCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleNotificationSent(ctx context.Context, cmd domain.SendNotificationCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleAddReaction(ctx context.Context, cmd domain.AddReactionCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleRemoveReaction(ctx context.Context, cmd domain.RemoveReactionCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleSetChannelTopic(ctx context.Context, cmd domain.SetChannelTopicCommand) error {
	return nil
}

// TODO: implement
func (s *ChannelService) HandleArchiveChannel(ctx context.Context, cmd domain.ArchiveChannelCommand) error {
	return nil
}

// TODO:  implement
func (s *ChannelService) HandleUnarchiveChannel(ctx context.Context, cmd domain.UnarchiveChannelCommand) error {
	return nil
}
