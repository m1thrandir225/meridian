package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
	"github.com/m1thrandir225/meridian/internal/integration/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type IntegrationService struct {
	repo           persistence.IntegrationRepository
	tokenGenerator TokenGenerator
	publisher      kafka.EventPublisher
	logger         *logging.Logger
}

func NewIntegrationService(
	repo persistence.IntegrationRepository,
	tokenGenerator TokenGenerator,
	publisher kafka.EventPublisher,
	logger *logging.Logger,
) *IntegrationService {
	return &IntegrationService{
		repo:           repo,
		tokenGenerator: tokenGenerator,
		publisher:      publisher,
		logger:         logger,
	}
}

func (s *IntegrationService) RegisterIntegration(ctx context.Context, cmd domain.RegisterIntegrationCommand) (*domain.Integration, string, error) {
	logger := s.logger.WithMethod("RegisterIntegration")
	logger.Info("Registering integration")

	creatorUserID := domain.UserIDRef(cmd.CreatorUserID)
	targetChannels := make([]domain.ChannelIDRef, len(cmd.TargetChannels))
	for i, channelID := range cmd.TargetChannels {
		targetChannels[i] = domain.ChannelIDRef(channelID)
	}

	rawToken, hashedToken, err := s.tokenGenerator.Generate()
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	integration, err := domain.NewIntegration(cmd.ServiceName, creatorUserID, targetChannels, *hashedToken, rawToken)
	if err != nil {
		logger.Error("Failed to create integration", zap.Error(err))
		return nil, "", fmt.Errorf("registration validation failed: %w", err)
	}

	if err := s.repo.Save(ctx, integration); err != nil {
		logger.Error("Failed to save integration", zap.Error(err))
		return nil, "", fmt.Errorf("failed to save integration: %w", err)
	}

	s.dispatchEvents(ctx, integration)
	logger.Info("Successfully registered Integration", zap.String("integration_id", integration.ID.String()))
	return integration, rawToken, nil
}

func (s *IntegrationService) ValidateApiToken(ctx context.Context, rawToken string) (isValid bool, integrationID string, channels []string, err error) {
	logger := s.logger.WithMethod("ValidateApiToken")
	logger.Info("Validating API token")

	if rawToken == "" {
		logger.Error("Raw token cannot be empty")
		return false, "", nil, errors.New("raw token cannot be empty")
	}

	tokenHash := domain.GenerateLookupHash(rawToken)

	integration, err := s.repo.FindByTokenLookupHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrIntegrationNotFound) {
			logger.Error("Integration not found", zap.String("token_hash", tokenHash))
			return false, "", nil, nil
		}
		logger.Error("Error during token lookup", zap.Error(err))
		return false, "", nil, fmt.Errorf("error during token lookup: %w", err)
	}

	if integration.IsRevoked {
		logger.Error("Error integration is revoked", zap.String("integration_id", integration.ID.String()))
		return false, integration.ID.String(), nil, nil
	}

	logger.Info("API token is valid", zap.String("integration_id", integration.ID.String()))
	return true, integration.ID.String(), integration.TargetChannelIDsAsStringSlice(), nil
}

func (s *IntegrationService) RevokeToken(ctx context.Context, cmd domain.RevokeTokenCommand) error {
	logger := s.logger.WithMethod("RevokeToken")
	logger.Info("Revoking token")

	integrationID, err := domain.NewIntegrationIDFromString(cmd.IntegrationID)
	if err != nil {
		logger.Error("Failed to parse integration ID", zap.Error(err))
		return err
	}
	integration, err := s.repo.FindByID(ctx, integrationID.Value())
	if err != nil {
		logger.Error("Failed to find integration", zap.Error(err))
		return err
	}

	requestorId := domain.UserIDRef(cmd.RequestorID)

	if integration.CreatorUserID != requestorId {
		logger.Error("Forbidden", zap.String("integration_id", integration.ID.String()))
		return domain.ErrForbidden
	}

	if err := integration.Revoke(); err != nil {
		if errors.Is(err, domain.ErrIntegrationRevoked) {
			logger.Error("Integration is already revoked", zap.String("integration_id", integration.ID.String()))
			return nil
		}
		logger.Error("Integration revoke operation failed", zap.Error(err))
		return fmt.Errorf("revoke operation failed: %w", err)
	}

	if err := s.repo.Save(ctx, integration); err != nil {
		logger.Error("Failed to save revoked integration", zap.Error(err))
		return fmt.Errorf("failed to save revoked integration: %w", err)
	}

	s.dispatchEvents(ctx, integration)
	logger.Info("Successfully revoked token for integration", zap.String("integration_id", integration.ID.String()))
	return nil
}

