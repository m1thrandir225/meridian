package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateProfile(ctx context.Context, fields map[string]interface{}) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
