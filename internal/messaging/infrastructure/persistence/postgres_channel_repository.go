package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/m1thrandir225/meridian/internal/messaging/domain"
)

var _ ChannelRepository = (*PostgresChannelRepository)(nil)

type PostgresChannelRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresChannelRepository(pool *pgxpool.Pool) *PostgresChannelRepository {
	return &PostgresChannelRepository{
		pool: pool,
	}
}

// FIXME: implement
func (r *PostgresChannelRepository) Save(ctx context.Context, channel *models.Channel) error {
	return nil
}

// FIXME: implement
func (r *PostgresChannelRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	return nil, nil
}

// FIXME: implement
func (r *PostgresChannelRepository) FindMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]models.Message, error) {
	return nil, nil
}

// FIXME: implement
func (r *PostgresChannelRepository) Delete(ctx context.Context, id uuid.UUID) {}
