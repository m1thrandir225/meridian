package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
	"github.com/m1thrandir225/meridian/pkg/common"
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

	// Lock and update/insert integration with version checking
	lockQuery := `SELECT version FROM integrations WHERE id = $1 FOR UPDATE`
	var currentVersion int64
	err = tx.QueryRow(ctx, lockQuery, integration.ID.String()).Scan(&currentVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if integration.Version != 1 {
				return fmt.Errorf("cannot insert integration %s with version %d: %w", integration.ID.String(), integration.Version, err)
			}

			targetChannels := integration.TargetChannelIDsAsStringSlice()
			insertSQL := `
				INSERT INTO integrations (
					id, service_name, creator_user_id, api_token_hash,
					token_lookup_hash, target_channel_ids, created_at, is_revoked, version
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			_, err := tx.Exec(ctx, insertSQL,
				integration.ID.String(), integration.ServiceName, integration.CreatorUserID.String(),
				integration.HashedAPIToken.Hash(), integration.TokenLookupHash, targetChannels,
				integration.CreatedAt, integration.IsRevoked, integration.Version)
			if err != nil {
				return fmt.Errorf("failed to insert integration: %w", err)
			}
		} else {
			return fmt.Errorf("error locking integration %s: %w", integration.ID.String(), err)
		}
	} else {
		if currentVersion != integration.Version-1 {
			return fmt.Errorf("concurrency conflict saving integration %s: expected version %d, found %d: %w", integration.ID.String(), currentVersion, integration.Version, common.ErrConcurrency)
		}

		targetChannels := integration.TargetChannelIDsAsStringSlice()
		updateSQL := `
			UPDATE integrations SET
				service_name = $2, api_token_hash = $3, target_channel_ids = $4,
				is_revoked = $5, version = $6
			WHERE id = $1 AND version = $7`
		cmdTag, err := tx.Exec(ctx, updateSQL,
			integration.ID.String(), integration.ServiceName, integration.HashedAPIToken.Hash(),
			targetChannels, integration.IsRevoked, integration.Version, currentVersion)
		if err != nil {
			return fmt.Errorf("error updating integration %s: %w", integration.ID.String(), err)
		}

		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("integration %s update affected %d rows, expected 1 (possible lost update): %w", integration.ID.String(), cmdTag.RowsAffected(), common.ErrConcurrency)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresIntegrationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Integration, error) {
	query := `SELECT id, service_name, creator_user_id, api_token_hash,
	                 token_lookup_hash, target_channel_ids, created_at, is_revoked, version
	          FROM integrations WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id.String())
	return r.scanIntegration(row)
}

func (r *PostgresIntegrationRepository) FindByTokenLookupHash(ctx context.Context, lookupHash string) (*domain.Integration, error) {
	query := `SELECT id, service_name, creator_user_id, api_token_hash,
	                 token_lookup_hash, target_channel_ids, created_at, is_revoked, version
	          FROM integrations WHERE token_lookup_hash = $1`
	row := r.db.QueryRow(ctx, query, lookupHash)
	return r.scanIntegration(row)
}

func (r *PostgresIntegrationRepository) FindByCreatorUserID(ctx context.Context, creatorUserID uuid.UUID) ([]*domain.Integration, error) {
	query := `SELECT id, service_name, creator_user_id, api_token_hash,
	                 token_lookup_hash, target_channel_ids, created_at, is_revoked, version
	          FROM integrations WHERE creator_user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, creatorUserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query integrations: %w", err)
	}
	defer rows.Close()

	var integrations []*domain.Integration
	for rows.Next() {
		integration, err := r.scanIntegration(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan integration: %w", err)
		}
		integrations = append(integrations, integration)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return integrations, nil
}

func (r *PostgresIntegrationRepository) scanIntegration(row pgx.Row) (*domain.Integration, error) {
	var integ domain.Integration
	var id uuid.UUID
	var creatorID string
	var tokenHash string
	var targetChannels []string

	err := row.Scan(&id, &integ.ServiceName, &creatorID, &tokenHash, &integ.TokenLookupHash,
		&targetChannels, &integ.CreatedAt, &integ.IsRevoked, &integ.Version)

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

func (r *PostgresIntegrationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM integrations WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete integration %s: %w", id.String(), err)
	}
	if cmdTag.RowsAffected() != 1 {
		return domain.ErrIntegrationNotFound
	}
	return nil
}
