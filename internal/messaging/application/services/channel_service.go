package services

import "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"

type ChannelService struct {
	repo persistence.ChannelRepository
}

func NewChannelService(repo persistence.ChannelRepository) *ChannelService {
	return &ChannelService{
		repo: repo,
	}
}
