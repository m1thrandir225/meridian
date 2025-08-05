package domain

import "errors"

var (
	ErrIntegrationNotFound = errors.New("integration not found")
	ErrServiceNameEmpty    = errors.New("service name cannot be empty")
	ErrNoTargetChannels    = errors.New("integration must target at least one channel")
	ErrCreatorIDEmpty      = errors.New("creator user ID cannot be empty")
	ErrIntegrationRevoked  = errors.New("integration is already revoked")
	ErrTokenGenerationFail = errors.New("failed to generate API token")
	ErrInvalidUUIDFormat   = errors.New("invalid UUID format")
	ErrForbidden           = errors.New("action forbidden for the requesting user")
	ErrNotFound            = ErrIntegrationNotFound
)
