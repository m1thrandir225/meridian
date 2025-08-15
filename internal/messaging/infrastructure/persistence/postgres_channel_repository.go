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

func (r *PostgresChannelRepository) Save(ctx context.Context, channel *models.Channel) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beggining transaction: %w", err)
	}
	defer tx.Rollback(ctx)

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

	deleteMembershipQuery := `DELETE FROM members WHERE channel_id = $1`
	_, err = tx.Exec(ctx, deleteMembershipQuery, channel.ID)
	if err != nil {
		return fmt.Errorf("error deleting old members for channel %s: %w", channel.ID, err)
	}
	if len(channel.Members) > 0 {
		memberRows := [][]interface{}{}
		for _, member := range channel.Members {
			memberRows = append(memberRows, []interface{}{
				channel.ID, member.GetId(), member.GetRole(), member.GetJoinedAt(), member.GetLastRead(),
			})
		}
		copyCount, err := tx.CopyFrom(
			ctx,
			pgx.Identifier{"members"},
			[]string{"channel_id", "user_id", "role", "joined_at", "last_read"},
			pgx.CopyFromRows(memberRows),
		)
		if err != nil {
			return fmt.Errorf("error inserting members for channel %s using copy: %w", channel.ID, err)
		}

		if int(copyCount) != len(channel.Members) {
			return fmt.Errorf("expected %d members to be inserted for channel %s, but got %d", len(channel.Members), channel.ID, copyCount)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting transaction for channel %s: %w", channel.ID, err)
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
		var channel models.Channel
		var topic *string
		var lastMsgTime *time.Time

		err := rows.Scan(
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
			return nil, fmt.Errorf("error scanning channel for user %s: %w", userID, err)
		}

		if topic != nil {
			channel.Topic = *topic
		}
		if lastMsgTime != nil {
			channel.LastMessageTime = *lastMsgTime
		}

		// load members for each channel
		memberQuery := `
			SELECT user_id, role, joined_at, last_read
			FROM members
			WHERE channel_id = $1
		`
		memberRows, err := r.pool.Query(ctx, memberQuery, channel.ID)
		if err != nil {
			return nil, fmt.Errorf("error querying members for channel %s: %w", channel.ID, err)
		}

		var rehydratedMembers []models.Member
		for memberRows.Next() {
			var memberUserID uuid.UUID
			var memberRole string
			var memberJoinedAt time.Time
			var memberLastRead sql.NullTime

			if err := memberRows.Scan(&memberUserID, &memberRole, &memberJoinedAt, &memberLastRead); err != nil {
				memberRows.Close()
				return nil, fmt.Errorf("error scanning member for channel %s: %w", channel.ID, err)
			}

			var actualLastRead time.Time
			if memberLastRead.Valid {
				actualLastRead = memberLastRead.Time
			}

			member := models.RehydrateMember(memberUserID, memberRole, memberJoinedAt, actualLastRead)
			rehydratedMembers = append(rehydratedMembers, member)
		}
		memberRows.Close()

		if err := memberRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating members for channel %s: %w", channel.ID, err)
		}

		channel.Members = rehydratedMembers
		channel.Messages = []models.Message{} // don't load them

		channels = append(channels, &channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating channels for user %s: %w", userID, err)
	}

	return channels, nil
}

