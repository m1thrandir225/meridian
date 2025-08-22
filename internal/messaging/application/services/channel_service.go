package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"

	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type ChannelService struct {
	repo              persistence.ChannelRepository
	eventPub          kafka.EventPublisher
	identityClient    *IdentityClient
	integrationClient *IntegrationClient
	logger            *logging.Logger
}

func NewChannelService(repo persistence.ChannelRepository, eventPub kafka.EventPublisher, identityClient *IdentityClient, integrationClient *IntegrationClient, logger *logging.Logger) *ChannelService {
	return &ChannelService{
		repo:              repo,
		eventPub:          eventPub,
		identityClient:    identityClient,
		integrationClient: integrationClient,
		logger:            logger,
	}
}

// HandleGetUserChannels returns the channels for a user
func (s *ChannelService) HandleGetUserChannels(ctx context.Context, cmd domain.GetUserChannelsCommand) ([]*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleGetUserChannels")
	logger.Info("Getting user channels")

	channels, err := s.repo.FindUserChannels(ctx, cmd.UserID)
	if err != nil {
		logger.Error("Failed to get user channels", zap.Error(err))
		return nil, err
	}

	logger.Info("User channels retrieved", zap.Int("count", len(channels)))
	return channels, nil
}

// HandleCreateChannel creates a new channel and publishes the events
func (s *ChannelService) HandleCreateChannel(ctx context.Context, cmd domain.CreateChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleCreateChannel")
	logger.Info("Creating channel")

	channel, err := domain.NewChannel(cmd.Name, cmd.Topic, cmd.CreatorUserID)
	if err != nil {
		logger.Error("Failed to create channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Channel created", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// HandleGetChannel returns the channel and publishes the events
func (s *ChannelService) HandleGetChannel(ctx context.Context, cmd domain.GetChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleGetChannel")
	logger.Info("Getting channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Channel retrieved", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// HandleAddBotToChannel adds a bot to a channel and publishes the events
func (s *ChannelService) HandleAddBotToChannel(ctx context.Context, cmd domain.AddBotToChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleAddBotToChannel")
	logger.Info("Adding bot to channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.AddBotMember(cmd.IntegrationID)
	if err != nil {
		logger.Error("Failed to add bot to channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Bot added to channel", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// HandleJoinChannel joins a channel and publishes the events
func (s *ChannelService) HandleJoinChannel(ctx context.Context, cmd domain.JoinChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleJoinChannel")
	logger.Info("Joining channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.AddMember(cmd.UserID)
	if err != nil {
		logger.Error("Failed to add member to channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Joined channel", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// HandleLeaveChannel leaves a channel
func (s *ChannelService) HandleLeaveChannel(ctx context.Context, cmd domain.LeaveChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleLeaveChannel")
	logger.Info("Leaving channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.RemoveMember(cmd.UserID)
	if err != nil {
		logger.Error("Failed to remove member from channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Left channel", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// HandleSetChannelTopic sets the topic of a channel
func (s *ChannelService) HandleSetChannelTopic(ctx context.Context, cmd domain.SetChannelTopicCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleSetChannelTopic")
	logger.Info("Setting channel topic")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	channel.SetTopic(cmd.UserID, cmd.Topic)

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Channel topic set", zap.String("channel_id", channel.ID.String()))
	return channel, err
}

// HandleArchiveChannel archives a channel
func (s *ChannelService) HandleArchiveChannel(ctx context.Context, cmd domain.ArchiveChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleArchiveChannel")
	logger.Info("Archiving channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.ArchiveChannel(cmd.UserID)
	if err != nil {
		logger.Error("Failed to archive channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Channel archived", zap.String("channel_id", channel.ID.String()))
	return channel, err
}

// HandleUnarchiveChannel unarchives a channel
func (s *ChannelService) HandleUnarchiveChannel(ctx context.Context, cmd domain.UnarchiveChannelCommand) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleUnarchiveChannel")
	logger.Info("Unarchiving channel")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.UnarchiveChannel(cmd.UserID)
	if err != nil {
		logger.Error("Failed to unarchive channel", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	err = s.eventPub.PublishEvents(ctx, channel.GetPendingEvents())
	if err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Channel unarchived", zap.String("channel_id", channel.ID.String()))
	return channel, err
}

// HandleCreateChannelInvite creates a new channel invite and publishes the events
func (s *ChannelService) HandleCreateChannelInvite(
	ctx context.Context,
	cmd domain.CreateChannelInviteCommand,
) (*domain.Channel, *domain.ChannelInvite, error) {
	logger := s.logger.WithMethod("HandleCreateChannelInvite")
	logger.Info("Creating channel invite")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, nil, err
	}

	invite, err := channel.CreateInvite(
		cmd.CreatedByUserID,
		cmd.ExpiresAt,
		cmd.MaxUses,
	)
	if err != nil {
		logger.Error("Failed to create channel invite", zap.Error(err))
		return nil, nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, nil, err
	}

	if err := s.eventPub.PublishEvents(ctx, channel.GetPendingEvents()); err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, nil, err
	}
	channel.ClearPendingEvents()

	if err != nil {
		logger.Error("Failed to create channel invite", zap.Error(err))
		return nil, nil, err
	}

	logger.Info("Channel invite created", zap.String("channel_id", channel.ID.String()))
	return channel, invite, nil
}

func (s *ChannelService) HandleAcceptChannelInvite(
	ctx context.Context,
	cmd domain.AcceptChannelInviteCommand,
) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleAcceptChannelInvite")
	logger.Info("Accepting channel invite")

	channel, err := s.repo.FindByInviteCode(ctx, cmd.InviteCode)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.AcceptInvite(cmd.InviteCode, cmd.UserID)
	if err != nil {
		logger.Error("Failed to accept invite", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	if err := s.eventPub.PublishEvents(ctx, channel.GetPendingEvents()); err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}
	channel.ClearPendingEvents()

	logger.Info("Invite accepted", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

func (s *ChannelService) HandleGetChannelInvites(
	ctx context.Context,
	cmd domain.GetChannelInvitesCommand,
) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleGetChannelInvites")
	logger.Info("Getting channel invites")

	channel, err := s.repo.FindById(ctx, cmd.ChannelID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	logger.Info("Channel invites retrieved", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

func (s *ChannelService) HandleDeactivateChannelInvite(
	ctx context.Context,
	cmd domain.DeactivateChannelInviteCommand,
) (*domain.Channel, error) {
	logger := s.logger.WithMethod("HandleDeactivateChannelInvite")
	logger.Info("Deactivating channel invite")

	channel, err := s.repo.FindByInviteID(ctx, cmd.InviteID)
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	err = channel.DeactivateInvite(cmd.InviteID, cmd.UserID)
	if err != nil {
		logger.Error("Failed to deactivate invite", zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, channel); err != nil {
		logger.Error("Failed to save channel", zap.Error(err))
		return nil, err
	}

	if err := s.eventPub.PublishEvents(ctx, channel.GetPendingEvents()); err != nil {
		logger.Error("Failed to publish events", zap.Error(err))
		return nil, err
	}

	channel.ClearPendingEvents()

	logger.Info("Invite deactivated", zap.String("channel_id", channel.ID.String()))
	return channel, nil
}

// getChannelMembers returns the users and integration bots for a channel
func (s *ChannelService) getChannelMembers(ctx context.Context, channel *domain.Channel) ([]*domain.User, []*domain.IntegrationBot, error) {
	logger := s.logger.WithMethod("getChannelMembers")
	logger.Info("Getting channel members")

	userIDs := make([]string, 0)
	integrationIDs := make([]string, 0)

	for _, member := range channel.Members {
		if member.GetRole() == "bot" {
			integrationIDs = append(integrationIDs, member.GetId().String())
		} else {
			userIDs = append(userIDs, member.GetId().String())
		}
	}

	logger.Info("Channel members retrieved", zap.Int("user_count", len(userIDs)), zap.Int("integration_count", len(integrationIDs)))

	var users []*domain.User
	var integrationBots []*domain.IntegrationBot

	// Fetch users
	if len(userIDs) > 0 {
		usersResp, err := s.identityClient.GetUsers(ctx, userIDs)
		if err != nil {
			logger.Error("Failed to fetch user information", zap.Error(err))
			return nil, nil, fmt.Errorf("failed to fetch user information: %w", err)
		}

		for _, user := range usersResp.Users {
			userID, err := uuid.Parse(user.Id)
			if err != nil {
				logger.Error("Failed to parse user ID", zap.Error(err))
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
			logger.Error("Failed to fetch integration information", zap.Error(err))
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
				logger.Error("Failed to parse integration ID", zap.Error(err))
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

	logger.Info("Channel members retrieved", zap.Int("user_count", len(users)), zap.Int("integration_count", len(integrationBots)))
	return users, integrationBots, nil
}

// ReturnChannelDTO returns a channel DTO
func (s *ChannelService) ReturnChannelDTO(ctx context.Context, channel *domain.Channel) (*domain.ChannelDTO, error) {
	logger := s.logger.WithMethod("ReturnChannelDTO")
	logger.Info("Returning channel DTO")

	users, integrationBots, err := s.getChannelMembers(ctx, channel)
	if err != nil {
		logger.Error("Failed to get channel members", zap.Error(err))
		return nil, err
	}

	dto := domain.ToChannelDTO(channel, users, integrationBots)
	return &dto, nil
}

// ReturnChannelDTOs returns a list of channel DTOs
func (s *ChannelService) ReturnChannelDTOs(ctx context.Context, channels []*domain.Channel) ([]domain.ChannelDTO, error) {
	logger := s.logger.WithMethod("ReturnChannelDTOs")
	logger.Info("Returning channel DTOs")

	dtos := make([]domain.ChannelDTO, len(channels))
	for i, channel := range channels {
		dto, err := s.ReturnChannelDTO(ctx, channel)
		if err != nil {
			logger.Error("Failed to return channel DTO", zap.Error(err))
			return nil, err
		}
		dtos[i] = *dto
	}
	logger.Info("Channel DTOs returned", zap.Int("count", len(dtos)))
	return dtos, nil
}
