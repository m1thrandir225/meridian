package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
)

type PostgresIntegrationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresIntegrationRepository(pool *pgxpool.Pool) *PostgresIntegrationRepository {
	return &PostgresIntegrationRepository{
		db: pool,
	}
}

func (r *PostgresIntegrationRepository) Save(ctx context.Context, integration *domain.Integration) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	targetChannels := integration.TargetChannelIDsAsStringSlice()

	updateSQL := `
		UPDATE integrations SET
			service_name = $2, api_token_hash = $3, target_channel_ids = $4, is_revoked = $5
		WHERE id = $1`
	tag, err := tx.Exec(ctx, updateSQL,
		integration.ID.String(), integration.ServiceName, integration.HashedAPIToken.Hash(),
		targetChannels, integration.IsRevoked)

	if err == nil && tag.RowsAffected() == 1 {
		return tx.Commit(ctx)
	}

	insertSQL := `
		INSERT INTO integrations (
			id, service_name, creator_user_id, api_token_hash,
			target_channel_ids, created_at, is_revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(ctx, insertSQL,
		integration.ID.String(), integration.ServiceName, integration.CreatorUserID.String(),
		integration.HashedAPIToken.Hash(), targetChannels,
		integration.CreatedAt, integration.IsRevoked)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("failed to save due to unique constraint: %w", err)
		}
		return fmt.Errorf("failed to insert integration: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresIntegrationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Integration, error) {
	query := `SELECT id, service_name, creator_user_id, api_token_hash,
	                 token_lookup_hash, target_channel_ids, created_at, is_revoked
	          FROM integrations WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id.String())
	return r.scanIntegration(row)
}

func (r *PostgresIntegrationRepository) FindByTokenLookupHash(ctx context.Context, lookupHash string) (*domain.Integration, error) {
	query := `SELECT id, service_name, creator_user_id, api_token_hash,
	                 token_lookup_hash, target_channel_ids, created_at, is_revoked
	          FROM integrations WHERE token_lookup_hash = $1`
	row := r.db.QueryRow(ctx, query, lookupHash)
	return r.scanIntegration(row)
}

func (r *PostgresIntegrationRepository) scanIntegration(row pgx.Row) (*domain.Integration, error) {
	var integ domain.Integration
	var id uuid.UUID
	var creatorID string
	var tokenHash string
	var targetChannels []string

	err := row.Scan(&id, &integ.ServiceName, &creatorID, &tokenHash, &integ.TokenLookupHash,
		&targetChannels, &integ.CreatedAt, &integ.IsRevoked)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrIntegrationNotFound
		}
		return nil, fmt.Errorf("database scan error: %w", err)
	}

	integrationId, err := domain.NewIntegrationIDFromString(id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create integration ID: %w", err)
	}
	integ.CreatorUserID = domain.UserIDRef(creatorID)
	apiToken, err := domain.NewAPIToken(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create API token: %w", err)
	}
	integ.TargetChannelIDs = make([]domain.ChannelIDRef, len(targetChannels))
	for i, ch := range targetChannels {
		integ.TargetChannelIDs[i] = domain.ChannelIDRef(ch)
	}

	integ.ID = *integrationId
	integ.HashedAPIToken = *apiToken

	return &integ, nil
}
