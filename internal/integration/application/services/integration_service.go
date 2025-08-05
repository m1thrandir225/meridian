package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/m1thrandir225/meridian/internal/integration/domain"
	"github.com/m1thrandir225/meridian/internal/integration/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/kafka"
)

type IntegrationService struct {
	repo           persistence.IntegrationRepository
	tokenGenerator TokenGenerator
	publisher      kafka.EventPublisher
}

func NewIntegrationService(repo persistence.IntegrationRepository, tokenGenerator TokenGenerator, publisher kafka.EventPublisher) *IntegrationService {
	return &IntegrationService{
		repo:           repo,
		tokenGenerator: tokenGenerator,
		publisher:      publisher,
	}
}

func (s *IntegrationService) RegisterIntegration(ctx context.Context, cmd domain.RegisterIntegrationCommand) (*domain.Integration, string, error) {
	creatorUserID := domain.UserIDRef(cmd.CreatorUserID)
	targetChannels := make([]domain.ChannelIDRef, len(cmd.TargetChannels))
	for i, channelID := range cmd.TargetChannels {
		targetChannels[i] = domain.ChannelIDRef(channelID)
	}

	rawToken, hashedToken, err := s.tokenGenerator.Generate()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	integration, err := domain.NewIntegration(cmd.ServiceName, creatorUserID, targetChannels, hashedToken, rawToken)
	if err != nil {
		return nil, "", fmt.Errorf("registration validation failed: %w", err)
	}

	if err := s.repo.Save(ctx, integration); err != nil {
		return nil, "", fmt.Errorf("failed to save integration: %w", err)
	}

	s.dispatchEvents(ctx, integration)
	log.Printf("Successfully registered Integration %s", integration.ID.String())
	return integration, rawToken, nil
}

func (s *IntegrationService) ValidateApiToken(ctx context.Context, rawToken string) (isValid bool, integrationID string, channels []string, err error) {
	if rawToken == "" {
		return false, "", nil, errors.New("raw token cannot be empty")
	}

	integration, err := s.repo.FindByToken(ctx, rawToken)
	if err != nil {
		if errors.Is(err, domain.ErrIntegrationNotFound) {
			return false, "", nil, nil
		}
		return false, "", nil, fmt.Errorf("error during token lookup: %w", err)
	}

	if integration.IsRevoked {
		return false, integration.ID.String(), nil, nil
	}

	return true, integration.ID.String(), integration.TargetChannelIDsAsStringSlice(), nil
}

func (s *IntegrationService) RevokeToken(ctx context.Context, cmd domain.RevokeTokenCommand) error {
	integrationID, err := domain.NewIntegrationIDFromString(cmd.IntegrationID)
	if err != nil {
		return err
	}
	integration, err := s.repo.FindByID(ctx, integrationID.Value())
	if err != nil {
		return err
	}

	requestorId := domain.UserIDRef(cmd.RequestorID)

	if integration.CreatorUserID != requestorId {
		return domain.ErrForbidden
	}

	if err := integration.Revoke(); err != nil {
		if errors.Is(err, domain.ErrIntegrationRevoked) {
			return nil
		}
		return fmt.Errorf("revoke operation failed: %w", err)
	}

	if err := s.repo.Save(ctx, integration); err != nil {
		return fmt.Errorf("failed to save revoked integration: %w", err)
	}

	s.dispatchEvents(ctx, integration)
	log.Printf("Successfully revoked token for integration %s", integration.ID.String())
	return nil
}

func (s *IntegrationService) dispatchEvents(ctx context.Context, integration *domain.Integration) {
	if err := s.publisher.PublishEvents(ctx, integration.Events()); err != nil {
		log.Printf("CRITICAL: Failed to publish domain events %+v: %v", integration.Events(), err)
	}
	integration.ClearEvents()
}
