package domain // ChannelActivity represents channel engagement metrics
import (
	"time"

	"github.com/google/uuid"
)

type ChannelActivity struct {
	ID            MetricID
	ChannelID     uuid.UUID
	MessagesCount int64
	MembersCount  int64
	LastMessageAt time.Time
	ActivityScore float64
	Version       int64
}
