package domain

import (
	"time"

	"github.com/m1thrandir225/meridian/pkg/common"
)

type UserRegisteredEvent struct {
	common.BaseDomainEvent
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Timestamp time.Time `json:"registered_at"`
}

type UserAuthenticatedEvent struct {
	common.BaseDomainEvent
	UserID              string    `json:"user_id"`
	Username            string    `json:"username"`
	AuthenticationToken string    `json:"authentication_token"`
	Timestamp           time.Time `json:"timestamp"`
}

type UserProfileUpdatedEvent struct {
	common.BaseDomainEvent
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
	Timestamp     time.Time              `json:"timestamp"`
}

type UserDeletedEvent struct {
	common.BaseDomainEvent
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
