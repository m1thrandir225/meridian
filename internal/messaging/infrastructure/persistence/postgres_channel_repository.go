package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/pkg/common"
)

var _ ChannelRepository = (*PostgresChannelRepository)(nil)

type PostgresChannelRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresChannelRepository(pool *pgxpool.Pool) *PostgresChannelRepository {
	return &PostgresChannelRepository{
		pool: pool,
	}
}

// Helper method to scan basic channel data
func (r *PostgresChannelRepository) scanChannelBasic(row pgx.Row) (*models.Channel, error) {
	var channel models.Channel
	var topic *string
	var lastMsgTime *time.Time

	err := row.Scan(
		&channel.ID,
		&channel.Name,
		&topic,
		&channel.CreatorUserID,
		&channel.CreationTime,
		&lastMsgTime,
		&channel.IsArchived,
		&channel.Version,
	)
	if err != nil {
		return nil, err
	}

	if topic != nil {
		channel.Topic = *topic
	}
	if lastMsgTime != nil {
		channel.LastMessageTime = *lastMsgTime
	}

	return &channel, nil
}

// Helper method to load members for a channel
func (r *PostgresChannelRepository) loadMembers(ctx context.Context, channelID uuid.UUID) ([]models.Member, error) {
	query := `
		SELECT user_id, role, joined_at, last_read
		FROM members
		WHERE channel_id = $1
		ORDER BY joined_at ASC
	`

	rows, err := r.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("error querying members for channel %s: %w", channelID, err)
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var memberUserID uuid.UUID
		var memberRole string
		var memberJoinedAt time.Time
		var memberLastRead sql.NullTime

		if err := rows.Scan(&memberUserID, &memberRole, &memberJoinedAt, &memberLastRead); err != nil {
			return nil, fmt.Errorf("error scanning member for channel %s: %w", channelID, err)
		}

		var actualLastRead time.Time
		if memberLastRead.Valid {
			actualLastRead = memberLastRead.Time
		}

		member := models.RehydrateMember(memberUserID, memberRole, memberJoinedAt, actualLastRead)
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members for channel %s: %w", channelID, err)
	}

	return members, nil
}

