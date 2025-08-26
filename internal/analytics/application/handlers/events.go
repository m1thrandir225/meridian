package handlers

import "time"

type UserRegisteredEvent struct {
	ID          string    `json:"ID"`
	Name        string    `json:"Name"`
	Time        time.Time `json:"Time"`
	AggrID      string    `json:"AggrID"`
	AggrType    string    `json:"AggrType"`
	AggrVersion int64     `json:"AggrVersion"`
	UserID      string    `json:"UserID"`
	Username    string    `json:"Username"`
	Email       string    `json:"Email"`
	FirstName   string    `json:"FirstName"`
	LastName    string    `json:"LastName"`
	Timestamp   time.Time `json:"Timestamp"`
}

type UserProfileUpdatedEvent struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Time          time.Time              `json:"time"`
	AggrID        string                 `json:"aggr_id"`
	AggrType      string                 `json:"aggr_type"`
	AggrVersion   int64                  `json:"aggr_version"`
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
	Timestamp     time.Time              `json:"timestamp"`
}

// Fix the MessageSentEvent structure to match the actual domain event
type MessageSentEvent struct {
	ID            string  `json:"ID"`
	Name          string  `json:"Name"` // Go default: Name
	Time          string  `json:"Time"`
	AggrID        string  `json:"AggrID"`
	AggrType      string  `json:"AggrType"`
	AggrVersion   int64   `json:"AggrVersion"`
	MessageID     string  `json:"MessageID"`
	SenderUserID  *string `json:"SenderUserID"`
	IntegrationID *string `json:"IntegrationID"`
	Content       struct {
		Text      string   `json:"Text"`
		Mentions  []string `json:"Mentions"`
		Links     []string `json:"Links"`
		Formatted bool     `json:"Formatted"`
	} `json:"Content"` // Go default: Content
	Timestamp       string  `json:"Timestamp"`
	ParentMessageID *string `json:"ParentMessageID"`
}

// Fix the ChannelCreatedEvent structure
type ChannelCreatedEvent struct {
	ID            string    `json:"ID"`
	Name          string    `json:"Name"`
	Time          time.Time `json:"Time"`
	AggrID        string    `json:"AggrID"`
	AggrType      string    `json:"AggrType"`
	AggrVersion   int64     `json:"AggrVersion"`
	ChannelName   string    `json:"ChannelName"`
	CreatorUserID string    `json:"CreatorUserID"`
	Topic         string    `json:"Topic"`
}

type UserJoinedChannelEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Time        time.Time `json:"time"`
	AggrID      string    `json:"aggr_id"`
	AggrType    string    `json:"aggr_type"`
	AggrVersion int64     `json:"aggr_version"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
}

type UserLeftChannelEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Time        time.Time `json:"time"`
	AggrID      string    `json:"aggr_id"`
	AggrType    string    `json:"aggr_type"`
	AggrVersion int64     `json:"aggr_version"`
	UserID      string    `json:"user_id"`
}

type ReactionAddedEvent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Time         time.Time `json:"time"`
	AggrID       string    `json:"aggr_id"`
	AggrType     string    `json:"aggr_type"`
	AggrVersion  int64     `json:"aggr_version"`
	ReactionID   string    `json:"reaction_id"`
	MessageID    string    `json:"message_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type ReactionRemovedEvent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Time         time.Time `json:"time"`
	AggrID       string    `json:"aggr_id"`
	AggrType     string    `json:"aggr_type"`
	AggrVersion  int64     `json:"aggr_version"`
	MessageID    string    `json:"message_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
}

type IntegrationRegisteredEvent struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Time             time.Time `json:"time"`
	AggrID           string    `json:"aggr_id"`
	AggrType         string    `json:"aggr_type"`
	AggrVersion      int64     `json:"aggr_version"`
	IntegrationID    string    `json:"integrationId"`
	ServiceName      string    `json:"serviceName"`
	CreatorUserID    string    `json:"creatorUserId"`
	TargetChannelIDs []string  `json:"targetChannelIds"`
	RegisteredAt     time.Time `json:"registeredAt"`
}
