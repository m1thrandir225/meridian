package handlers

import "time"

type MessageSentEvent struct {
	ID            string  `json:"ID"`
	Name          string  `json:"Name"`
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
	} `json:"Content"`
	Timestamp       string  `json:"Timestamp"`
	ParentMessageID *string `json:"ParentMessageID"`
}

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
	ID          string    `json:"ID"`
	Name        string    `json:"Name"`
	Time        time.Time `json:"Time"`
	AggrID      string    `json:"AggrID"`
	AggrType    string    `json:"AggrType"`
	AggrVersion int64     `json:"AggrVersion"`
	UserID      string    `json:"UserID"`
	Role        string    `json:"Role"`
	JoinedAt    time.Time `json:"JoinedAt"`
}

type UserLeftChannelEvent struct {
	ID          string    `json:"ID"`
	Name        string    `json:"Name"`
	Time        time.Time `json:"Time"`
	AggrID      string    `json:"AggrID"`
	AggrType    string    `json:"AggrType"`
	AggrVersion int64     `json:"AggrVersion"`
	UserID      string    `json:"UserID"`
}

type ReactionAddedEvent struct {
	ID           string    `json:"ID"`
	Name         string    `json:"Name"`
	Time         time.Time `json:"Time"`
	AggrID       string    `json:"AggrID"`
	AggrType     string    `json:"AggrType"`
	AggrVersion  int64     `json:"AggrVersion"`
	ReactionID   string    `json:"ReactionID"`
	MessageID    string    `json:"MessageID"`
	UserID       string    `json:"UserID"`
	ReactionType string    `json:"ReactionType"`
	Timestamp    time.Time `json:"Timestamp"`
}

type ReactionRemovedEvent struct {
	ID           string    `json:"ID"`
	Name         string    `json:"Name"`
	Time         time.Time `json:"Time"`
	AggrID       string    `json:"AggrID"`
	AggrType     string    `json:"AggrType"`
	AggrVersion  int64     `json:"AggrVersion"`
	MessageID    string    `json:"MessageID"`
	UserID       string    `json:"UserID"`
	ReactionType string    `json:"ReactionType"`
}

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
	ID            string                 `json:"ID"`
	Name          string                 `json:"Name"`
	Time          time.Time              `json:"Time"`
	AggrID        string                 `json:"AggrID"`
	AggrType      string                 `json:"AggrType"`
	AggrVersion   int64                  `json:"AggrVersion"`
	UserID        string                 `json:"UserID"`
	UpdatedFields map[string]interface{} `json:"UpdatedFields"`
	Timestamp     time.Time              `json:"Timestamp"`
}

type IntegrationRegisteredEvent struct {
	ID               string    `json:"ID"`
	Name             string    `json:"Name"`
	Time             time.Time `json:"Time"`
	AggrID           string    `json:"AggrID"`
	AggrType         string    `json:"AggrType"`
	AggrVersion      int64     `json:"AggrVersion"`
	IntegrationID    string    `json:"IntegrationID"`
	ServiceName      string    `json:"ServiceName"`
	CreatorUserID    string    `json:"CreatorUserID"`
	TargetChannelIDs []string  `json:"TargetChannelIDs"`
	RegisteredAt     time.Time `json:"RegisteredAt"`
}