// Helper method to load invites for a channel
func (r *PostgresChannelRepository) loadInvites(ctx context.Context, channelID uuid.UUID) ([]models.ChannelInvite, error) {
	query := `
		SELECT id, channel_id, created_by_user_id, invite_code, expires_at, max_uses, current_uses, created_at, is_active
		FROM channel_invites
		WHERE channel_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("error querying invites for channel %s: %w", channelID, err)
	}
	defer rows.Close()

	var invites []models.ChannelInvite
	for rows.Next() {
		var invite models.ChannelInvite
		var maxUses *int

		err := rows.Scan(
			&invite.ID,
			&invite.ChannelID,
			&invite.CreatedByUserID,
			&invite.InviteCode,
			&invite.ExpiresAt,
			&maxUses,
			&invite.CurrentUses,
			&invite.CreatedAt,
			&invite.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning invite for channel %s: %w", channelID, err)
		}

		invite.MaxUse = maxUses
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invites for channel %s: %w", channelID, err)
	}

	return invites, nil
}

// Helper method to load messages for a channel
func (r *PostgresChannelRepository) loadMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]models.Message, error) {
	query := `
		SELECT id, channel_id, sender_user_id, integration_id,
		       content_text, content_mentions, content_link, content_formatted,
		       created_at, parent_message_id
		FROM messages
		WHERE channel_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.pool.Query(ctx, query, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying messages for channel %s: %w", channelID, err)
	}
	defer rows.Close()

	var messages []models.Message
	var messageIDs []uuid.UUID

	for rows.Next() {
		var messageId uuid.UUID
		var senderUserID, integrationID, parentMessageID *uuid.UUID
		var mentions []uuid.UUID
		var links []string
		var text string
		var timestamp time.Time
		var isFormatted bool

		err := rows.Scan(
			&messageId,
			&channelID,
			&senderUserID,
			&integrationID,
			&text,
			&mentions,
			&links,
			&isFormatted,
			&timestamp,
			&parentMessageID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message for channel %s: %w", channelID, err)
		}

		content := models.RehydrateMessageContent(text, mentions, links, isFormatted)

		msg := models.RehydrateMessage(
			messageId,
			channelID,
			senderUserID,
			integrationID,
			parentMessageID,
			content,
			[]models.Reaction{},
			timestamp,
		)

		messages = append(messages, msg)
		messageIDs = append(messageIDs, messageId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages for channel %s: %w", channelID, err)
	}

	// Load reactions for all messages if we have any
	if len(messages) > 0 {
		if err := r.loadReactionsForMessages(ctx, messages, messageIDs); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

// Helper method to load reactions for messages
func (r *PostgresChannelRepository) loadReactionsForMessages(ctx context.Context, messages []models.Message, messageIDs []uuid.UUID) error {
	query := `
		SELECT id, message_id, user_id, reaction_type, created_at
		FROM reactions
		WHERE message_id = ANY($1)
		ORDER BY message_id, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, messageIDs)
	if err != nil {
		return fmt.Errorf("error querying reactions for messages: %w", err)
	}
	defer rows.Close()

	reactionsByMessageID := make(map[uuid.UUID][]models.Reaction)
	for rows.Next() {
		var reactionID uuid.UUID
		var reactionMessageID uuid.UUID
		var reactionUserID uuid.UUID
		var reactionType string
		var reactionTimestamp time.Time

		err := rows.Scan(
			&reactionID,
			&reactionMessageID,
			&reactionUserID,
			&reactionType,
			&reactionTimestamp,
		)
		if err != nil {
			return fmt.Errorf("error scanning reaction: %w", err)
		}

		reaction := models.RehydrateReaction(reactionID, reactionMessageID, reactionUserID, reactionType, reactionTimestamp)
		reactionsByMessageID[reactionMessageID] = append(reactionsByMessageID[reactionMessageID], reaction)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating reactions: %w", err)
	}

	// Attach reactions to messages
	for i := range messages {
		if reactions, ok := reactionsByMessageID[messages[i].GetId()]; ok {
			messages[i].SetLoadedReactions(reactions)
		}
	}

	return nil
}

// Helper method to save members using COPY
func (r *PostgresChannelRepository) saveMembers(ctx context.Context, tx pgx.Tx, channelID uuid.UUID, members []models.Member) error {
	if len(members) == 0 {
		return nil
	}

	// Delete existing members
	deleteQuery := `DELETE FROM members WHERE channel_id = $1`
	_, err := tx.Exec(ctx, deleteQuery, channelID)
	if err != nil {
		return fmt.Errorf("error deleting old members for channel %s: %w", channelID, err)
	}

	// Insert new members using COPY
	memberRows := make([][]any, len(members))
	for i, member := range members {
		memberRows[i] = []any{
			channelID,
			member.GetId(),
			member.GetRole(),
			member.GetJoinedAt(),
			member.GetLastRead(),
		}
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"members"},
		[]string{"channel_id", "user_id", "role", "joined_at", "last_read"},
		pgx.CopyFromRows(memberRows),
	)
	if err != nil {
		return fmt.Errorf("error inserting members for channel %s using copy: %w", channelID, err)
	}

	if int(copyCount) != len(members) {
		return fmt.Errorf("expected %d members to be inserted for channel %s, but got %d", len(members), channelID, copyCount)
	}

	return nil
}

// Helper method to save invites
func (r *PostgresChannelRepository) saveInvites(ctx context.Context, tx pgx.Tx, channelID uuid.UUID, invites []models.ChannelInvite) error {
	if len(invites) == 0 {
		return nil
	}

	// Delete existing invites
	deleteQuery := `DELETE FROM channel_invites WHERE channel_id = $1`
	_, err := tx.Exec(ctx, deleteQuery, channelID)
	if err != nil {
		return fmt.Errorf("error deleting old invites for channel %s: %w", channelID, err)
	}

	// Insert new invites
	insertQuery := `
		INSERT INTO channel_invites (id, channel_id, created_by_user_id, invite_code, expires_at, max_uses, current_uses, created_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, invite := range invites {
		var maxUses *int
		if invite.GetMaxUse() != nil {
			maxUses = invite.GetMaxUse()
		}

		_, err := tx.Exec(ctx, insertQuery,
			invite.GetID(),
			invite.GetChannelID(),
			invite.GetCreatedByUserID(),
			invite.GetInviteCode(),
			invite.GetExpiresAt(),
			maxUses,
			invite.GetCurrentUses(),
			invite.GetCreatedAt(),
			invite.GetIsActive(),
		)
		if err != nil {
			return fmt.Errorf("error inserting invite %s for channel %s: %w", invite.GetID(), channelID, err)
		}
	}

	return nil
}

func (r *PostgresChannelRepository) Save(ctx context.Context, channel *models.Channel) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock and update/insert channel
	lockQuery := `SELECT version FROM channels WHERE id = $1 FOR UPDATE`
	var currentVersion int64
	err = tx.QueryRow(ctx, lockQuery, channel.ID).Scan(&currentVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if channel.Version != 1 {
				return fmt.Errorf("cannot insert channel %s with version %d: %w", channel.ID, channel.Version, err)
			}
			insertQuery := `
				INSERT INTO channels(id, name, topic, creator_user_id, creation_time, last_message_time, is_archived, version)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			`
			_, err := tx.Exec(ctx, insertQuery, channel.ID, channel.Name, channel.Topic, channel.CreatorUserID, channel.CreationTime, channel.LastMessageTime, channel.IsArchived, channel.Version)
			if err != nil {
				return fmt.Errorf("error inserting channel %s: %w", channel.ID, err)
			}
		} else {
			return fmt.Errorf("error locking channel %s: %w", channel.ID, err)
		}
	} else {
		if currentVersion != channel.Version-1 {
			return fmt.Errorf("concurrency conflict saving channel %s: expected version %d, found %d: %w", channel.ID, currentVersion, channel.Version, common.ErrConcurrency)
		}

		updateQuery := `
			UPDATE channels SET name = $1, topic = $2, last_message_time = $3, is_archived = $4, version = $5
			WHERE id = $6 and version = $7
		`
		cmdTag, err := tx.Exec(ctx, updateQuery, channel.Name, channel.Topic, channel.LastMessageTime, channel.IsArchived, channel.Version, channel.ID, currentVersion)
		if err != nil {
			return fmt.Errorf("error updating channel %s: %w", channel.ID, err)
		}

		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("channel %s update affected %d rows, expected 1 (possible lost update): %w", channel.ID, cmdTag.RowsAffected(), common.ErrConcurrency)
		}
	}

	if err := r.saveMembers(ctx, tx, channel.ID, channel.Members); err != nil {
		return err
	}

	if err := r.saveInvites(ctx, tx, channel.ID, channel.Invites); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction for channel %s: %w", channel.ID, err)
	}
	return nil
}

