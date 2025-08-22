package domain

import (
	"time"

	"github.com/m1thrandir225/meridian/pkg/common"
)

type UserRegisteredEvent struct {
	common.BaseDomainEvent
	UserID    string    `json:"user_id"`
	Username  string    `json:"Username"`
	Email     string    `json:"UserEmail"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Timestamp time.Time `json:"registered_at"`
}

type UserAuthenticatedEvent struct {
	common.BaseDomainEvent
	UserID              string    `json:"user_id"`
	Username            string    `json:"Username"`
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

func CreateUserRegisteredEvent(user *User) UserRegisteredEvent {
	base := common.NewBaseDomainEvent("UserRegistered", user.ID.value, user.Version, "User")

	return UserRegisteredEvent{
		BaseDomainEvent: base,
		UserID:          user.ID.value.String(),
		Username:        user.Username.String(),
		Email:           user.Email.String(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Timestamp:       time.Now(),
	}
}

func CreateUserAuthenticatedEvent(user *User, token string) UserAuthenticatedEvent {
	base := common.NewBaseDomainEvent("UserAuthenticated", user.ID.value, user.Version, "User")

	return UserAuthenticatedEvent{
		BaseDomainEvent:     base,
		UserID:              user.ID.value.String(),
		Username:            user.Username.String(),
		AuthenticationToken: token,
		Timestamp:           time.Now(),
	}
}

func CreateUserProfileUpdated(user *User, updatedFields map[string]interface{}) UserProfileUpdatedEvent {
	base := common.NewBaseDomainEvent("UserProfileUpdated", user.ID.value, user.Version, "User")

	return UserProfileUpdatedEvent{
		BaseDomainEvent: base,
		UserID:          user.ID.value.String(),
		UpdatedFields:   updatedFields,
		Timestamp:       time.Now(),
	}
}

func CreateUserDeletedEvent(user *User) UserDeletedEvent {
	base := common.NewBaseDomainEvent("UserDeleted", user.ID.value, user.Version, "User")

	return UserDeletedEvent{
		BaseDomainEvent: base,
		UserID:          user.ID.value.String(),
		Timestamp:       time.Now(),
	}
}
