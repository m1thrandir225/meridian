package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/internal/identity/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

var (
	ErrUserExists      = errors.New("user with given username or email already exists")
	ErrAuthFailed      = domain.ErrAuthentication
	ErrUserNotFound    = persistence.ErrUserNotFound
	ErrTokenGeneration = errors.New("failed to generate authentication token")
)

type UserService struct {
	repo              persistence.UserRepository
	tokenGenerator    AuthTokenGenerator
	authTokenValidity time.Duration
}

func NewUserService(
	repository persistence.UserRepository,
	tokenGenerator AuthTokenGenerator,
	tokenValidity time.Duration,
) *UserService {
	return &UserService{
		repo:              repository,
		tokenGenerator:    tokenGenerator,
		authTokenValidity: tokenValidity,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, cmd domain.RegisterUserCommand) (*domain.User, error) {
	uName, err := domain.NewUsername(cmd.Username)
	if err != nil {
		return nil, err
	}

	eMail, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	existingUser, _ := s.repo.FindByUsername(ctx, uName.String())
	if existingUser != nil {
		return nil, domain.ErrUsernameTaken
	}
	existingUser, _ = s.repo.FindByEmail(ctx, eMail.String())
	if existingUser != nil {
		return nil, domain.ErrEmailTaken
	}

	user, err := domain.NewUser(
		cmd.Username,
		cmd.Email,
		cmd.FirstName,
		cmd.LastName,
		cmd.Password,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	log.Printf("User regriistered: %s (%s)", user.Username, user.ID.String())
	return user, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, cmd domain.AuthenticateUserCommand) (string, *auth.TokenClaims, error) {
	var user *domain.User
	var err error

	trimmedIdentifier := strings.TrimSpace(cmd.LoginIdentifier)

	if strings.Contains(trimmedIdentifier, "@") {
		email, mailErr := domain.NewEmail(trimmedIdentifier)
		if mailErr != nil {
			return "", nil, domain.ErrEmailInvalid
		}
		user, err = s.repo.FindByEmail(ctx, email.String())
	} else {
		username, nameErr := domain.NewUsername(trimmedIdentifier)
		if nameErr != nil {
			return "", nil, domain.ErrUsernameInvalid
		}
		user, err = s.repo.FindByUsername(ctx, username.String())
	}

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", nil, ErrAuthFailed
		}
		return "", nil, fmt.Errorf("error retreiving user: %w", err)
	}

	return "", nil, nil
	// TODO
}
