package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByRefreshTokenHash(ctx context.Context, hash string) (*domain.User, error)
}
