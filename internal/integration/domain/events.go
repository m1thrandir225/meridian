package domain

import "time"

type IntegrationRegisteredEvent struct {
	IntegrationID    string    `json:"integrationId"`
	ServiceName      string    `json:"serviceName"`
	CreatorUserID    string    `json:"creatorUserId"`
	TargetChannelIDs []string  `json:"targetChannelIds"`
	RegisteredAt     time.Time `json:"registeredAt"`
}

type APITokenRevokedEvent struct {
	IntegrationID string    `json:"integrationId"`
	RevokedAt     time.Time `json:"revokedAt"`
}

type IntegrationTargetChannelsUpdatedEvent struct {
	IntegrationID         string    `json:"integrationId"`
	UpdatedTargetChannels []string  `json:"updatedTargetChannels"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
