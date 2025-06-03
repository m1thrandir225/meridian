package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
)

var _ UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) UpdateProfile(ctx context.Context, fields map[string]interface{}) (*domain.User, error) {
	return nil, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}
