package domain

import "errors"

var (
	ErrUserIDInvalid   = errors.New("user ID is invalid")
	ErrUsernameInvalid = errors.New("Username is invalid (length, characters)")
	ErrUsernameTaken   = errors.New("Username is already taken")
	ErrEmailInvalid    = errors.New("UserEmail format is invalid")
	ErrEmailTaken      = errors.New("UserEmail is already taken")
	ErrPasswordPolicy  = errors.New("password does not meet policy requirements")
	ErrAuthentication  = errors.New("authentication failed: invalid credentials")
	ErrUserNotFound    = errors.New("the specified user was not found")
	ErrAuthFailed      = errors.New("authentication failed")
)
