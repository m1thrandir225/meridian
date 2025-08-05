package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/common"
)

var _ UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: pool,
	}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	lockQuery := `SELECT version FROM users WHERE id = $1 FOR UPDATE`
	var currentVersion int64

	err = tx.QueryRow(ctx, lockQuery, user.ID.String()).Scan(&currentVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if user.Version != 1 {
				return fmt.Errorf("cannot insert user %s with version %d: %w", &user.ID, user.Version, err)
			}

			insertQuery := `
			INSERT INTO users(id, username, first_name, last_name, email, password, version, registartion_time)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			`
			_, err := tx.Exec(
				ctx,
				insertQuery,
				user.ID.String(),
				user.Username.String(),
				user.FirstName,
				user.LastName,
				user.Email.String(),
				user.PasswordHash.String(),
				user.Version,
				user.RegistrationTime,
			)
			if err != nil {
				return fmt.Errorf("error inserting user %s: %w", &user.ID, err)
			}
		} else {
			return fmt.Errorf("error locking user %s: %w", &user.ID, err)
		}
	} else {
		expectedVersion := user.Version - 1
		if currentVersion != expectedVersion {
			return fmt.Errorf("concurrency conflict saving user %s: expected version %d, found version %d: %w", &user.ID, currentVersion, user.Version, err)
		}

		updateQuery := `
			UPDATE users SET username=$1, first_name = $2, last_name = $3, email = $4, password = $5, version = $6
		WHERE id = $7 AND version = $8
		`

		cmdTag, err := tx.Exec(
			ctx,
			updateQuery,
			user.Username.String(),
			user.FirstName,
			user.LastName,
			user.Email.String(),
			user.PasswordHash.String(),
			user.Version,
			user.ID.String(),
			expectedVersion,
		)
		if err != nil {
			return fmt.Errorf("error updating user %s: %w", &user.ID, err)
		}

		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("user %s update affected %d rows, expected 1( possible lost update): %w", &user.ID, cmdTag.RowsAffected(), common.ErrConcurrency)
		}
	}
	deleteTokensSQL := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err = tx.Exec(ctx, deleteTokensSQL, user.ID.String())
	if err != nil {
		return fmt.Errorf("failed to clear existing refresh tokens: %w", err)
	}
	for _, rt := range user.RefreshTokens {
		insertTokenSQL := `
			INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked, created_at, device, ip_address)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := tx.Exec(
			ctx,
			insertTokenSQL,
			rt.ID.String(),
			rt.UserID.String(),
			rt.TokenHash,
			rt.ExpiresAt,
			rt.IsRevoked,
			rt.CreatedAt,
			rt.Device,
			rt.IPAddress,
		)
		if err != nil {
			return fmt.Errorf("failed to insert refresh token %s: %w", rt.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return r.findByField(ctx, "id", id.String())
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.findByField(ctx, "username", username)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findByField(ctx, "email", email)
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	deleteQuery := `DELETE FROM users WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, deleteQuery, id.String())
	if err != nil {
		return fmt.Errorf("error deleting user %s: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s was not found for deletion: %w", id, common.ErrNotFound)
	}

	return nil
}

func (r *PostgresUserRepository) FindByRefreshTokenHash(ctx context.Context, hash string) (*domain.User, error) {
	query := `
		SELECT u.id FROM users u
		JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token_hash = $1 AND rt.is_revoked = FALSE AND rt.expires_at > NOW()`
	var userId uuid.UUID
	err := r.db.QueryRow(ctx, query, hash).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by refresh token hash: %w", err)
	}

	if userId == uuid.Nil {
		return nil, domain.ErrUserIDInvalid
	}
	return r.FindById(ctx, userId)
}

func (r *PostgresUserRepository) findByField(ctx context.Context, fieldName string, value any) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT id, username, first_name, last_name, email, password, version, registartion_time
		FROM users
		WHERE %s = $1`, fieldName)

	row := r.db.QueryRow(ctx, query, value)
	return r.scanUser(row)
}

func (r *PostgresUserRepository) scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	var userId uuid.UUID
	var usernameStr, emailStr, passwordStr string
	var regTime time.Time

	err := row.Scan(
		&userId,
		&usernameStr,
		&user.FirstName,
		&user.LastName,
		&emailStr,
		&passwordStr,
		&user.Version,
		&regTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	domId, _ := domain.UserIDFromString(userId.String())
	domUsername, errUN := domain.NewUsername(usernameStr)
	if errUN != nil {
		return nil, errUN
	}
	domEmail, errEM := domain.NewEmail(emailStr)
	if errEM != nil {
		return nil, errEM
	}
	domPassHash, errPH := domain.FromHashedString(passwordStr)
	if errPH != nil {
		return nil, errPH
	}

	user.Email = domEmail
	user.ID = *domId
	user.Username = domUsername
	user.PasswordHash = domPassHash
	user.RegistrationTime = regTime

	tokensQuery := `SELECT id, user_id, token_hash, expires_at, is_revoked, created_at, device, ip_address
	                FROM refresh_tokens WHERE user_id = $1 AND is_revoked = FALSE`
	rows, err := r.db.Query(context.Background(), tokensQuery, user.ID.String())

	defer rows.Close()
	var tokens []*domain.RefreshToken
	for rows.Next() {
		var rt domain.RefreshToken
		var userID uuid.UUID
		err := rows.Scan(&rt.ID, &userID, &rt.TokenHash, &rt.ExpiresAt, &rt.IsRevoked, &rt.CreatedAt, &rt.Device, &rt.IPAddress)
		if err != nil {
			log.Printf("ERROR: Failed to scan refresh token row for user %s: %v", user.ID.String(), err)
			continue
		}
		userIDVO, _ := domain.UserIDFromString(userID.String())

		rt.UserID = *userIDVO
		tokens = append(tokens, &rt)
	}
	user.RefreshTokens = tokens
	return &user, nil
}