func (r *PostgresChannelRepository) FindUserChannels(ctx context.Context, userID uuid.UUID) ([]*models.Channel, error) {
	query := `
		SELECT DISTINCT c.id, c.name, c.topic, c.creator_user_id, c.creation_time, c.last_message_time, c.is_archived, c.version
		FROM channels c
		LEFT JOIN members m ON c.id = m.channel_id
		WHERE c.creator_user_id = $1 OR m.user_id = $1
		ORDER BY c.last_message_time DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying for channels for user %s: %w", userID, err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel, err := r.scanChannelBasic(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning channel for user %s: %w", userID, err)
		}

		// Load members for each channel
		members, err := r.loadMembers(ctx, channel.ID)
		if err != nil {
			return nil, err
		}

		channel.Members = members
		channel.Messages = []models.Message{}
		channel.Invites = []models.ChannelInvite{}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating channels for user %s: %w", userID, err)
	}

	return channels, nil
}

func (r *PostgresChannelRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	query := `
		SELECT id, name, topic, creator_user_id, creation_time, last_message_time, is_archived, version
		FROM channels
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	channel, err := r.scanChannelBasic(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel with ID %s not found: %w", id, common.ErrNotFound)
		}
		return nil, fmt.Errorf("error scanning channel %s: %w", id, err)
	}

	// Load members and invites using helper methods
	members, err := r.loadMembers(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	invites, err := r.loadInvites(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	channel.Members = members
	channel.Invites = invites
	channel.Messages = []models.Message{}

	return channel, nil
}

func (r *PostgresChannelRepository) FindMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]models.Message, error) {
	return r.loadMessages(ctx, channelID, limit, offset)
}

