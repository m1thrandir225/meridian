package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"

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
	logger               *logging.Logger
}

func NewUserService(
	repository persistence.UserRepository,
	tokenGenerator AuthTokenGenerator,
	tokenValidity time.Duration,
	refreshTokenValidity time.Duration,
	eventPublisher kafka.EventPublisher,
	logger *logging.Logger,
) *IdentityService {
	return &IdentityService{
		repo:                 repository,
		tokenGenerator:       tokenGenerator,
		AuthTokenValidity:    tokenValidity,
		RefreshTokenValidity: refreshTokenValidity,
		publisher:            eventPublisher,
		logger:               logger,
	}
}

func (s *IdentityService) RegisterUser(ctx context.Context, cmd domain.RegisterUserCommand) (*domain.User, error) {
	logger := s.logger.WithMethod("RegisterUser")
	logger.Info("Registering user")

	uName, err := domain.NewUsername(cmd.Username)
	if err != nil {
		logger.Error("Error creating username", zap.Error(err))
		return nil, err
	}

	eMail, err := domain.NewEmail(cmd.Email)
	if err != nil {
		logger.Error("Error creating email", zap.Error(err))
		return nil, err
	}

	existingUser, _ := s.repo.FindByUsername(ctx, uName.String())
	if existingUser != nil {
		logger.Error("Username already taken", zap.String("username", uName.String()))
		return nil, domain.ErrUsernameTaken
	}
	existingUser, _ = s.repo.FindByEmail(ctx, eMail.String())
	if existingUser != nil {
		logger.Error("Email already taken", zap.String("email", eMail.String()))
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
		logger.Error("Error saving user", zap.Error(err))
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		logger.Error("Error publishing UserRegisteredEvent", zap.String("user_id", user.ID.String()), zap.Error(err))
	}

	user.ClearEvents()

	logger.Info("User registered", zap.String("username", user.Username.String()), zap.String("user_id", user.ID.String()))
	return user, nil
}