func (r *PostgresChannelRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	queryChannel := `
		SELECT id, name, topic, creator_user_id, creation_time, last_message_time, is_archived, version
		FROM channels
		WHERE id = $1
	`
	var channel models.Channel
	var topic *string
	var lastMsgTime *time.Time

	err := r.pool.QueryRow(ctx, queryChannel, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel with ID %s not found: %w", id, common.ErrNotFound)
		}
		return nil, fmt.Errorf("error scanning channel %s: %w", id, err)
	}

	if topic != nil {
		channel.Topic = *topic
	}

	if lastMsgTime != nil {
		channel.LastMessageTime = *lastMsgTime
	}
	queryMembers := `
		SELECT user_id, role, joined_at, last_read
		FROM members
		WHERE channel_id = $1
	`
	rows, err := r.pool.Query(ctx, queryMembers, id)
	if err != nil {
		return nil, fmt.Errorf("error querrying members for channel %s: %w", id, err)
	}
	var rehydratedMembers []models.Member
	for rows.Next() {
		var memberUserID uuid.UUID
		var memberRole string
		var memberJoinedAt time.Time
		var memberLastRead sql.NullTime

		err := rows.Scan(&memberUserID, &memberRole, &memberJoinedAt, &memberLastRead)
		if err != nil {
			return nil, fmt.Errorf("error scanning member for channel %s: %w", id, err)
		}
		var actualLastRead time.Time
		if memberLastRead.Valid {
			actualLastRead = memberLastRead.Time
		}
		member := models.RehydrateMember(memberUserID, memberRole, memberJoinedAt, actualLastRead)
		rehydratedMembers = append(rehydratedMembers, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members for channel %s: %w", id, err)
	}
	channel.Members = rehydratedMembers
	channel.Messages = []models.Message{} // dont load them
	return &channel, nil
}

func (r *PostgresChannelRepository) FindMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]models.Message, error) {
	queryMessages := `
		SELECT id, channel_id, sender_user_id, integration_id,
		       content_text, content_mentions, content_link, content_formatted,
		       created_at, parent_message_id
		FROM messages
		WHERE channel_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.pool.Query(ctx, queryMessages, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying messages for channel %s: %w", channelID, err)
	}
	defer rows.Close()

	messages := []models.Message{}
	messageIDs := []uuid.UUID{}

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

		var pSenderUserID, pIntegrationID, pParentMessageID *uuid.UUID
		if senderUserID != nil {
			pSenderUserID = senderUserID
		}
		if integrationID != nil {
			pIntegrationID = integrationID
		}
		if parentMessageID != nil {
			pParentMessageID = parentMessageID
		}

		msg := models.RehydrateMessage(
			messageId,
			channelID,
			pSenderUserID,
			pIntegrationID,
			pParentMessageID,
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
	rows.Close()

	if len(messages) == 0 {
		return messages, nil
	}

	queryReactions := `
		SELECT id, message_id, user_id, reaction_type, created_at
		FROM reactions
		WHERE message_id = ANY($1) -- Use ANY with an array of UUIDs
		ORDER BY message_id, created_at ASC`

	reactionRows, err := r.pool.Query(ctx, queryReactions, messageIDs)
	if err != nil {
		return nil, fmt.Errorf("error querying reactions for messages: %w", err)
	}
	defer reactionRows.Close()

	reactionsByMessageID := make(map[uuid.UUID][]models.Reaction)
	for reactionRows.Next() {
		var reactionID uuid.UUID
		var reactionMessageID uuid.UUID
		var reactionUserID uuid.UUID
		var reactionType string
		var reactionTimestamp time.Time

		err := reactionRows.Scan(
			&reactionID,
			&reactionMessageID,
			&reactionUserID,
			&reactionType,
			&reactionTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reaction: %w", err)
		}

		reaction := models.RehydrateReaction(reactionID, reactionMessageID, reactionUserID, reactionType, reactionTimestamp)
		reactionsByMessageID[reactionMessageID] = append(reactionsByMessageID[reactionMessageID], reaction)
	}
	if err := reactionRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reactions: %w", err)
	}

	for i := range messages {
		if reactions, ok := reactionsByMessageID[messages[i].GetId()]; ok {
			messages[i].SetLoadedReactions(reactions)
		}
	}

	return messages, nil
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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
		VALUES ($1, $2, $3, $4, $5)`

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
		WHERE message_id = $1 AND user_id = $2 AND reaction_type = $3`

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
		ORDER BY created_at ASC`

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
