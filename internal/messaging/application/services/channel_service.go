package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/kafka"

	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type ChannelService struct {
	repo              persistence.ChannelRepository
	eventPub          kafka.EventPublisher
	identityClient    *IdentityClient
	integrationClient *IntegrationClient
}

func NewChannelService(repo persistence.ChannelRepository, eventPub kafka.EventPublisher, identityClient *IdentityClient, integrationClient *IntegrationClient) *ChannelService {
	return &ChannelService{
		repo:              repo,
		eventPub:          eventPub,
		identityClient:    identityClient,
		integrationClient: integrationClient,
	}
}

// HandleGetUserChannels returns the channels for a user
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

	return channel, nil
}

// handleGetChannel returns the channel and publishes the events
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

// HandleAddBotToChannel adds a bot to a channel and publishes the events
func (s *ChannelService) HandleAddBotToChannel(ctx context.Context, cmd domain.AddBotToChannelCommand) (*domain.Channel, error) {
	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		return nil, err
	}

	err = channel.AddBotMember(cmd.IntegrationID)
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
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		return nil, err
	}
	channel.ClearPendingEvents()
	return channel, nil
}

// HandleLeaveChannel leaves a channel
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

// HandleSetChannelTopic sets the topic of a channel
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

// HandleArchiveChannel archives a channel
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

// HandleUnarchiveChannel unarchives a channel
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

// getChannelMembers returns the users and integration bots for a channel
func (s *ChannelService) getChannelMembers(ctx context.Context, channel *domain.Channel) ([]*domain.User, []*domain.IntegrationBot, error) {
	userIDs := make([]string, 0)
	integrationIDs := make([]string, 0)

	for _, member := range channel.Members {
		if member.GetRole() == "bot" {
			integrationIDs = append(integrationIDs, member.GetId().String())
		} else {
			userIDs = append(userIDs, member.GetId().String())
		}
	}

	var users []*domain.User
	var integrationBots []*domain.IntegrationBot

	// Fetch users
	if len(userIDs) > 0 {
		usersResp, err := s.identityClient.GetUsers(ctx, userIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch user information: %w", err)
		}

		for _, user := range usersResp.Users {
			userID, err := uuid.Parse(user.Id)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse user ID: %w", err)
			}
			domainUser := domain.NewUser(userID, user.Username, user.FirstName, user.LastName, user.Email)
			users = append(users, domainUser)
		}
	}

	// Fetch integration bots
	if len(integrationIDs) > 0 {
		integrations, err := s.integrationClient.GetIntegrations(ctx, integrationIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch integration information: %w", err)
		}

		for _, integration := range integrations {
			// Parse the created_at time
			createdAt, err := time.Parse(time.RFC3339, integration.CreatedAt)
			if err != nil {
				log.Printf("Failed to parse integration created_at time: %v", err)
				createdAt = time.Now()
			}

			integrationID, err := uuid.Parse(integration.Id)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse integration ID: %w", err)
			}

			integrationBot := domain.NewIntegrationBot(
				integrationID,
				integration.ServiceName,
				createdAt,
				integration.IsRevoked,
			)
			integrationBots = append(integrationBots, integrationBot)
		}
	}

	return users, integrationBots, nil
}

func (s *ChannelService) ReturnChannelDTO(ctx context.Context, channel *domain.Channel) (*domain.ChannelDTO, error) {
	users, integrationBots, err := s.getChannelMembers(ctx, channel)
	if err != nil {
		return nil, err
	}

	dto := domain.ToChannelDTO(channel, users, integrationBots)
	return &dto, nil
}

func (s *ChannelService) ReturnChannelDTOs(ctx context.Context, channels []*domain.Channel) ([]domain.ChannelDTO, error) {
	dtos := make([]domain.ChannelDTO, len(channels))
	for i, channel := range channels {
		dto, err := s.ReturnChannelDTO(ctx, channel)
		if err != nil {
			return nil, err
		}
		dtos[i] = *dto
	}
	return dtos, nil
}