func (s *IdentityService) AuthenticateUser(ctx context.Context, cmd domain.AuthenticateUserCommand) (accessToken string, refreshToken string, claims *auth.TokenClaims, err error) {
	logger := s.logger.WithMethod("AuthenticateUser")
	logger.Info("Authenticating user")

	var user *domain.User

	trimmedIdentifier := strings.TrimSpace(cmd.LoginIdentifier)

	if strings.Contains(trimmedIdentifier, "@") {
		email, mailErr := domain.NewEmail(trimmedIdentifier)
		if mailErr != nil {
			logger.Error("Error creating email", zap.Error(mailErr))
			return "", "", nil, domain.ErrEmailInvalid
		}
		user, err = s.repo.FindByEmail(ctx, email.String())
		if err != nil {
			logger.Error("Error retrieving user", zap.Error(err))
			return "", "", nil, fmt.Errorf("error retrieving user: %w", err)
		}
	} else {
		username, nameErr := domain.NewUsername(trimmedIdentifier)
		if nameErr != nil {
			logger.Error("Error creating username", zap.Error(nameErr))
			return "", "", nil, domain.ErrUsernameInvalid
		}
		user, err = s.repo.FindByUsername(ctx, username.String())
		if err != nil {
			logger.Error("Error retrieving user", zap.Error(err))
			return "", "", nil, fmt.Errorf("error retrieving user: %w", err)
		}
	}

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", zap.String("identifier", trimmedIdentifier))
			return "", "", nil, ErrAuthFailed
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return "", "", nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if user == nil {
		logger.Error("User is nil but no error returned")
		return "", "", nil, fmt.Errorf("user is nil but no error returned")
	}

	if err := user.Authenticate(cmd.Password); err != nil {
		logger.Error("Authentication failed", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", nil, ErrAuthFailed
	}

	accessToken, claims, err = s.tokenGenerator.GenerateToken(user, s.AuthTokenValidity)
	if err != nil {
		logger.Error("Error generating token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", nil, ErrTokenGeneration
	}

	refreshToken, err = user.IssueRefreshToken(cmd.Device, cmd.IPAddress, s.RefreshTokenValidity)
	if err != nil {
		logger.Error("Error issuing refresh token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		logger.Error("Error saving user", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", nil, fmt.Errorf("failed to save user state after token refresh: %w", err)
	}

	event := domain.CreateUserAuthenticatedEvent(user, accessToken)
	if err := s.publisher.PublishEvent(ctx, event); err != nil {
		logger.Error("Error publishing UserAuthenticatedEvent", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	logger.Info("User authenticated", zap.String("user_id", user.ID.String()))

	return accessToken, refreshToken, claims, nil
}

func (s *IdentityService) GetUser(ctx context.Context, cmd domain.GetUserCommand) (*domain.User, error) {
	logger := s.logger.WithMethod("GetUser")
	logger.Info("Getting user")

	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", zap.String("user_id", userId.String()))
			return nil, ErrUserNotFound
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	logger.Info("User retrieved", zap.String("user_id", user.ID.String()))
	return user, nil
}

func (s *IdentityService) GetUsers(ctx context.Context, cmd domain.GetUsersCommand) ([]*domain.User, error) {
	logger := s.logger.WithMethod("GetUsers")
	logger.Info("Getting users")
	if len(cmd.UserIds) == 0 {
		return []*domain.User{}, nil
	}
	userIds := make([]uuid.UUID, len(cmd.UserIds))
	for i, id := range cmd.UserIds {
		userId, err := uuid.Parse(id)
		if err != nil {
			logger.Error("Invalid user ID", zap.Error(err))
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
		userIds[i] = userId
	}
	users, err := s.repo.FindByIds(ctx, userIds)
	if err != nil {
		logger.Error("Error retrieving users", zap.Error(err))
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	logger.Info("Users retrieved", zap.Int("count", len(users)))
	return users, nil
}

func (s *IdentityService) UpdateUserProfile(ctx context.Context, cmd domain.UpdateUserProfileCommand) (*domain.User, error) {
	logger := s.logger.WithMethod("UpdateUserProfile")
	logger.Info("Updating user profile")

	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", zap.String("user_id", userId.String()))
			return nil, ErrUserNotFound
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	err = user.UpdateProfile(cmd.NewUsername, cmd.NewEmail, cmd.NewFirstName, cmd.NewLastName)
	if err != nil {
		logger.Error("Error updating user profile", zap.String("user_id", user.ID.String()), zap.Error(err))
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		logger.Error("Error saving user", zap.String("user_id", user.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		logger.Error("Error publishing UserUpdatedEvent", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	logger.Info("User profile updated", zap.String("user_id", user.ID.String()))
	user.ClearEvents()

	return user, nil
}

func (s *IdentityService) UpdateUserPassword(ctx context.Context, cmd domain.UpdateUserPasswordCommand) error {
	logger := s.logger.WithMethod("UpdateUserPassword")
	logger.Info("Updating user password")

	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", zap.String("user_id", userId.String()))
			return ErrUserNotFound
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return fmt.Errorf("error retrieving user: %w", err)
	}

	err = user.UpdatePassword(cmd.NewPassword)
	if err != nil {
		logger.Error("Error updating user password", zap.String("user_id", user.ID.String()), zap.Error(err))
		return err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		logger.Error("Error saving user", zap.String("user_id", user.ID.String()), zap.Error(err))
		return fmt.Errorf("failed to save user: %w", err)
	}

	if err := s.publisher.PublishEvents(ctx, user.Events()); err != nil {
		logger.Error("Error publishing UserUpdatedEvent", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	logger.Info("User password updated", zap.String("user_id", user.ID.String()))
	user.ClearEvents()

	return nil
}

func (s *IdentityService) DeleteUser(ctx context.Context, cmd domain.DeleteUserCommand) error {
	logger := s.logger.WithMethod("DeleteUser")
	logger.Info("Deleting user")

	userId, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("User not found", zap.String("user_id", userId.String()))
			return ErrUserNotFound
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return err
	}

	if err := s.repo.Delete(ctx, userId); err != nil {
		logger.Error("Error deleting user", zap.String("user_id", user.ID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	event := domain.CreateUserDeletedEvent(user)
	if err := s.publisher.PublishEvent(ctx, event); err != nil {
		logger.Error("Error publishing UserDeletedEvent", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	logger.Info("User deleted", zap.String("user_id", user.ID.String()))
	return nil
}

func (s *IdentityService) RefreshAuthentication(ctx context.Context, cmd domain.RefreshTokenCommand) (newAccessToken string, newRefreshToken string, err error) {
	logger := s.logger.WithMethod("RefreshAuthentication")
	logger.Info("Refreshing authentication")

	hash := sha256.Sum256([]byte(cmd.RawRefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	user, err := s.repo.FindByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Error("Attempt to use invalid or revoked refresh token", zap.String("token_hash", tokenHash))
			return "", "", ErrAuthFailed
		}
		return "", "", fmt.Errorf("error finding user by refresh token: %w", err)
	}
	if _, err := user.UseRefreshToken(cmd.RawRefreshToken); err != nil {
		logger.Error("Failed to use valid refresh token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", ErrAuthFailed
	}
	newAccessToken, _, err = s.tokenGenerator.GenerateToken(user, s.AuthTokenValidity)
	if err != nil {
		logger.Error("Error generating token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", err
	}

	newRefreshToken, err = user.IssueRefreshToken(cmd.Device, cmd.IPAddress, s.RefreshTokenValidity)
	if err != nil {
		logger.Error("Error issuing refresh token", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		logger.Error("Error saving user", zap.String("user_id", user.ID.String()), zap.Error(err))
		return "", "", fmt.Errorf("failed to save user state after token refresh: %w", err)
	}

	logger.Info("Authentication refreshed", zap.String("user_id", user.ID.String()))
	return newAccessToken, newRefreshToken, nil
}
