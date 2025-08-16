package handlers

type RegisterIntegrationRequest struct {
	ServiceName      string   `json:"service_name" binding:"required"`
	TargetChannelIDs []string `json:"target_channel_ids" binding:"required"`
}

type RevokeIntegrationRequest struct {
	IntegrationID string `json:"integration_id" binding:"required"`
}

type WebhookMessageRequest struct {
	ContentText     string            `json:"content_text" binding:"required"`
	TargetChannelID string            `json:"target_channel_id,omitempty"` // Optional, defaults to first target channel
	ParentMessageID *string           `json:"parent_message_id,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type CallbackMessageRequest struct {
	ContentText     string            `json:"content_text" binding:"required"`
	TargetChannelID string            `json:"target_channel_id,omitempty"`
	ParentMessageID *string           `json:"parent_message_id,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}
