package domain

import "github.com/google/uuid"

type userID struct {
	value uuid.UUID
}

func NewUserID() (*userID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &userID{
		value: id,
	}, nil
}

func (id *userID) String() string {
	return id.value.String()
}

func UserIDFromString(input string) (*userID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &userID{
		value: id,
	}, nil
}
