package domain

import (
	"strings"
	"time"

	"github.com/m1thrandir225/meridian/pkg/common"
)

const API_TOKEN_BYTES = 32

// Integration is the aggregate root for the integration domain
// It represents a single integration with a service, such as a bot or a webhook
type Integration struct {
	ID               IntegrationID
	ServiceName      string
	CreatorUserID    UserIDRef
	HashedAPIToken   APIToken
	TokenLookupHash  string
	CreatedAt        time.Time
	IsRevoked        bool
	TargetChannelIDs []ChannelIDRef
	Version          int64
	events           []common.DomainEvent
}

func NewIntegration(serviceName string, creatorID UserIDRef, targetChannels []ChannelIDRef, apiToken APIToken, rawToken string) (*Integration, error) {
	if strings.TrimSpace(serviceName) == "" {
		return nil, ErrServiceNameEmpty
	}
	if len(targetChannels) == 0 {
		return nil, ErrNoTargetChannels
	}
	if creatorID == "" {
		return nil, ErrCreatorIDEmpty
	}
	id, err := NewIntegrationID()
	if err != nil {
		return nil, err
	}

	integration := &Integration{
		ID:               *id,
		ServiceName:      serviceName,
		CreatorUserID:    creatorID,
		HashedAPIToken:   apiToken,
		CreatedAt:        time.Now(),
		TokenLookupHash:  GenerateLookupHash(rawToken),
		IsRevoked:        false,
		TargetChannelIDs: targetChannels,
		Version:          1,
		events:           make([]common.DomainEvent, 0),
	}

	integration.addEvent(CreateIntegrationRegisteredEvent(integration))

	return integration, nil
}

func (i *Integration) addEvent(event common.DomainEvent) {
	i.events = append(i.events, event)
}

func (i *Integration) Events() []common.DomainEvent {
	return i.events
}

func (i *Integration) ClearEvents() {
	i.events = nil
}

func (i *Integration) Revoke() error {
	if i.IsRevoked {
		return ErrIntegrationRevoked
	}
	i.IsRevoked = true
	i.Version++

	i.addEvent(CreateAPITokenRevokedEvent(i))
	return nil
}

func (i *Integration) UpdateTargetChannels(newChannels []ChannelIDRef) error {
	if len(newChannels) == 0 {
		return ErrNoTargetChannels
	}
	i.TargetChannelIDs = newChannels
	i.Version++

	i.addEvent(CreateIntegrationTargetChannelsUpdatedEvent(i))
	return nil
}

func (i *Integration) TargetChannelIDsAsStringSlice() []string {
	strs := make([]string, len(i.TargetChannelIDs))
	for j, ref := range i.TargetChannelIDs {
		strs[j] = string(ref)
	}
	return strs
}
