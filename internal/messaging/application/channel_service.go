package application

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
