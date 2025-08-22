package domain

import (
	"time"

	"github.com/m1thrandir225/meridian/pkg/common"
)

type IntegrationRegisteredEvent struct {
	common.BaseDomainEvent
	IntegrationID    string    `json:"integrationId"`
	ServiceName      string    `json:"serviceName"`
	CreatorUserID    string    `json:"creatorUserId"`
	TargetChannelIDs []string  `json:"targetChannelIds"`
	RegisteredAt     time.Time `json:"registeredAt"`
}

type APITokenRevokedEvent struct {
	common.BaseDomainEvent
	IntegrationID string    `json:"integrationId"`
	RevokedAt     time.Time `json:"revokedAt"`
}

type IntegrationTargetChannelsUpdatedEvent struct {
	common.BaseDomainEvent
	IntegrationID         string    `json:"integrationId"`
	UpdatedTargetChannels []string  `json:"updatedTargetChannels"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

func CreateIntegrationRegisteredEvent(integration *Integration) IntegrationRegisteredEvent {
	base := common.NewBaseDomainEvent("IntegrationRegistered", integration.ID.value, integration.Version, "Integration")

	return IntegrationRegisteredEvent{
		BaseDomainEvent:  base,
		IntegrationID:    integration.ID.String(),
		ServiceName:      integration.ServiceName,
		CreatorUserID:    integration.CreatorUserID.String(),
		TargetChannelIDs: integration.TargetChannelIDsAsStringSlice(),
		RegisteredAt:     integration.CreatedAt,
	}
}

func CreateAPITokenRevokedEvent(integration *Integration) APITokenRevokedEvent {
	base := common.NewBaseDomainEvent("APITokenRevoked", integration.ID.value, integration.Version, "Integration")

	return APITokenRevokedEvent{
		BaseDomainEvent: base,
		IntegrationID:   integration.ID.String(),
		RevokedAt:       time.Now(),
	}
}

func CreateIntegrationTargetChannelsUpdatedEvent(integration *Integration) IntegrationTargetChannelsUpdatedEvent {
	base := common.NewBaseDomainEvent("IntegrationTargetChannelsUpdated", integration.ID.value, integration.Version, "Integration")

	return IntegrationTargetChannelsUpdatedEvent{
		BaseDomainEvent:       base,
		IntegrationID:         integration.ID.String(),
		UpdatedTargetChannels: integration.TargetChannelIDsAsStringSlice(),
		UpdatedAt:             time.Now(),
	}
}