func (s *IntegrationService) GetIntegration(ctx context.Context, cmd domain.GetIntegrationCommand) (*domain.Integration, error) {
	logger := s.logger.WithMethod("GetIntegration")
	logger.Info("Getting integration")

	integrationID, err := domain.NewIntegrationIDFromString(cmd.IntegrationID)
	if err != nil {
		logger.Error("Failed to parse integration ID", zap.Error(err))
		return nil, err
	}

	integration, err := s.repo.FindByID(ctx, integrationID.Value())
	if err != nil {
		logger.Error("Failed to find integration", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully retrieved integration", zap.String("integration_id", integration.ID.String()))
	return integration, nil
}

func (s *IntegrationService) ListIntegrations(ctx context.Context, cmd domain.ListIntegrationsCommand) ([]*domain.Integration, error) {
	logger := s.logger.WithMethod("ListIntegrations")
	logger.Info("Listing integrations")

	if cmd.CreatorUserID == "" {
		logger.Error("Creator user ID cannot be empty")
		return nil, errors.New("creator user ID cannot be empty")
	}

	creatorUserID, err := uuid.Parse(cmd.CreatorUserID)
	if err != nil {
		logger.Error("Invalid creator user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid creator user ID: %w", err)
	}

	integrations, err := s.repo.FindByCreatorUserID(ctx, creatorUserID)
	if err != nil {
		logger.Error("Failed to find integrations", zap.Error(err))
		return nil, fmt.Errorf("failed to find integrations: %w", err)
	}

	logger.Info("Successfully retrieved integrations", zap.Int("count", len(integrations)))
	return integrations, nil
}

func (s *IntegrationService) UpdateIntegration(ctx context.Context, cmd domain.UpdateIntegrationCommand) (*domain.Integration, error) {
	logger := s.logger.WithMethod("UpdateIntegration")
	logger.Info("Updating integration")

	integrationID, err := domain.NewIntegrationIDFromString(cmd.IntegrationID)
	if err != nil {
		logger.Error("Failed to parse integration ID", zap.Error(err))
		return nil, fmt.Errorf("invalid integration ID: %w", err)
	}
	integration, err := s.repo.FindByID(ctx, integrationID.Value())
	if err != nil {
		logger.Error("Failed to find integration", zap.Error(err))
		return nil, err
	}

	requestorId, err := domain.NewUserIDRef(cmd.RequestorID)
	if err != nil {
		logger.Error("Failed to parse requestor ID", zap.Error(err))
		return nil, fmt.Errorf("invalid requestor ID: %w", err)
	}

	if integration.CreatorUserID != requestorId {
		logger.Error("Forbidden", zap.String("integration_id", integration.ID.String()))
		return nil, domain.ErrForbidden
	}
	if integration.IsRevoked {
		logger.Error("Integration is revoked", zap.String("integration_id", integration.ID.String()))
		return nil, domain.ErrIntegrationRevoked
	}

	targetChannels := make([]domain.ChannelIDRef, len(cmd.TargetChannelIDs))
	for i, channelID := range cmd.TargetChannelIDs {
		targetChannels[i] = domain.ChannelIDRef(channelID)
	}

	if err := integration.UpdateTargetChannels(targetChannels); err != nil {
		logger.Error("failed to update target channels", zap.Error(err))
		return nil, fmt.Errorf("failed to update target channels: %w", err)
	}

	if err := s.repo.Save(ctx, integration); err != nil {
		logger.Error("failed to save updated integration", zap.Error(err))
		return nil, fmt.Errorf("failed to save updated integration: %w", err)
	}

	s.dispatchEvents(ctx, integration)
	logger.Info("Successfully updated integration", zap.String("integration_id", integration.ID.String()))
	return integration, nil
}

func (s *IntegrationService) dispatchEvents(ctx context.Context, integration *domain.Integration) {
	if err := s.publisher.PublishEvents(ctx, integration.Events()); err != nil {
		log.Printf("CRITICAL: Failed to publish domain events %+v: %v", integration.Events(), err)
	}
	integration.ClearEvents()
}
