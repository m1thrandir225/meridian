package handlers

import (
	"time"

	"github.com/m1thrandir225/meridian/internal/integration/domain"
)

type IntegrationDTO struct {
	ServiceName     string   `json:"service_name"`
	TargetChannels  []string `json:"target_channels"`
	Token           string   `json:"token"`
	TokenLookupHash string   `json:"token_lookup_hash"`
	IsRevoked       bool     `json:"is_revoked"`
	CreatedAt       string   `json:"created_at"`
	ID              string   `json:"id"`
}

func ToIntegrationDTO(integration *domain.Integration, token string) IntegrationDTO {
	return IntegrationDTO{
		ID:              integration.ID.String(),
		ServiceName:     integration.ServiceName,
		TargetChannels:  integration.TargetChannelIDsAsStringSlice(),
		Token:           token,
		TokenLookupHash: integration.TokenLookupHash,
		IsRevoked:       integration.IsRevoked,
		CreatedAt:       integration.CreatedAt.Format(time.RFC3339),
	}
}
