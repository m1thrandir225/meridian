package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/kafka"
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
	ErrUserNotFound    = domain.ErrUserNotFound
	ErrTokenGeneration = errors.New("failed to generate authentication token")
)

type IdentityService struct {
	repo              persistence.UserRepository
	tokenGenerator    AuthTokenGenerator
	authTokenValidity time.Duration
	publisher         kafka.EventPublisher
}

func NewUserService(
	repository persistence.UserRepository,
	tokenGenerator AuthTokenGenerator,
	tokenValidity time.Duration,
	eventPublisher kafka.EventPublisher,
) *IdentityService {
	return &IdentityService{
		repo:              repository,
		tokenGenerator:    tokenGenerator,
		authTokenValidity: tokenValidity,
		publisher:         eventPublisher,
	}
}

func (s *IdentityService) RegisterUser(ctx context.Context, cmd domain.RegisterUserCommand) (*domain.User, error) {
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

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		log.Printf("ERROR publishing UserRegisteredEvent for %s: %v", user.ID.String(), err)
	}

	user.ClearEvents()

	log.Printf("User registered: %s (%s)", user.Username, user.ID.String())
	return user, nil
}

func (s *IdentityService) AuthenticateUser(ctx context.Context, cmd domain.AuthenticateUserCommand) (string, *auth.TokenClaims, error) {
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
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", nil, ErrAuthFailed
		}
		return "", nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if err := user.Authenticate(cmd.Password); err != nil {
		return "", nil, ErrAuthFailed
	}

	tokenString, claims, err := s.tokenGenerator.GenerateToken(user, s.authTokenValidity)
	if err != nil {
		log.Printf("ERROR generating token for user %s: %v", user.ID.String(), err)
		return "", nil, ErrTokenGeneration
	}

	event := domain.CreateUserAuthenticatedEvent(user, tokenString)
	if err := s.publisher.PublishEvent(ctx, event); err != nil {
		log.Printf("ERROR publishing UserAuthenticatedEvent for %s: %v", user.ID.String(), err)
	}

	return tokenString, claims, nil
}

func (s *IdentityService) GetUser(ctx context.Context, cmd domain.GetUserCommand) (*domain.User, error) {
	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return user, nil
}

func (s *IdentityService) UpdateUserProfile(ctx context.Context, cmd domain.UpdateUserProfileCommand) (*domain.User, error) {
	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	err = user.UpdateProfile(cmd.NewUsername, cmd.NewEmail, cmd.NewFirstName, cmd.NewLastName)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		log.Printf("ERROR publishing UserUpdatedEvent for %s: %v", user.ID.String(), err)
	}
	user.ClearEvents()

	return user, nil
}

func (s *IdentityService) UpdateUserPassword(ctx context.Context, cmd domain.UpdateUserPasswordCommand) error {
	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ErrUserNotFound
		}
	}

	err = user.UpdatePassword(cmd.NewPassword)
	if err != nil {
		return err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		log.Printf("ERROR publishing UserUpdatedEvent for %s: %v", user.ID.String(), err)
	}
	user.ClearEvents()

	return nil
}

func (s *IdentityService) DeleteUser(ctx context.Context, cmd domain.DeleteUserCommand) error {
	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, userId); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	event := domain.CreateUserDeletedEvent(user)
	if err := s.publisher.PublishEvent(ctx, event); err != nil {
		log.Printf("ERROR publishing UserDeletedEvent for %s: %v", user.ID.String(), err)
	}
	return nil
}
