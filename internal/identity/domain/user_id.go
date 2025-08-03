package domain

import "github.com/google/uuid"

type UserID struct {
	value uuid.UUID
}

func NewUserID() (*UserID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &UserID{
		value: id,
	}, nil
}

func (id *UserID) String() string {
	return id.value.String()
}

func UserIDFromString(input string) (*UserID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &UserID{
		value: id,
	}, nil
}
