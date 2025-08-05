package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/kafka"

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
	repo                 persistence.UserRepository
	tokenGenerator       AuthTokenGenerator
	AuthTokenValidity    time.Duration
	RefreshTokenValidity time.Duration
	publisher            kafka.EventPublisher
}

func NewUserService(
	repository persistence.UserRepository,
	tokenGenerator AuthTokenGenerator,
	tokenValidity time.Duration,
	refreshTokenValidity time.Duration,
	eventPublisher kafka.EventPublisher,
) *IdentityService {
	return &IdentityService{
		repo:                 repository,
		tokenGenerator:       tokenGenerator,
		AuthTokenValidity:    tokenValidity,
		RefreshTokenValidity: refreshTokenValidity,
		publisher:            eventPublisher,
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

func (s *IdentityService) AuthenticateUser(ctx context.Context, cmd domain.AuthenticateUserCommand) (accessToken string, refreshToken string, claims *auth.TokenClaims, err error) {
	var user *domain.User

	trimmedIdentifier := strings.TrimSpace(cmd.LoginIdentifier)

	if strings.Contains(trimmedIdentifier, "@") {
		email, mailErr := domain.NewEmail(trimmedIdentifier)
		if mailErr != nil {
			return "", "", nil, domain.ErrEmailInvalid
		}
		user, err = s.repo.FindByEmail(ctx, email.String())
	} else {
		username, nameErr := domain.NewUsername(trimmedIdentifier)
		if nameErr != nil {
			return "", "", nil, domain.ErrUsernameInvalid
		}
		user, err = s.repo.FindByUsername(ctx, username.String())
	}

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", "", nil, ErrAuthFailed
		}
		return "", "", nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if err := user.Authenticate(cmd.Password); err != nil {
		return "", "", nil, ErrAuthFailed
	}

	accessToken, claims, err = s.tokenGenerator.GenerateToken(user, s.AuthTokenValidity)
	if err != nil {
		log.Printf("ERROR generating token for user %s: %v", user.ID.String(), err)
		return "", "", nil, ErrTokenGeneration
	}

	refreshToken, err = user.IssueRefreshToken(cmd.Device, cmd.IPAddress, s.RefreshTokenValidity)
	if err != nil {
		return "", "", nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return "", "", nil, fmt.Errorf("failed to save user state after token refresh: %w", err)
	}

	event := domain.CreateUserAuthenticatedEvent(user, accessToken)
	if err := s.publisher.PublishEvent(ctx, event); err != nil {
		log.Printf("ERROR publishing UserAuthenticatedEvent for %s: %v", user.ID.String(), err)
	}

	return accessToken, refreshToken, claims, nil
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

func (s *IdentityService) RefreshAuthentication(ctx context.Context, cmd domain.RefreshTokenCommand) (newAccessToken string, newRefreshToken string, err error) {
	hash := sha256.Sum256([]byte(cmd.RawRefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	user, err := s.repo.FindByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Printf("SECURITY: Attempt to use invalid or revoked refresh token. Hash: %s", tokenHash)
			return "", "", ErrAuthFailed
		}
		return "", "", fmt.Errorf("error finding user by refresh token: %w", err)
	}
	if _, err := user.UseRefreshToken(cmd.RawRefreshToken); err != nil {
		log.Printf("SECURITY: Failed to use valid refresh token for user %s. Possible race or expiry. Error: %v", user.ID.String(), err)
		return "", "", ErrAuthFailed
	}
	newAccessToken, _, err = s.tokenGenerator.GenerateToken(user, s.AuthTokenValidity)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = user.IssueRefreshToken(cmd.Device, cmd.IPAddress, s.RefreshTokenValidity)
	if err != nil {
		return "", "", err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return "", "", fmt.Errorf("failed to save user state after token refresh: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}
