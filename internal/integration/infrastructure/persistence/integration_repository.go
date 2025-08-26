package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
)

type IntegrationRepository interface {
	Save(ctx context.Context, integration *domain.Integration) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Integration, error)
	FindByTokenLookupHash(ctx context.Context, lookupHash string) (*domain.Integration, error)
	FindByCreatorUserID(ctx context.Context, creatorUserID uuid.UUID) ([]*domain.Integration, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
