package persistence

import (
	"context"

	"github.com/google/uuid"
	models "github.com/m1thrandir225/meridian/internal/messaging/domain"
)

type ChannelRepository interface {
	Save(ctx context.Context, channel *models.Channel) error
	FindById(ctx context.Context, id uuid.UUID) (*models.Channel, error)
	FindMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]models.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