func (r *PostgresChannelRepository) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Channel, error) {
	query := `
		SELECT c.id, c.name, c.topic, c.creator_user_id, c.creation_time, c.last_message_time, c.is_archived, c.version
		FROM channels c
		JOIN channel_invites ci ON c.id = ci.channel_id
		WHERE ci.invite_code = $1
	`

	row := r.pool.QueryRow(ctx, query, inviteCode)
	channel, err := r.scanChannelBasic(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel with invite code %s not found: %w", inviteCode, common.ErrNotFound)
		}
		return nil, fmt.Errorf("error scanning channel for invite code %s: %w", inviteCode, err)
	}

	// Load members and invites
	members, err := r.loadMembers(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	invites, err := r.loadInvites(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	channel.Members = members
	channel.Invites = invites
	channel.Messages = []models.Message{}

	return channel, nil
}

func (r *PostgresChannelRepository) FindByInviteID(ctx context.Context, inviteID uuid.UUID) (*models.Channel, error) {
	query := `
		SELECT c.id, c.name, c.topic, c.creator_user_id, c.creation_time, c.last_message_time, c.is_archived, c.version
		FROM channels c
		JOIN channel_invites ci ON c.id = ci.channel_id
		WHERE ci.id = $1
	`

	row := r.pool.QueryRow(ctx, query, inviteID)
	channel, err := r.scanChannelBasic(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel with invite ID %s not found: %w", inviteID, common.ErrNotFound)
		}
		return nil, fmt.Errorf("error scanning channel for invite ID %s: %w", inviteID, err)
	}

	// Load members and invites
	members, err := r.loadMembers(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	invites, err := r.loadInvites(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	channel.Members = members
	channel.Invites = invites
	channel.Messages = []models.Message{}

	return channel, nil
}

func (r *PostgresChannelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM channels WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting channel %s: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("channel with ID %s was not found for deletion: %w", id, common.ErrNotFound)
	}
	return nil
}

func (r *PostgresChannelRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (
			id, channel_id, sender_user_id, integration_id,
			content_text, content_mentions, content_link, content_formatted,
			created_at, parent_message_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var senderID, integrationID, parentID interface{}
	if message.GetSenderUserId() != nil {
		senderID = *message.GetSenderUserId()
	}
	if message.GetIntegrationId() != nil {
		integrationID = *message.GetIntegrationId()
	}
	if message.GetParentMessageId() != nil {
		parentID = *message.GetParentMessageId()
	}

	_, err := r.pool.Exec(ctx, query,
		message.GetId(),
		message.GetChannelId(),
		senderID,
		integrationID,
		message.GetContent().GetText(),
		message.GetContent().GetMentions(),
		message.GetContent().GetLinks(),
		message.GetContent().GetIsFormatted(),
		message.GetCreatedAt(),
		parentID,
	)
	if err != nil {
		return fmt.Errorf("error inserting message %s for channel %s: %w", message.GetId(), message.GetChannelId(), err)
	}
	return nil
}

func (r *PostgresChannelRepository) SaveReaction(ctx context.Context, reaction *models.Reaction) error {
	query := `
		INSERT INTO reactions (id, message_id, user_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		reaction.GetId(),
		reaction.GetMessageId(),
		reaction.GetUserId(),
		reaction.GetReactionType(),
		reaction.GetCreatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("reaction already exists for message %s by user %s with type %s: %w",
				reaction.GetMessageId(), reaction.GetUserId(), reaction.GetReactionType(), common.ErrConflict)
		}
		return fmt.Errorf("error inserting reaction %s: %w", reaction.GetId(), err)
	}
	return nil
}

func (r *PostgresChannelRepository) DeleteReaction(ctx context.Context, messageID, userID uuid.UUID, reactionType string) error {
	query := `
		DELETE FROM reactions
		WHERE message_id = $1 AND user_id = $2 AND reaction_type = $3
	`

	cmdTag, err := r.pool.Exec(ctx, query, messageID, userID, reactionType)
	if err != nil {
		return fmt.Errorf("error deleting reaction (type: %s) from message %s by user %s: %w",
			reactionType, messageID, userID, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reaction (type: %s) not found on message %s for user %s to delete: %w",
			reactionType, messageID, userID, common.ErrNotFound)
	}
	return nil
}

func (r *PostgresChannelRepository) FindReactionsByMessageID(ctx context.Context, messageID uuid.UUID) ([]models.Reaction, error) {
	query := `
		SELECT id, message_id, user_id, reaction_type, created_at
		FROM reactions
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying reactions for message %s: %w", messageID, err)
	}
	defer rows.Close()

	var reactions []models.Reaction
	for rows.Next() {
		var reactionID uuid.UUID
		var msgID uuid.UUID
		var userID uuid.UUID
		var reactionType string
		var timestamp time.Time

		err := rows.Scan(&reactionID, &msgID, &userID, &reactionType, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("error scanning reaction for message %s: %w", messageID, err)
		}
		reaction := models.RehydrateReaction(reactionID, msgID, userID, reactionType, timestamp)
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reactions for message %s: %w", messageID, err)
	}
	return reactions, nil
}
